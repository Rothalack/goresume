package config

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	DB                       *sql.DB
	WarcraftlogsClientId     string
	WarcraftlogsClientSecret string
)

func GetEnvs() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	WarcraftlogsClientId = os.Getenv("WARCRAFTLOGS_CLIENT_ID")
	WarcraftlogsClientSecret = os.Getenv("WARCRAFTLOGS_CLIENT_SECRET")
}

func GetEnvironment(router *gin.Engine) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	env := os.Getenv("GO_ENV")
	port := os.Getenv("PORT")
	debug := os.Getenv("DEBUG")

	initMySQL()

	if env == "prod" {
		prod(router, port)
	} else {
		dev(router, debug, port)
	}
}

func prod(router *gin.Engine, port string) {
	certificates := map[string]tls.Certificate{
		"brian-schaaf.com":     loadCertificate("./keys/brian-schaaf.com.pem", "./keys/brian-schaaf.com.key"),
		"www.brian-schaaf.com": loadCertificate("./keys/brian-schaaf.com.pem", "./keys/brian-schaaf.com.key"),
		"frothy.dev":           loadCertificate("./keys/frothy.dev.pem", "./keys/frothy.dev.key"),
		"www.frothy.dev":       loadCertificate("./keys/frothy.dev.pem", "./keys/frothy.dev.key"),
	}

	tlsConfig := &tls.Config{
		GetCertificate: func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			cert, ok := certificates[clientHello.ServerName]
			if !ok {
				log.Printf("No certificate found for domain: %s", clientHello.ServerName)
				return nil, nil
			}
			return &cert, nil
		},
	}

	server := &http.Server{
		Addr:      fmt.Sprintf(":%s", port),
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	log.Println("Starting prod server...")
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func dev(router *gin.Engine, debug string, port string) {
	log.Println("Starting dev server...")
	if debug == "true" {
		router.Use(gin.Logger())
		router.Use(gin.Recovery())
	}
	router.Run(fmt.Sprintf(":%s", port))
}

func loadCertificate(certFile, keyFile string) tls.Certificate {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load certificate: %v", err)
	}
	return cert
}

func initMySQL() {
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")
	dbPort := os.Getenv("MYSQL_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(localhost:%s)/%s?parseTime=true", dbUser, dbPass, dbPort, dbName)

	var dbErr error
	DB, dbErr = sql.Open("mysql", dsn)
	if dbErr != nil {
		log.Fatalf("Error connecting to database: %v", dbErr)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	DB.SetMaxOpenConns(50)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(time.Minute * 5)
}
