package models

import "time"

type User struct {
	ID           int       `json:"id" db:"id"`
	PhoneNumber  string    `json:"phone_number" db:"phone_number"`
	RegisteredAt time.Time `json:"registered_at" db:"registered_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
