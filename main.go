package main

import (
	"goresume/config"
	"goresume/controllers/commands"
	"goresume/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	config.GetEnvs()

	router := gin.Default()

	router.Use(config.PrometheusMiddleware())
	config.SetupPrometheus(router)

	routes.Routes(router)

	config.GetEnvironment(router)

	RunCommands()
}

func RunCommands() {
	command := os.Args[1]

	switch command {
	case "sync_base_data":
		commands.SyncBaseDataCommand()
	default:
		log.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
