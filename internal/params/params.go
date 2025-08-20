package params

import (
	"time"

	"github.com/amir-mirjalili/go-user-authentication/internal/models"
)

type OTPRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required,e164" example:"+1234567890"`
}

type OTPVerifyRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required,e164" example:"+1234567890"`
	Code        string `json:"code" binding:"required,len=6,numeric" example:"123456"`
}

type AuthResponse struct {
	Token string      `json:"token" `
	User  models.User `json:"user"`
}

type UserListResponse struct {
	Users      []models.User `json:"users"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"total_pages"`
}

type UserRegisterRequest struct {
	PhoneNumber  string    `json:"phone_number" db:"phone_number"`
	RegisteredAt time.Time `json:"registered_at" db:"registered_at"`
}
