package entities

import (
	"time"
)

type Size struct {
	Id           int       `db:"id,key,auto"`   //unique key id+difficulty_id
	DifficultyId int       `db:"difficulty_id"` //unique key id+difficulty_id
	Size         string    `db:"size"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
