package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/amir-mirjalili/go-user-authentication/internal/models"
)

type OtpRepository struct {
	Postgres *sql.DB
}

func (p *OtpRepository) CountOTPRequests(phoneNumber string, since time.Time) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM otps WHERE phone_number = $1 AND created_at >= $2`
	err := p.Postgres.QueryRow(query, phoneNumber, since).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func NewOtpRepository(postgres *sql.DB) *OtpRepository {
	return &OtpRepository{
		Postgres: postgres,
	}
}

func (p *OtpRepository) SaveOTP(otp *models.OTP) error {
	_, err := p.Postgres.Exec(`INSERT INTO otp_requests (phone_number) VALUES ($1)`, otp.PhoneNumber)
	if err != nil {
		return err
	}

	query := `INSERT INTO otps (phone_number, code, expires_at) VALUES ($1, $2, $3)
			  ON CONFLICT (phone_number)
			  DO UPDATE SET code = $2, expires_at = $3, created_at = NOW()`
	_, err = p.Postgres.Exec(query, otp.PhoneNumber, otp.Code, otp.ExpiresAt)
	return err
}

func (p *OtpRepository) GetOTP(phoneNumber string) (*models.OTP, error) {
	otp := &models.OTP{}
	query := `SELECT phone_number, code, expires_at, created_at FROM otps WHERE phone_number = $1`
	err := p.Postgres.QueryRow(query, phoneNumber).Scan(&otp.PhoneNumber, &otp.Code, &otp.ExpiresAt, &otp.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("OTP not found")
		}
		return nil, err
	}
	return otp, nil
}

func (p *OtpRepository) DeleteOTP(phoneNumber string) error {
	_, err := p.Postgres.Exec(`DELETE FROM otps WHERE phone_number = $1`, phoneNumber)
	return err
}

func (p *UserRepository) CountOTPRequests(phoneNumber string, since time.Time) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM otp_requests WHERE phone_number = $1 AND requested_at > $2`
	err := p.Postgres.QueryRow(query, phoneNumber, since).Scan(&count)
	return count, err
}
