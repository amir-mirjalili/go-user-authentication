package services

import (
	"time"

	"github.com/amir-mirjalili/go-user-authentication/internal/params"
	"github.com/amir-mirjalili/go-user-authentication/internal/pkg/jwt"
)

type AuthService struct {
	otpService  *OTPService
	userService *UserService
	jwtSecret   string
}

func NewAuthService(otpService *OTPService, userService *UserService, jwtSecret string) *AuthService {
	return &AuthService{
		otpService:  otpService,
		userService: userService,
		jwtSecret:   jwtSecret,
	}
}

func (s *AuthService) SendOTP(phoneNumber string) error {
	_, err := s.otpService.GenerateOTP(phoneNumber)
	return err
}

func (s *AuthService) VerifyOTPAndLogin(phoneNumber, code string) (*params.AuthResponse, error) {
	// Verify OTP
	valid, err := s.otpService.VerifyOTP(phoneNumber, code)
	if err != nil || !valid {
		return nil, err
	}

	// Check if user exists
	user, err := s.userService.GetUserByPhone(phoneNumber)
	if err != nil {
		// User doesn't exist, create new user
		input := &params.UserRegisterRequest{
			PhoneNumber:  phoneNumber,
			RegisteredAt: time.Now(),
		}
		user, err = s.userService.CreateUser(input)
		if err != nil {
			return nil, err
		}
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID, user.PhoneNumber, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &params.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}
