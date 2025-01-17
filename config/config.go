package config

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Prod(router *gin.Engine, port string) {
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

func Dev(router *gin.Engine, debug string, port string) {
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
