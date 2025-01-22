package entities

import (
	"time"
)

type User struct {
	ID        int       `db:"id,key,auto"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
