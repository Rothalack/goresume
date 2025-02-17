package entities

import "time"

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	UserName     *string   `json:"user_name,omitempty"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	GoogleID     *string   `json:"google_id,omitempty"`
}
