package entities

import (
	"time"
)

type Server struct {
	Id             int       `db:"id,key,auto"`
	RegionId       int       `db:"region_id"` //unique key region_id+server_id
	ServerId       int       `db:"server_id"` //unique key region_id+server_id
	ServerName     string    `db:"server_name"`
	NormalizedName string    `db:"normalized_name"`
	Slug           string    `db:"slug"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
