package commands

import (
	"goresume/controllers/warcraftlogs"
)

func SyncBaseDataCommand() {
	warcraftlogs.SyncData()
}
