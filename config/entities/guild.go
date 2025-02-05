package entities

import (
	"time"
)

type Guild struct {
	Id        int       `db:"id"`
	GuildName string    `db:"guild_name"`
	GuildID   int       `db:"guild_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
