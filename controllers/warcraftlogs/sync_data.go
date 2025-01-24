package warcraftlogs

import (
	"fmt"
	"goresume/config"
	"goresume/config/entities"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func SyncData() {
	rows, err := config.DB.Query("SELECT id, api_url, note FROM games;")
	if err != nil {
		log.Fatalf("Failed to query games: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var game entities.Game
		if err := rows.Scan(&game.Id, &game.ApiUrl, &game.Note); err != nil {
			log.Printf("error scanning game row: %v", err)
			continue
		}

		SetGameContext(game.ApiUrl)
		log.Printf("running syncExpansions. gameId: %v", game.Id)
		syncExpansions(game.Id)
		log.Printf("running syncRegions. gameId: %v", game.Id)
		syncRegions(game.Id)
	}
}

func syncExpansions(gameId int) ([]Expansion, error) {
	expansions, err := GetExpansions()
	if err != nil {
		return nil, fmt.Errorf("failed to get regions: %v", err)
	}

	tx, err := config.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %v", err)
	}

	for _, expansion := range expansions {
		_, err := tx.Exec(`
            INSERT IGNORE INTO expansions (game_id, expansion_id, expansion_name)
            VALUES (?, ?, ?)
        `, gameId, expansion.Id, expansion.Name)
		if err != nil {
			log.Printf("failed to insert expansion %d: %v", expansion.Id, err)
			continue
		}

		log.Printf("running syncZones. expansionId: %v", expansion.Id)
		syncZones(expansion.Id, expansion.Zones)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit expansions transaction: %v", err)
	}

	return expansions, nil
}

func syncRegions(gameId int) ([]Region, error) {
	regions, err := GetRegions()
	if err != nil {
		return nil, fmt.Errorf("failed to get regions: %v", err)
	}

	tx, err := config.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %v", err)
	}

	for _, region := range regions {
		_, err := tx.Exec(`
            INSERT IGNORE INTO regions (game_id, region_id, name, compact_name, slug)
            VALUES (?, ?, ?, ?, ?)
        `, gameId, region.Id, region.Name, region.CompactName, region.Slug)
		if err != nil {
			log.Printf("failed to insert region %d: %v", region.Id, err)
			continue
		}

		log.Printf("running syncServers. region.Id: %v", region.Id)
		syncServers(region.Id)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit region transaction: %v", err)
	}

	return regions, nil
}

func syncServers(regionId int) error {
	servers, err := GetServersFromRegion(regionId)
	if err != nil {
		return fmt.Errorf("failed to get servers: %v", err)
	}

	tx, err := config.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to open db: %v", err)
	}

	for _, server := range servers {
		_, err := tx.Exec(`
            INSERT IGNORE INTO servers (region_id, server_id, server_name, normalized_name, slug)
            VALUES (?, ?, ?, ?, ?)
        `, regionId, server.Id, server.Name, server.NormalizedName, server.Slug)
		if err != nil {
			log.Printf("failed to insert server %d: %v", server.Id, err)
			continue
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit server transaction: %v", err)
	}

	return nil
}

func syncZones(expansionId int, zones []Zone) ([]Zone, error) {
	tx, _ := config.DB.Begin()

	for _, zone := range zones {
		_, err := tx.Exec(`
            INSERT IGNORE INTO zones (zone_id, expansion_id, zone_name)
            VALUES (?, ?, ?)
        `, zone.Id, expansionId, zone.Name)
		if err != nil {
			log.Printf("failed to insert zone %d: %v", zone.Id, err)
			continue
		}

		log.Printf("running syncPartitions. zone.Id: %v", zone.Id)
		syncPartitions(zone)
		log.Printf("running syncDifficulties. zone.Id: %v", zone.Id)
		syncDifficulties(zone)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit zone transaction: %v", err)
	}

	return zones, nil
}

func syncPartitions(zone Zone) error {
	tx, _ := config.DB.Begin()

	for _, paritition := range zone.Partitions {
		_, err := tx.Exec(`
            INSERT IGNORE INTO parititions (partition_id, zone_id, name, compact_name, default_status)
            VALUES (?, ?, ?, ?, ?)
        `, paritition.Id, zone.Id, paritition.Name, paritition.CompactName, paritition.DefaultStatus)
		if err != nil {
			log.Printf("failed to insert partition %d: %v", zone.Id, err)
			continue
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit partition transaction: %v", err)
	}

	return nil
}

func syncDifficulties(zone Zone) error {
	tx, _ := config.DB.Begin()

	for _, difficulty := range zone.Difficulties {
		_, err := tx.Exec(`
            INSERT IGNORE INTO difficulties (difficulty_id, zone_id, name)
            VALUES (?, ?, ?)
        `, difficulty.Id, zone.Id, difficulty.Name)
		if err != nil {
			log.Printf("failed to insert difficulty %d: %v", zone.Id, err)
			continue
		}
		for _, size := range difficulty.Sizes {
			_, err := tx.Exec(`
				INSERT IGNORE INTO sizes (difficulty_id, size)
				VALUES (?, ?)
			`, difficulty.Id, size)
			if err != nil {
				log.Printf("failed to insert size %d: %v", zone.Id, err)
				continue
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit difficulties transaction: %v", err)
	}

	return nil
}
