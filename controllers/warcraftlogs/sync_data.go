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
		log.Println("running syncClassesFactions")
	}
}

func syncExpansions(gameId int) ([]Expansion, error) {
	expansions, err := GetExpansions()
	if err != nil {
		return nil, fmt.Errorf("failed to get regions: %v", err)
	}

	for _, expansion := range expansions {
		_, err := config.DB.Exec(`
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

	return expansions, nil
}

func syncRegions(gameId int) ([]Region, error) {
	regions, err := GetRegions()
	if err != nil {
		return nil, fmt.Errorf("failed to get regions: %v", err)
	}

	for _, region := range regions {
		_, err := config.DB.Exec(`
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

	return regions, nil
}

func syncServers(regionId int) error {
	servers, err := GetServersFromRegion(regionId)

	if err != nil {
		return fmt.Errorf("failed to get servers: %v", err)
	}

	for _, server := range servers {
		_, err := config.DB.Exec(`
            INSERT IGNORE INTO servers (region_id, server_id, server_name, normalized_name, slug)
            VALUES (?, ?, ?, ?, ?)
        `, regionId, server.Id, server.Name, server.NormalizedName, server.Slug)
		if err != nil {
			log.Printf("failed to insert server %d: %v", server.Id, err)
			continue
		}
	}

	return nil
}

func syncZones(expansionId int, zones []Zone) ([]Zone, error) {
	for _, zone := range zones {
		_, err := config.DB.Exec(`
            INSERT IGNORE INTO zones (zone_id, expansion_id, zone_name)
            VALUES (?, ?, ?)
        `, zone.Id, expansionId, zone.Name)
		if err != nil {
			log.Printf("failed to insert zone %d: %v", zone.Id, err)
			continue
		}

		log.Printf("running syncDifficulties. zone.Id: %v", zone.Id)
		syncDifficulties(zone)
	}

	return zones, nil
}

func syncDifficulties(zone Zone) error {
	for _, difficulty := range zone.Difficulties {
		log.Printf("for difficulty in difficulties. difficulty: %v", difficulty)
		_, err := config.DB.Exec(`
            INSERT IGNORE INTO difficulties (difficulty_id, zone_id, name)
            VALUES (?, ?, ?)
        `, difficulty.Id, zone.Id, difficulty.Name)
		if err != nil {
			log.Printf("failed to insert difficulty %d: %v", zone.Id, err)
			continue
		}
		for _, size := range difficulty.Sizes {
			log.Printf("for size in difficulty.Sizes. size: %v", size)
			_, err := config.DB.Exec(`
				INSERT IGNORE INTO sizes (difficulty_id, zone_id, size)
				VALUES (?, ?, ?)
			`, difficulty.Id, zone.Id, size)
			if err != nil {
				log.Printf("failed to insert size %d: %v", zone.Id, err)
				continue
			}
		}
	}

	return nil
}
