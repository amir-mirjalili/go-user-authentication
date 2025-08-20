package models

import "time"

type OTP struct {
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Code        string    `json:"code" db:"code"`
	ExpiresAt   time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
