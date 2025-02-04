package entities

import (
	"time"
)

type Encounter struct {
	Id            int       `db:"id,key,auto"`
	EncounterId   int       `db:"encounter_id"`
	EncounterName string    `db:"encounter_name"`
	ZoneId        int       `db:"zone_id"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
