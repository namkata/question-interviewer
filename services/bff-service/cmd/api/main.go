package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	bff_http "github.com/question-interviewer/bff-service/internal/adapters/http"
)

func main() {
	log.Println("Starting BFF Service...")

	// Configuration
	practiceServiceURL := os.Getenv("PRACTICE_SERVICE_URL")
	if practiceServiceURL == "" {
		practiceServiceURL = "http://localhost:8080"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Handlers
	handler := bff_http.NewBFFHandler(practiceServiceURL)

	// Router Setup
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // In production, replace with frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	handler.RegisterRoutes(r)

	// Start Server
	log.Printf("BFF Service listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
