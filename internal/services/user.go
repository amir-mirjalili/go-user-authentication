package services

import (
	"math"

	"github.com/amir-mirjalili/go-user-authentication/internal/models"
	"github.com/amir-mirjalili/go-user-authentication/internal/params"
)

type UserRepository interface {
	CreateUser(user *params.UserRegisterRequest) (*models.User, error)
	GetUserByPhone(phoneNumber string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	ListUsers(page, limit int, search string) ([]models.User, int, error)
}

type UserService struct {
	repository UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repository: repo}
}

func (s *UserService) GetUser(id int) (*models.User, error) {
	return s.repository.GetUserByID(id)
}

func (s *UserService) GetUserByPhone(phoneNumber string) (*models.User, error) {
	return s.repository.GetUserByPhone(phoneNumber)
}

func (s *UserService) ListUsers(page, limit int, search string) (*params.UserListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := s.repository.ListUsers(page, limit, search)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &params.UserListResponse{
		Users:      users,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *UserService) CreateUser(user *params.UserRegisterRequest) (*models.User, error) {
	entity, err := s.repository.CreateUser(user)
	if err != nil {
		return &models.User{}, err
	}
	return entity, nil
}
