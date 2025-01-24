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
	tokenUrl string
	ApiUrl   string
)

type GuildResponse struct {
	GuildData struct {
		Guild Guild `json:"guild"`
	} `json:"guildData"`
}

type Guild struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Faction     struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"faction"`
	Server struct {
		Name   string `json:"name"`
		Region struct {
			Id          int    `json:"id"`
			CompactName string `json:"compactName"`
			Name        string `json:"name"`
			Slug        string `json:"slug"`
		} `json:"region"`
		Slug string `json:"slug"`
	} `json:"server"`
	CompetitionMode bool `json:"competitionMode"`
	StealthMode     bool `json:"stealthMode"`
}

type GuildZoneRankings struct {
	Progress struct {
		World  int `json:"world"`
		Region int `json:"region"`
		Server int `json:"server"`
	} `json:"progress"`
	Speed struct {
		World  int `json:"world"`
		Region int `json:"region"`
		Server int `json:"server"`
	} `json:"speed"`
	CompleteRaidSpeed struct {
		World  int `json:"world"`
		Region int `json:"region"`
		Server int `json:"server"`
	} `json:"completeRaidSpeed"`
}

func SetGameContext(baseUrl string) {
	tokenUrl = baseUrl + "oauth/token"
	ApiUrl = baseUrl + "api/v2/client"
}

func getAccessToken() (string, error) {
	clientID := config.WarcraftlogsClientId
	clientSecret := config.WarcraftlogsClientSecret

	req, err := http.NewRequest("POST", tokenUrl, nil)
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

func GetGuild(guildName, guildRegion, guildServer string) (*GuildResponse, error) {
	accessToken, err := getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %v", err)
	}

	client := graphql.NewClient(ApiUrl)

	req := graphql.NewRequest(`
        query GuildQuery($name: String!, $serverSlug: String!, $serverRegion: String!) {
            guildData {
                guild(name: $name, serverSlug: $serverSlug, serverRegion: $serverRegion) {
                    id
                    name
                    description
                    faction {
						id
						name
					}
                    server {
                        name
                        region {
							id
							compactName
							name
							slug
						}
                        slug
                    }
                    competitionMode
                    stealthMode
                }
            }
        }
    `)

	req.Var("name", guildName)
	req.Var("serverSlug", guildServer)
	req.Var("serverRegion", guildRegion)

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	var resp GuildResponse
	ctx := context.Background()
	if err := client.Run(ctx, req, &resp); err != nil {
		log.Printf("GraphQL Request Error: %v", err)
		return nil, fmt.Errorf("error querying guild: %w", err)
	}

	prettyJSON, _ := json.MarshalIndent(resp, "", "  ")
	log.Printf("Full Raw Response: %s", string(prettyJSON))

	return &resp, nil
}

func GetGuildZoneRanking(guildId, zoneId int) (*GuildZoneRankings, error) {
	accessToken, err := getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %v", err)
	}

	client := graphql.NewClient(ApiUrl)

	req := graphql.NewRequest(`
        query GuildZoneRanking($guildId: Int!, $zoneId: Int) {
			guildData {
				guild(id: $guildId) {
					zoneRanking(zoneId: $zoneId) {
						progress(size: 20) {
							world
							region
							server
						}
						speed(size: 20, difficulty: 4) {
							world
							region
							server
						}
						completeRaidSpeed(size: 20, difficulty: 4) {
							world
							region
							server
						}
					}
				}
			}
		}
    `)

	req.Var("guildId", guildId)
	req.Var("zoneId", zoneId)

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	var resp GuildZoneRankings
	ctx := context.Background()
	if err := client.Run(ctx, req, &resp); err != nil {
		log.Printf("GraphQL Request Error: %v", err)
		return nil, fmt.Errorf("error querying guild: %w", err)
	}

	prettyJSON, _ := json.MarshalIndent(resp, "", "  ")
	log.Printf("Full Raw Response: %s", string(prettyJSON))

	return &resp, nil
}

// var resp map[string]interface{}
// if err := client.Run(context.Background(), req, &resp); err != nil {
// 	log.Printf("querying guild: %v", err)
// }

// prettyJSON, err := json.MarshalIndent(resp, "", "    ")
// if err != nil {
// 	log.Printf("Error pretty printing JSON: %v", err)
// } else {
// 	fmt.Printf("\n=== Query Response ===\n%s\n", string(prettyJSON))
// }

// prettyJSON, _ := json.MarshalIndent(resp, "", "  ")
// log.Printf("Full Raw Response: %s", string(prettyJSON))

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
