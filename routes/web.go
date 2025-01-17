package routes

import (
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
}
