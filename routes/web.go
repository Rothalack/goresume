package routes

import (
	"goresume/controllers/warcraftlogs"
	"html/template"
	"net/http"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	var CustomConfig = goview.Config{
		Root:         "resources/views",
		Extension:    ".html",
		Master:       "layouts/master",
		Partials:     []string{},
		Funcs:        make(template.FuncMap),
		DisableCache: false,
		Delims:       goview.Delims{Left: "{{", Right: "}}"},
	}

	router.HTMLRender = ginview.New(CustomConfig)

	router.Static("/static", "./static")
	router.StaticFile("/favicon.ico", "./static/images/favicon.ico")
	router.StaticFile("/sitemap.xml", "./sitemap.xml")
	router.StaticFile("/security.txt", "./security.txt")
	router.StaticFile("/.well-known/security.txt", "./security.txt")
	router.StaticFile("/humans.txt", "./humans.txt")
	router.StaticFile("/ads.txt", "./ads.txt")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "home", gin.H{
			"title": "Brian Schaaf",
		})
	})

	router.GET("/resume", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "resume", gin.H{
			"title": "Resume",
		})
	})

	router.GET("/gohard", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "gohard", gin.H{
			"title": "Go Hard",
		})
	})

	router.GET("/rankings", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "rankings", gin.H{
			"title": "Rankings",
		})
	})

	router.GET("/api/logs-data", func(c *gin.Context) {
		resp, err := warcraftlogs.GetData()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": resp,
		})
	})

	router.GET("/api/ranking-data", func(c *gin.Context) {
		var req warcraftlogs.RankingRequest

		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		resp, guildId, err := warcraftlogs.GetRanking(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data":    resp.GuildData.Guild.ZoneRanking,
			"guildId": guildId,
		})
	})

	router.GET("/api/char-data", func(c *gin.Context) {
		var req warcraftlogs.CharRequest

		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		resp, err := warcraftlogs.GetChars(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": resp,
		})
	})

	// guildName := c.DefaultQuery("guildName", "")
	// gameId := c.DefaultQuery("game", "")
	// serverId := c.DefaultQuery("server", "")
	// regionId := c.DefaultQuery("region", "")
	// expansionId := c.DefaultQuery("expansion", "")
	// zoneId := c.DefaultQuery("zone", "")
	// difficultyId := c.DefaultQuery("difficulty", "")

	// router.GET("/tools", func(ctx *gin.Context) {
	// 	ctx.HTML(http.StatusOK, "tools", gin.H{
	// 		"title": "Tools",
	// 	})
	// })

	// router.GET("/gamin", func(ctx *gin.Context) {
	// 	ctx.HTML(http.StatusOK, "gamin", gin.H{
	// 		"title": "Gamin",
	// 	})
	// })

	// router.GET("/cars", func(ctx *gin.Context) {
	// 	ctx.HTML(http.StatusOK, "cars", gin.H{
	// 		"title": "Cars",
	// 	})
	// })

	// router.GET("/testvue", func(ctx *gin.Context) {
	// 	ctx.HTML(http.StatusOK, "testvue", gin.H{
	// 		"title": "testvue",
	// 	})
	// })

	// router.GET("/api/find-regions", func(c *gin.Context) {
	// 	resp, err := warcraftlogs.GetRegions()
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{
	// 			"error": err.Error(),
	// 		})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"data": resp.WorldData,
	// 	})
	// })

	// router.GET("/api/find-expansions", func(c *gin.Context) {
	// 	resp, err := warcraftlogs.GetExpansions()
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{
	// 			"error": err.Error(),
	// 		})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"data": resp,
	// 	})
	// })

	// router.GET("/api/find-server", func(c *gin.Context) {
	// 	regionIdStr := c.Query("regionId")

	// 	regionId, err := strconv.Atoi(regionIdStr)
	// 	if err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{
	// 			"error": "Invalid Region",
	// 		})
	// 		return
	// 	}

	// 	serversResp, err := warcraftlogs.GetServersFromRegion(regionId)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{
	// 			"error": err.Error(),
	// 		})
	// 		return
	// 	}

	// 	c.JSON(http.StatusOK, gin.H{
	// 		"data": serversResp.WorldData.Region.Servers.Data,
	// 	})
	// })

	// router.GET("/api/find-guild", func(c *gin.Context) {
	// 	guildName := c.Query("guild")
	// 	guildRegion := c.Query("guildRegion")
	// 	guildServer := c.Query("guildServer")

	// 	data, err := warcraftlogs.GetGuild(guildName, guildRegion, guildServer)

	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{
	// 			"error": err.Error(),
	// 		})
	// 		return
	// 	}

	// 	fmt.Println(data)

	// 	c.JSON(http.StatusOK, gin.H{
	// 		"guildName": guildName,
	// 		"region":    guildRegion,
	// 		"realm":     guildServer,
	// 	})
	// })

	// router.POST("/api/add-character", func(c *gin.Context) {
	// 	charName := c.PostForm("character")

	// 	// For now, return mock data

	// 	// data := warcraftlogs.GetCharacter()
	// 	// fmt.Println(data)
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"name":  charName,
	// 		"level": 60,
	// 		"class": "Mage",
	// 		"guild": charName,
	// 	})
	// })
}
