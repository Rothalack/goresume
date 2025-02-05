package entities

import (
	"time"
)

type Game struct {
	Id        int       `db:"id,key,auto"`
	GameName  string    `db:"game_name"`
	ApiUrl    string    `db:"api_url"`
	Note      string    `db:"note"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
