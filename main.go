package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

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

	router.GET("/gohard", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "gohard", gin.H{
			"title": "Go Hard",
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
		Addr:      ":4443",
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	log.Println("Starting server on :4443...")
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func loadCertificate(certFile, keyFile string) tls.Certificate {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load certificate: %v", err)
	}
	return cert
}
