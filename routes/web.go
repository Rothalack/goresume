package routes

import (
	"goresume/controllers/warcraftlogs"
	"log"
	"net/http"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	router.HTMLRender = ginview.Default()

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

	router.GET("/tools", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "tools", gin.H{
			"title": "Tools",
		})
	})

	router.GET("/gamin", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "gamin", gin.H{
			"title": "Gamin",
		})
	})

	router.GET("/cars", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "cars", gin.H{
			"title": "Cars",
		})
	})

	router.GET("/rankings", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "rankings", gin.H{
			"title": "Rankings",
		})
	})

	router.GET("/api/find-regions", func(c *gin.Context) {
		data, err := warcraftlogs.GetRegions()
		if err != nil {
			log.Fatalf("Failed to get region data: %v", err)
		}
		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})

	router.GET("/api/find-server", func(c *gin.Context) {
		guildRegion := c.Query("regionId")

		warcraftlogs.GetServersFromRegion(guildRegion)
	})

	router.GET("/api/find-guild", func(c *gin.Context) {
		// guildName := c.Query("guild")
		// guildServer := c.Query("guildServer")

		// warcraftlogs.GetRegions()
		// warcraftlogs.GetServersFromRegion(6, 100, 1)
		// warcraftlogs.GetGuild()

		// guildName := c.Query("guild")
		// guildServer := c.Query("guildServer")

		// data, err := warcraftlogs.GetGuild(guildName, guildServer)
		// if err != nil {
		// 	log.Fatalf("Failed to get guild data: %v", err)
		// }
		// fmt.Println(data)
		// c.JSON(http.StatusOK, gin.H{
		// 	"guildName": guildName,
		// 	"region":    "US",
		// 	"realm":     "Stormrage",
		// })
	})

	router.POST("/api/add-character", func(c *gin.Context) {
		charName := c.PostForm("character")

		// For now, return mock data

		// data := warcraftlogs.GetCharacter()
		// fmt.Println(data)
		c.JSON(http.StatusOK, gin.H{
			"name":  charName,
			"level": 60,
			"class": "Mage",
			"guild": charName,
		})
	})
}
