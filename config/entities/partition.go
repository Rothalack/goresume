package entities

import (
	"time"
)

type Partition struct {
	Id            int       `db:"id,key,auto"`
	PartitionId   int       `db:"partition_id"` //unique key partition_id+zone_id
	ZoneId        int       `db:"zone_id"`      //unique key partition_id+zone_id
	Name          string    `db:"name"`
	CompactName   string    `db:"compact_name"`
	DefaultStatus bool      `db:"default_status"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
