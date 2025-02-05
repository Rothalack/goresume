package entities

import (
	"time"
)

type Difficulty struct {
	Id           int       `db:"id,key,auto"`
	DifficultyId int       `db:"difficulty_id"` //unique key difficulty_id+zone_id
	ZoneId       int       `db:"zone_id"`       //unique key difficulty_id+zone_id
	Name         string    `db:"name"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
