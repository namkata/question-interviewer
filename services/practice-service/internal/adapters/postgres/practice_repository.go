package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/question-interviewer/practice-service/internal/domain"
	"github.com/question-interviewer/practice-service/internal/ports"
)

type PracticeRepository struct {
	db *sql.DB
}

func NewPracticeRepository(db *sql.DB) ports.PracticeRepository {
	return &PracticeRepository{
		db: db,
	}
}

func (r *PracticeRepository) CreateSession(ctx context.Context, session *domain.PracticeSession) error {
	query := `
		INSERT INTO practice_sessions (id, user_id, score, started_at, ended_at, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		session.ID,
		session.UserID,
		session.Score,
		session.StartedAt,
		session.EndedAt,
		session.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to create practice session: %w", err)
	}
	return nil
}

func (r *PracticeRepository) GetSession(ctx context.Context, id uuid.UUID) (*domain.PracticeSession, error) {
	query := `
		SELECT id, user_id, score, started_at, ended_at, status
		FROM practice_sessions
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var s domain.PracticeSession
	err := row.Scan(
		&s.ID,
		&s.UserID,
		&s.Score,
		&s.StartedAt,
		&s.EndedAt,
		&s.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return &s, nil
}

func (r *PracticeRepository) UpdateSession(ctx context.Context, session *domain.PracticeSession) error {
	query := `
		UPDATE practice_sessions
		SET score = $1, ended_at = $2, status = $3
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query,
		session.Score,
		session.EndedAt,
		session.Status,
		session.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}
	return nil
}

func (r *PracticeRepository) GetRandomQuestionID(ctx context.Context, topicID *uuid.UUID, level *string) (uuid.UUID, error) {
	// If topic/level not specified or no questions found with specific level, try broader search
	query := `SELECT id FROM questions WHERE status = 'published'`
	args := []interface{}{}
	argIdx := 1

	if topicID != nil {
		query += fmt.Sprintf(" AND topic_id = $%d", argIdx)
		args = append(args, *topicID)
		argIdx++
	}

	// Try with level first
	var finalQuery string
	var finalArgs []interface{}

	if level != nil && *level != "" {
		levelQuery := query + fmt.Sprintf(" AND level = $%d", argIdx)
		levelArgs := append(args, *level)
		finalQuery = levelQuery + " ORDER BY RANDOM() LIMIT 1"
		finalArgs = levelArgs

		// Check if any exist
		var checkID uuid.UUID
		err := r.db.QueryRowContext(ctx, finalQuery, finalArgs...).Scan(&checkID)
		if err == nil {
			return checkID, nil
		}
	}

	// If we are here, either level was nil or no questions found for that level.
	// Fallback: Ignore level, just match topic
	finalQuery = query + " ORDER BY RANDOM() LIMIT 1"
	finalArgs = args

	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, finalQuery, finalArgs...).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get random question: %w", err)
	}
	return id, nil
}

func (r *PracticeRepository) CreateAttempt(ctx context.Context, attempt *domain.PracticeAttempt) error {
	query := `
		INSERT INTO practice_attempts (id, session_id, question_id, user_answer, score, feedback, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		attempt.ID,
		attempt.SessionID,
		attempt.QuestionID,
		attempt.UserAnswer,
		attempt.Score,
		attempt.Feedback,
		attempt.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create attempt: %w", err)
	}
	return nil
}

func (r *PracticeRepository) GetQuestionContent(ctx context.Context, questionID uuid.UUID) (string, string, string, string, error) {
	// Join with topics table to get topic name if needed, but for now assuming we just need question fields
	// But wait, topic name is in topics table.
	// Let's assume questions table has content and level.
	// And we need topic name.

	query := `
		SELECT q.content, t.name, q.level, COALESCE(q.correct_answer, '')
		FROM questions q
		LEFT JOIN topics t ON q.topic_id = t.id
		WHERE q.id = $1
	`
	var content, topic, level, correctAnswer string
	// Handle potential NULLs if topic is missing
	var topicName sql.NullString

	row := r.db.QueryRowContext(ctx, query, questionID)
	if err := row.Scan(&content, &topicName, &level, &correctAnswer); err != nil {
		return "", "", "", "", fmt.Errorf("failed to get question content: %w", err)
	}

	if topicName.Valid {
		topic = topicName.String
	} else {
		topic = "General"
	}

	return content, topic, level, correctAnswer, nil
}

func (r *PracticeRepository) CreateQuestion(ctx context.Context, q *domain.Question) error {
	// First ensure topic exists or get it (simplified: just use a default topic if not found or create one)
	// For MVP, let's assume we look up topic by name or insert it.

	// Check if topic exists
	var topicID uuid.UUID
	err := r.db.QueryRowContext(ctx, "SELECT id FROM topics WHERE name = $1", q.TopicName).Scan(&topicID)
	if err != nil {
		// Create topic if not exists
		topicID = uuid.New()
		_, err = r.db.ExecContext(ctx, "INSERT INTO topics (id, name, description) VALUES ($1, $2, $3)",
			topicID, q.TopicName, "Auto-generated topic")
		if err != nil {
			return fmt.Errorf("failed to create topic: %w", err)
		}
	}

	// Use content excerpt as title if not provided (though we don't have title in domain yet)
	// Since we made title nullable in migration 000006, we can skip it or provide a default.
	// We will insert title as first 50 chars of content for backward compatibility or just leave it null if allowed.
	// But let's check schema: migration 000006 makes title nullable.

	_, err = r.db.ExecContext(ctx,
		`INSERT INTO questions (id, topic_id, content, level, correct_answer, title) 
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		q.ID, topicID, q.Content, q.Level, q.CorrectAnswer, "Generated Question")

	if err != nil {
		return fmt.Errorf("failed to insert question: %w", err)
	}

	return nil
}

func (r *PracticeRepository) GetTopicIDByName(ctx context.Context, name string) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, "SELECT id FROM topics WHERE name = $1", name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, fmt.Errorf("topic not found: %s", name)
		}
		return uuid.Nil, fmt.Errorf("failed to get topic id: %w", err)
	}
	return id, nil
}
