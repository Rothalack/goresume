package warcraftlogs

import (
	"context"
	"fmt"
	"log"

	"github.com/machinebox/graphql"
)

type WorldData struct {
	Regions []Region `json:"regions"`
}

type Region struct {
	Id          int      `json:"id"`
	CompactName string   `json:"compactName"`
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Servers     []Server `json:"servers"`
}

type RegionsResponse struct {
	WorldData struct {
		Regions []Region `json:"regions"`
	} `json:"worldData"`
}

type Server struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	NormalizedName string `json:"normalizedName"`
	Slug           string `json:"slug"`
	Region         struct {
		Slug string `json:"slug"`
	} `json:"region"`
}

type ServerPagination struct {
	Data     []Server `json:"data"`
	LastPage int      `json:"last_page"`
}

type ServersResponse struct {
	WorldData struct {
		Region struct {
			Servers ServerPagination `json:"servers"`
		} `json:"region"`
	} `json:"worldData"`
}

type Expansion struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Zones []Zone `json:"zones"`
}

type ExpansionsResponse struct {
	WorldData struct {
		Expansions []Expansion `json:"expansions"`
	} `json:"worldData"`
}

type Zone struct {
	Id           int          `json:"id"`
	Name         string       `json:"name"`
	Partitions   []Partition  `json:"paritions"`
	Difficulties []Difficulty `json:"difficulties"`
}

type Partition struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	CompactName   string `json:"compactName"`
	DefaultStatus bool   `json:"default"`
}

type Difficulty struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Sizes []int  `json:"sizes"`
}

func GetExpansions() ([]Expansion, error) {
	accessToken, err := getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %v", err)
	}

	client := graphql.NewClient(ApiUrl)

	req := graphql.NewRequest(`
        query {
            worldData {
                expansions {
                    id
					name
					zones {
						id
						difficulties {
							id
							name
							sizes
						}
						name
						partitions {
							id
							name
							compactName
							default
						}
					}
                }
            }
        }
    `)

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	var resp ExpansionsResponse
	ctx := context.Background()
	if err := client.Run(ctx, req, &resp); err != nil {
		log.Printf("GraphQL Request GetExpansions Error: %v", err)
		return nil, fmt.Errorf("error querying regions: %w", err)
	}

	return resp.WorldData.Expansions, nil
}

func GetZoneIds() (*RegionsResponse, error) {
	accessToken, err := getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %v", err)
	}

	client := graphql.NewClient(ApiUrl)

	req := graphql.NewRequest(`
        query {
            worldData {
                zones {
                    id
                    slug
                }
            }
        }
    `)

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	var resp RegionsResponse
	ctx := context.Background()
	if err := client.Run(ctx, req, &resp); err != nil {
		log.Printf("GraphQL Request GetZoneIds Error: %v", err)
		return nil, fmt.Errorf("error querying zones: %w", err)
	}

	return &resp, nil
}

func GetRegions() ([]Region, error) {
	accessToken, err := getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %v", err)
	}

	client := graphql.NewClient(ApiUrl)

	req := graphql.NewRequest(`
        query {
            worldData {
                regions {
                    id
					compactName
					name
                    slug
                }
            }
        }
    `)

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	var resp RegionsResponse
	ctx := context.Background()
	if err := client.Run(ctx, req, &resp); err != nil {
		log.Printf("GraphQL Request GetRegions Error: %v", err)
		return nil, fmt.Errorf("error querying regions: %w", err)
	}

	if len(resp.WorldData.Regions) == 0 {
		return nil, fmt.Errorf("no regions returned")
	}

	return resp.WorldData.Regions, nil
}

func GetServersFromRegion(regionId int) ([]Server, error) {
	accessToken, err := getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %v", err)
	}

	client := graphql.NewClient(ApiUrl)

	compiledServers := []Server{}
	page := 1

	for {
		req := graphql.NewRequest(`
			query ServersQuery($regionID: Int!, $page: Int!) {
				worldData {
					region(id: $regionID) {
						servers(limit: 100, page: $page) {
							data {
								id
								name
								normalizedName
								slug
								region {
									slug
								}
							}
						}
					}
				}
			}
		`)

		req.Var("regionID", regionId)
		req.Var("page", page)

		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Content-Type", "application/json")

		var resp ServersResponse
		ctx := context.Background()
		if err := client.Run(ctx, req, &resp); err != nil {
			log.Printf("GraphQL Request GetServersFromRegion Error: %v", err)
			return nil, fmt.Errorf("failed to get query servers: %w", err)
		}

		compiledServers = append(compiledServers, resp.WorldData.Region.Servers.Data...)
		if page >= resp.WorldData.Region.Servers.LastPage {
			break
		}
		page++
	}

	return compiledServers, nil
}
