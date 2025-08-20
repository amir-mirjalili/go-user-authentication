package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/amir-mirjalili/go-user-authentication/internal/models"
	"github.com/amir-mirjalili/go-user-authentication/internal/params"
)

type UserRepository struct {
	Postgres *sql.DB
}

func NewUserRepository(postgres *sql.DB) *UserRepository {
	return &UserRepository{
		Postgres: postgres,
	}
}

func (p *UserRepository) CreateUser(user *params.UserRegisterRequest) (*models.User, error) {
	query := `
		INSERT INTO users (phone_number, registered_at)
		VALUES ($1, $2)
		RETURNING id, phone_number, registered_at
	`

	var createdUser models.User
	err := p.Postgres.QueryRow(query, user.PhoneNumber, user.RegisteredAt).
		Scan(&createdUser.ID, &createdUser.PhoneNumber, &createdUser.RegisteredAt)
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (p *UserRepository) GetUserByPhone(phoneNumber string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, phone_number, registered_at, created_at, updated_at FROM users WHERE phone_number = $1`
	err := p.Postgres.QueryRow(query, phoneNumber).Scan(&user.ID, &user.PhoneNumber, &user.RegisteredAt, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (p *UserRepository) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, phone_number, registered_at, created_at, updated_at FROM users WHERE id = $1`
	err := p.Postgres.QueryRow(query, id).Scan(&user.ID, &user.PhoneNumber, &user.RegisteredAt, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (p *UserRepository) ListUsers(page, limit int, search string) ([]models.User, int, error) {
	var users []models.User
	var total int

	// Count total
	countQuery := `SELECT COUNT(*) FROM users WHERE phone_number ILIKE '%' || $1 || '%'`
	err := p.Postgres.QueryRow(countQuery, search).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get users
	offset := (page - 1) * limit
	query := `SELECT id, phone_number, registered_at, created_at, updated_at 
			  FROM users 
			  WHERE phone_number ILIKE '%' || $1 || '%' 
			  ORDER BY id 
			  LIMIT $2 OFFSET $3`

	rows, err := p.Postgres.Query(query, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.PhoneNumber, &user.RegisteredAt, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}
