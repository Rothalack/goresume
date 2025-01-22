package main

import (
	"goresume/config"
	"goresume/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.GetEnvs()

	router := gin.Default()

	routes.Routes(router)

	config.GetEnvironment(router)
}
