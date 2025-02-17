package routes

import (
	"goresume/config/middleware"
	"goresume/controllers/auth"
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

	router.POST("/auth/register", auth.Register)
	router.GET("/register", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "auth/register", gin.H{
			"title": "Register",
		})
	})
	router.POST("/auth/login", auth.Login)
	router.POST("/auth/logout", auth.Logout)

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

	router.GET("/test", func(ctx *gin.Context) {
		ctx.File("resources/views/test.html")
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

		resp, guildId, guildFaction, err := warcraftlogs.GetRanking(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data":         resp.GuildData.Guild.ZoneRanking,
			"guildId":      guildId,
			"guildFaction": guildFaction,
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

	protected := router.Group("/admin")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/dashboard", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			c.JSON(200, gin.H{"user_id": userID})
		})
	}
}
