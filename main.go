package main

import (
	"log"
	"os"

	"github.com/dw/firebase-studio-api-server-test/task/handlers"
	"github.com/dw/firebase-studio-api-server-test/task/repositories"
	"github.com/gin-gonic/gin"
)

func setupRoutes(r *gin.Engine) {
	// Initialize task repository and handler
	taskRepo := repositories.NewInMemoryTaskRepository()
	taskHandler := handlers.NewTaskHandler(taskRepo)

	// Task routes
	r.POST("/tasks", taskHandler.CreateTask)
	r.GET("/tasks", taskHandler.GetAllTasks)
	r.GET("/tasks/:id", taskHandler.GetTask)
	r.PUT("/tasks/:id", taskHandler.UpdateTask)
	r.DELETE("/tasks/:id", taskHandler.DeleteTask)

	// Health check route
	r.GET("/health", healthHandler)
}

func main() {
	log.Print("starting server...")
	r := gin.Default()

	setupRoutes(r)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
		log.Printf("Defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("Listening on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
