package warcraftlogs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"goresume/config"
	"io"
	"log"
	"net/http"

	"github.com/machinebox/graphql"
)

var (
	tokenURL = "https://www.warcraftlogs.com/oauth/token"
	apiURL   = "https://www.warcraftlogs.com/api/v2/client"
)

type GuildResponse struct {
	GuildData struct {
		Guild Guild `json:"guild"`
	} `json:"guildData"`
}

// Guild represents a WoW guild
type Guild struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	Server      Server            `json:"server"`
	Description string            `json:"description"`
	ZoneRanking GuildZoneRankings `json:"zoneRanking"`
}

// Server represents a WoW server
type Server struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Region string `json:"region"`
	Slug   string `json:"slug"`
}

// GuildZoneRankings represents rankings for a guild in a specific zone
type GuildZoneRankings struct {
	Progress json.RawMessage `json:"progress"` // Using RawMessage as the exact structure isn't provided
	Speed    json.RawMessage `json:"speed"`    // Using RawMessage as the exact structure isn't provided
}

// CharacterResponse wraps the character query response
type CharacterResponse struct {
	CharacterData struct {
		Character Character `json:"character"`
	} `json:"characterData"`
}

// Character represents a WoW character
type Character struct {
	ID           int             `json:"id"`
	Name         string          `json:"name"`
	Server       Server          `json:"server"`
	ClassID      int             `json:"classID"`
	Level        int             `json:"level"`
	ZoneRankings json.RawMessage `json:"zoneRankings"` // Using RawMessage as it returns JSON type
}

func getAccessToken() (string, error) {
	clientID := config.WarcraftlogsClientId
	clientSecret := config.WarcraftlogsClientSecret

	req, err := http.NewRequest("POST", tokenURL, nil)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqBody := "grant_type=client_credentials"
	req.Body = io.NopCloser(bytes.NewBufferString(reqBody))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status: %s, response: %s", resp.Status, body)
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

func GetCharacter() string {
	accessToken, err := getAccessToken()
	if err != nil {
		log.Fatalf("Failed to get access token test: %v", err)
	}

	return accessToken
}

// func GetGuild(guildName, guildServer string) (*Guild, error) {
// 	query := fmt.Sprintf(guildQuery, guildName, guildServer)

// 	requestBody := map[string]interface{}{
// 		"query": query,
// 	}
// 	body, err := json.Marshal(requestBody)
// 	if err != nil {
// 		return nil, err
// 	}

// 	accessToken, err := getAccessToken()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get access token: %w", err)
// 	}

// 	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Set("Authorization", "Bearer "+accessToken)
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
// 	}

// 	rawBody, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read response body: %w", err)
// 	}
// 	fmt.Println("Raw Response Body:", string(rawBody))

// 	var gqlResponse GuildResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&gqlResponse); err != nil {
// 		return nil, err
// 	}

// 	return &gqlResponse.Data.GuildData.Guild, nil
// }

func GetGuild() {
	accessToken, err := getAccessToken()
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	client := graphql.NewClient(apiURL)

	// req := graphql.NewRequest(`
	//     query GuildQuery($name: String!, $serverSlug: String!, $serverRegion: String!) {
	//         guildData {
	//             guild(name: $name, serverSlug: $serverSlug, serverRegion: $serverRegion) {
	//                 id
	//                 name
	//                 description
	//                 server {
	//                     id
	//                     name
	//                     region
	//                     slug
	//                 }
	//                 zoneRanking {
	//                     progress
	//                     speed
	//                 }
	//             }
	//         }
	//     }
	// `)

	req := graphql.NewRequest(`
	    query GuildQuery($id: Int!) {
	        guildData {
	            guild(id: $id) {
	                id
	                name
	                description
	                server {
	                    id
	                    name
	                    region
	                    slug
	                }
	                zoneRanking {
	                    progress
	                    speed
	                }
	            }
	        }
	    }
	`)

	req.Var("id", "488971")
	// req.Var("name", "Liquid")
	// req.Var("serverSlug", "Illidan")
	// req.Var("serverRegion", "US")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	var resp GuildResponse
	if err := client.Run(context.Background(), req, &resp); err != nil {
		// return nil, fmt.Errorf("querying guild: %w", err)
	}

	prettyJSON, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Printf("Error pretty printing JSON: %v", err)
	} else {
		fmt.Printf("\n=== Guild Query Response ===\n%s\n", string(prettyJSON))
	}

	// return &resp.GuildData.Guild, nil
}

func GetRegions() {
	accessToken, err := getAccessToken()
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	client := graphql.NewClient(apiURL)

	req := graphql.NewRequest(`
		query {
			worldData {
				regions {
					id
					name
					slug
				}
			}
		}
	`)

	req.Header.Set("Authorization", "Bearer "+accessToken)

	var resp map[string]interface{}
	if err := client.Run(context.Background(), req, &resp); err != nil {
		log.Printf("querying guild: %v", err)
	}

	prettyJSON, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Printf("Error pretty printing JSON: %v", err)
	} else {
		fmt.Printf("\n=== Query Response ===\n%s\n", string(prettyJSON))
	}
}

func GetServersFromRegion(regionID int, limit int, page int) {
	accessToken, err := getAccessToken()
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	client := graphql.NewClient(apiURL)

	req := graphql.NewRequest(`
		query ServersQuery($regionID: Int!, $limit: Int, $page: Int) {
			worldData {
				region(id: $regionID) {
					servers(limit: $limit, page: $page) {
						data {
							id
							name
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

	// Set variables for the query
	req.Var("regionID", regionID)
	req.Var("limit", limit)
	req.Var("page", page)

	req.Header.Set("Authorization", "Bearer "+accessToken)

	var resp map[string]interface{}
	if err := client.Run(context.Background(), req, &resp); err != nil {
		log.Printf("querying servers: %v", err)
	}

	prettyJSON, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Printf("Error pretty printing JSON: %v", err)
	} else {
		fmt.Printf("\n=== Servers Query Response ===\n%s\n", string(prettyJSON))
	}
}
