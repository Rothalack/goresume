package entities

import (
	"time"
)

type Zone struct {
	Id          int       `db:"id,key,auto"`
	ZoneId      int       `db:"zone_id"`      //unique key zone_id+expansion_id
	ExpansionId int       `db:"expansion_id"` //unique key zone_id+expansion_id
	ZoneName    string    `db:"zone_name"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
