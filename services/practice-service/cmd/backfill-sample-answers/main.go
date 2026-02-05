package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/question-interviewer/practice-service/internal/adapters/ai"
	"github.com/question-interviewer/practice-service/internal/adapters/postgres"
	"github.com/question-interviewer/practice-service/internal/services"
)

func main() {
	dbHost := getenvDefault("DB_HOST", "localhost")
	dbPort := getenvDefault("DB_PORT", "5432")
	dbUser := getenvDefault("DB_USER", "user")
	dbPassword := getenvDefault("DB_PASSWORD", "password")
	dbName := getenvDefault("DB_NAME", "question_db")
	aiServiceURL := getenvDefault("AI_SERVICE_URL", "http://localhost:8000")
	language := getenvDefault("LANGUAGE", "vi")
	limit := getenvIntDefault("LIMIT", 50)
	sleepMS := getenvIntDefault("SLEEP_MS", 250)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	repo := postgres.NewPracticeRepository(db)
	aiClient := ai.NewAIClient(aiServiceURL)
	svc := services.NewPracticeService(repo, aiClient, true)

	ctx := context.Background()

	rows, err := db.QueryContext(ctx, `
		SELECT id
		FROM questions
		WHERE COALESCE(sample_source, '') <> 'ai'
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		log.Fatalf("Failed to query questions: %v", err)
	}
	defer rows.Close()

	ids := make([]uuid.UUID, 0, limit)
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			log.Fatalf("Failed to scan id: %v", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Failed to iterate ids: %v", err)
	}

	for i, id := range ids {
		_, _, _, improved, err := svc.SuggestAnswer(ctx, id, "", language)
		if err != nil {
			log.Printf("Backfill failed for %s: %v", id, err)
		} else if strings.TrimSpace(improved) == "" {
			log.Printf("Backfill empty sample for %s", id)
		} else {
			log.Printf("Backfilled %d/%d: %s", i+1, len(ids), id)
		}

		if sleepMS > 0 {
			time.Sleep(time.Duration(sleepMS) * time.Millisecond)
		}
	}
}

func getenvDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func getenvIntDefault(key string, def int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

