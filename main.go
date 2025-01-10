package main

import (
	"net/http"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.HTMLRender = ginview.Default()

	router.Static("/static", "./static")
	router.StaticFile("/favicon.ico", "./static/images/favicon.ico")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "home", gin.H{
			"title": "Brian Schaaf",
			"add": func(a int, b int) int {
				return a + b
			},
		})
	})

	router.GET("/resume", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "resume", gin.H{
			"title": "Resume",
		})
	})

	router.GET("/gamin", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "gamin", gin.H{
			"title": "Gamin",
		})
	})

	router.Run(":8080")
}
