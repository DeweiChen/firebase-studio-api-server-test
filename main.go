package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Print("starting server...")
	r := gin.Default()
	r.GET("/", handler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
		log.Printf("Defaulting to port %s", port)
	}

	r.GET("/health", healthHandler)
	// Start HTTP server.
	log.Printf("Listening on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func handler(c *gin.Context) {
	c.String(200, "Hello, World!")
}

func healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}