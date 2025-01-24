package entities

import (
	"time"
)

type Expansion struct {
	Id            int       `db:"id,key,auto"`
	GameId        int       `db:"game_id"`      //unique key game_id+expansion_id
	ExpansionId   int       `db:"expansion_id"` //unique key game_id+expansion_id
	ExpansionName string    `db:"expansion_name"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
