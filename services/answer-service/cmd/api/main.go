package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	http_adapter "github.com/question-interviewer/answer-service/internal/adapters/http"
	"github.com/question-interviewer/answer-service/internal/adapters/postgres"
	"github.com/question-interviewer/answer-service/internal/services"
)

func main() {
	log.Println("Starting Answer Service...")

	// Database Connection
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Default fallback for local dev
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbUser == "" {
		dbUser = "user"
	}
	if dbPassword == "" {
		dbPassword = "password"
	}
	if dbName == "" {
		dbName = "question_db"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Printf("Warning: Failed to ping database: %v", err)
	} else {
		log.Println("Connected to Database")
	}

	// Dependency Injection
	repo := postgres.NewAnswerRepository(db)
	svc := services.NewAnswerService(repo)
	handler := http_adapter.NewAnswerHandler(svc)

	// Router Setup
	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	handler.RegisterRoutes(r)

	// Start Server
	log.Println("Server listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
