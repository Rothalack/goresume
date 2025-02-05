package entities

import (
	"time"
)

type Character struct {
	Id        int       `db:"id,key,auto"`
	CharName  string    `db:"char_name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
