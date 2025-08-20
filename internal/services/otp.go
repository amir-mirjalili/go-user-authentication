package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/amir-mirjalili/go-user-authentication/internal/models"
)

type OTPRepository interface {
	SaveOTP(otp *models.OTP) error
	GetOTP(phoneNumber string) (*models.OTP, error)
	DeleteOTP(phoneNumber string) error
	CountOTPRequests(phoneNumber string, since time.Time) (int, error)
}

type OTPService struct {
	repo OTPRepository
}

func NewOTPService(repo OTPRepository) *OTPService {
	return &OTPService{repo: repo}
}

func (s *OTPService) GenerateOTP(phoneNumber string) (*models.OTP, error) {
	// Check rate limiting
	since := time.Now().Add(-10 * time.Minute)
	count, err := s.repo.CountOTPRequests(phoneNumber, since)
	if err != nil {
		return nil, err
	}

	if count >= 3 {
		return nil, fmt.Errorf("rate limit exceeded: maximum 3 OTP requests per 10 minutes")
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	otp := &models.OTP{
		PhoneNumber: phoneNumber,
		Code:        code,
		ExpiresAt:   time.Now().Add(2 * time.Minute),
		CreatedAt:   time.Now(),
	}

	err = s.repo.SaveOTP(otp)
	if err != nil {
		return nil, err
	}

	fmt.Printf("OTP for %s: %s (expires at %s)\n", phoneNumber, code, otp.ExpiresAt.Format("15:04:05"))

	return otp, nil
}

func (s *OTPService) VerifyOTP(phoneNumber, code string) (bool, error) {
	otp, err := s.repo.GetOTP(phoneNumber)
	if err != nil {
		return false, fmt.Errorf("invalid OTP")
	}

	if time.Now().After(otp.ExpiresAt) {
		s.repo.DeleteOTP(phoneNumber)
		return false, fmt.Errorf("OTP expired")
	}

	if otp.Code != code {
		return false, fmt.Errorf("invalid OTP")
	}

	// Clean up used OTP
	s.repo.DeleteOTP(phoneNumber)
	return true, nil
}
