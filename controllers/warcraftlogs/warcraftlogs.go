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
	"time"

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
	Id      int         `json:"id"`
	Name    string      `json:"name"`
	Faction GameFaction `json:"faction"`
}

type Rank struct {
	Number     int    `json:"number"`
	Percentile *int   `json:"percentile"`
	Color      string `json:"color"`
}

type WorldRegionServerRankPositions struct {
	WorldRank  Rank `json:"worldRank"`
	RegionRank Rank `json:"regionRank"`
	ServerRank Rank `json:"serverRank"`
}

type GuildZoneRankings struct {
	GuildData struct {
		Guild struct {
			ZoneRanking struct {
				Progress          WorldRegionServerRankPositions `json:"progress"`
				Speed             WorldRegionServerRankPositions `json:"speed"`
				CompleteRaidSpeed WorldRegionServerRankPositions `json:"completeRaidSpeed"`
			} `json:"zoneRanking"`
		} `json:"guild"`
	} `json:"guildData"`
}

type RankingRequest struct {
	GuildName    string `form:"guild" binding:"required"`
	ApiUrl       string `form:"api_url" binding:"required"`
	ServerSlug   string `form:"server" binding:"required"`
	RegionSlug   string `form:"region" binding:"required"`
	ZoneId       int    `form:"zone" binding:"required"`
	DifficultyId int    `form:"difficulty"`
	Size         int    `form:"size"`
}

type CharRequest struct {
	GuildId int `form:"guild" binding:"required"`
	ZoneId  int `form:"zone" binding:"required"`
}

type GuildCharactersResponse struct {
	CharacterData struct {
		Characters struct {
			Data     []Character `json:"data"`
			LastPage int         `json:"last_page"`
		} `json:"characters"`
	} `json:"characterData"`
}

type GuildCharacters struct {
	GuildId    int `json:"guildId"`
	Difficulty int `json:"difficulty"`
	Size       int `json:"size"`
	ZoneId     int `json:"zoneId"`
	Limit      int `json:"limit"`
	Page       int `json:"page"`
}

type Character struct {
	Id      int         `json:"id"`
	Name    string      `json:"name"`
	Faction GameFaction `json:"faction"`
	ClassId int         `json:"classId"`
	Level   int         `json:"level"`
}

type CharRankings struct {
	Ranks        []Rank  `json:"ranks"`
	BestAmount   float64 `json:"bestAmount"`
	BestSpec     string  `json:"bestSpec"`
	MedianAmount float64 `json:"medianAmount"`
}

type CharRank struct {
	Rank       int     `json:"rank"`
	OutOf      int     `json:"outOf"`
	Amount     float64 `json:"amount"`
	Spec       string  `json:"spec"`
	Percentile float64 `json:"percentile"`
}

type GameFaction struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func SetGameContext(baseUrl string) {
	tokenUrl = baseUrl + "oauth/token"
	ApiUrl = baseUrl + "api/v2/client"
}

func getAccessToken() (string, error) {
	redisKey := config.RedisKeyPrefix + "warcraftlogs_access_token"

	token, err := config.RedisClient.Get(redisKey).Result()
	if err == nil && token != "" {
		fmt.Println("cached warcraftlogs access token used")
		return token, nil
	}

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

	if err := config.RedisClient.Set(redisKey, tokenResponse.AccessToken, time.Hour).Err(); err != nil {
		return "", fmt.Errorf("failed to set token in Redis: %w", err)
	}

	return tokenResponse.AccessToken, nil
}

func GetRanking(guildData RankingRequest) (*GuildZoneRankings, int, string, error) {
	guildName := guildData.GuildName
	serverSlug := guildData.ServerSlug
	regionSlug := guildData.RegionSlug
	baseUrl := guildData.ApiUrl
	zoneId := guildData.ZoneId
	difficultyId := guildData.DifficultyId
	size := guildData.Size

	SetGameContext(baseUrl)

	guild, err := GetGuild(guildName, regionSlug, serverSlug)
	guildId := guild.GuildData.Guild.Id
	guildFaction := guild.GuildData.Guild.Faction.Name
	if err != nil {
		return nil, 0, "", fmt.Errorf("failed to find guild: %w", err)
	}
	if guildId < 1 {
		return nil, 0, "", fmt.Errorf("failed to find guild")
	}

	accessToken, err := getAccessToken()
	if err != nil {
		return nil, 0, "", fmt.Errorf("failed to get access token: %w", err)
	}

	client := graphql.NewClient(ApiUrl)

	req := graphql.NewRequest(`
	    query GuildZoneRanking($guildId: Int!, $zoneId: Int!, $size: Int, $difficulty: Int) {
			guildData {
				guild(id: $guildId) {
					zoneRanking(zoneId: $zoneId) {
						progress(size: $size) {
							worldRank {
								number
								percentile
								color
							}
							regionRank {
								number
								percentile
								color
							}
							serverRank {
								number
								percentile
								color
							}
						}
						speed(size: $size, difficulty: $difficulty) {
							worldRank {
								number
								percentile
								color
							}
							regionRank {
								number
								percentile
								color
							}
							serverRank {
								number
								percentile
								color
							}
						}
						completeRaidSpeed(size: $size, difficulty: $difficulty) {
							worldRank {
								number
								percentile
								color
							}
							regionRank {
								number
								percentile
								color
							}
							serverRank {
								number
								percentile
								color
							}
						}
					}
				}
			}
		}
	`)

	req.Var("guildId", guildId)
	req.Var("zoneId", zoneId)
	req.Var("size", size)
	req.Var("difficulty", difficultyId)

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	var resp GuildZoneRankings
	ctx := context.Background()
	if err := client.Run(ctx, req, &resp); err != nil {
		log.Printf("GraphQL Request Error: %v", err)
		return nil, 0, "", fmt.Errorf("error querying guild: %w", err)
	}

	// prettyJSON, _ := json.MarshalIndent(resp, "", "  ")
	// log.Printf("Full Raw Response: %s", string(prettyJSON))

	return &resp, guildId, guildFaction, nil
}

func GetGuild(guildName, guildRegion, guildServer string) (*GuildResponse, error) {
	accessToken, err := getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	client := graphql.NewClient(ApiUrl)

	req := graphql.NewRequest(`
        query GuildQuery($name: String!, $serverSlug: String!, $serverRegion: String!) {
            guildData {
                guild(name: $name, serverSlug: $serverSlug, serverRegion: $serverRegion) {
                    id
                    name
					faction {
						id
						name
					}
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

	// prettyJSON, _ := json.MarshalIndent(resp, "", "  ")
	// log.Printf("Full Raw Response: %s", string(prettyJSON))

	return &resp, nil
}

func GetChars(requestData CharRequest) ([]Character, error) {
	guildId := requestData.GuildId

	accessToken, err := getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	fmt.Println(ApiUrl)
	client := graphql.NewClient(ApiUrl)

	compiledCharacters := []Character{}
	page := 1

	for {
		req := graphql.NewRequest(`
			query GuildCharacters($guildId: Int!, $page: Int!) {
				characterData {
					characters(guildID: $guildId, limit: 100, page: $page) {
						data {
							id
							name
							classID
							level
						}
						last_page
					}
				}
			}
		`)

		req.Var("guildId", guildId)
		req.Var("page", page)

		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Content-Type", "application/json")

		var resp GuildCharactersResponse
		ctx := context.Background()
		if err := client.Run(ctx, req, &resp); err != nil {
			log.Printf("GraphQL Request Error: %v", err)
			return nil, fmt.Errorf("error querying guild: %w", err)
		}

		// prettyJSON, _ := json.MarshalIndent(resp, "", "  ")
		// log.Printf("Full Raw Response: %s", string(prettyJSON))

		compiledCharacters = append(compiledCharacters, resp.CharacterData.Characters.Data...)
		if page >= resp.CharacterData.Characters.LastPage {
			break
		}
		page++
	}

	return compiledCharacters, nil
}
