package main

import (
	"goresume/config"
	"goresume/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	env := os.Getenv("GO_ENV")
	port := os.Getenv("PORT")
	debug := os.Getenv("DEBUG")

	router := gin.Default()

	routes.Routes(router)

	if env == "prod" {
		// Production setup (TLS, certificates, etc.)
		config.Prod(router, port)
	} else {
		// Development setup (debugging, logging)
		config.Dev(router, debug, port)
	}
}
