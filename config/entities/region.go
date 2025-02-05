package entities

import (
	"time"
)

type Region struct {
	Id          int       `db:"id,key,auto"`
	GameId      int       `db:"game_id"`   //unique key game_id+region_id
	RegionId    int       `db:"region_id"` //unique key game_id+region_id
	Name        string    `db:"name"`
	CompactName string    `db:"compact_name"`
	Slug        string    `db:"slug"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
