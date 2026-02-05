package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
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
	configJSON, err := json.Marshal(session.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal session config: %w", err)
	}

	query := `
		INSERT INTO practice_sessions (id, user_id, score, started_at, ended_at, status, topic_id, level, language, config)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err = r.db.ExecContext(ctx, query,
		session.ID,
		session.UserID,
		session.Score,
		session.StartedAt,
		session.EndedAt,
		session.Status,
		session.TopicID,
		session.Level,
		session.Language,
		configJSON,
	)
	if err != nil {
		return fmt.Errorf("failed to create practice session: %w", err)
	}
	return nil
}

func (r *PracticeRepository) GetSession(ctx context.Context, id uuid.UUID) (*domain.PracticeSession, error) {
	query := `
		SELECT id, user_id, score, started_at, ended_at, status, topic_id, level, language, config
		FROM practice_sessions
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var s domain.PracticeSession
	var configJSON []byte

	err := row.Scan(
		&s.ID,
		&s.UserID,
		&s.Score,
		&s.StartedAt,
		&s.EndedAt,
		&s.Status,
		&s.TopicID,
		&s.Level,
		&s.Language,
		&configJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if configJSON != nil {
		if err := json.Unmarshal(configJSON, &s.Config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal session config: %w", err)
		}
	} else {
		s.Config = make(map[string]interface{})
	}

	return &s, nil
}

func (r *PracticeRepository) UpdateSession(ctx context.Context, session *domain.PracticeSession) error {
	configJSON, err := json.Marshal(session.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal session config: %w", err)
	}

	query := `
		UPDATE practice_sessions
		SET score = $1, ended_at = $2, status = $3, config = $4
		WHERE id = $5
	`
	_, err = r.db.ExecContext(ctx, query,
		session.Score,
		session.EndedAt,
		session.Status,
		configJSON,
		session.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}
	return nil
}

func (r *PracticeRepository) GetRandomQuestionID(ctx context.Context, topicID *uuid.UUID, level *string, language string, config map[string]interface{}) (uuid.UUID, error) {
	// Base query
	query := `SELECT q.id FROM questions q`
	whereClauses := []string{"q.status = 'published'"}
	args := []interface{}{}
	argIdx := 1

	// 1. Topic ID (Direct filter)
	if topicID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("q.topic_id = $%d", argIdx))
		args = append(args, *topicID)
		argIdx++
	}

	// 2. Stacks/Roles from Config (Topic Name filter)
	// If topicID is nil, check if we have stacks in config
	if topicID == nil && config != nil {
		// Check for "tech_stacks" (frontend key) or "stacks" (legacy/fallback)
		var stacksInterface interface{}
		if s, ok := config["tech_stacks"]; ok {
			stacksInterface = s
		} else if s, ok := config["stacks"]; ok {
			stacksInterface = s
		}

		if stacks, ok := stacksInterface.([]interface{}); ok && len(stacks) > 0 {
			stackNames := make([]string, len(stacks))
			for i, s := range stacks {
				raw := fmt.Sprint(s)
				// Normalize stack names to match DB topics
				switch raw {
				case "Go":
					stackNames[i] = "Golang"
				case "Node.js":
					stackNames[i] = "NodeJS"
				case "PostgreSQL", "MongoDB", "Redis":
					stackNames[i] = "Data Layer"
				case "TypeScript":
					stackNames[i] = "JavaScript"
				default:
					stackNames[i] = raw
				}
			}

			// Join with topics table if we are filtering by topic name
			query += ` JOIN topics t ON q.topic_id = t.id`

			whereClauses = append(whereClauses, fmt.Sprintf("t.name = ANY($%d)", argIdx))
			// ANY needs a pq array compatible format or just passed as slice if driver supports it.
			// pgx stdlib usually handles []string as text array.
			args = append(args, stackNames)
			argIdx++
		}
	}

	// 3. Role from Config
	if config != nil {
		if role, ok := config["role"].(string); ok && role != "" {
			roundID, _ := config["round_id"].(string)

			if roundID == "devops_round" {
				whereClauses = append(whereClauses, fmt.Sprintf("(q.role = $%d OR q.role = 'Any' OR q.role = 'DevOps')", argIdx))
				args = append(args, role)
				argIdx++
			} else {
				whereClauses = append(whereClauses, fmt.Sprintf("(q.role = $%d OR q.role = 'Any')", argIdx))
				args = append(args, role)
				argIdx++
			}
		}
	}

	// 4. Language
	targetLang := "en"
	if language != "" {
		targetLang = language
	}
	whereClauses = append(whereClauses, fmt.Sprintf("q.language = $%d", argIdx))
	args = append(args, targetLang)
	argIdx++

	// Build WHERE string
	whereStr := " WHERE " + whereClauses[0]
	for i := 1; i < len(whereClauses); i++ {
		whereStr += " AND " + whereClauses[i]
	}

	// 4. Level Logic (With progression, then fallback)
	if level != nil && *level != "" {
		targetLevels := []string{*level}
		switch *level {
		case "Fresher":
			targetLevels = append(targetLevels, "Junior")
		case "Junior":
			targetLevels = append(targetLevels, "Mid")
		case "Mid":
			targetLevels = append(targetLevels, "Senior")
		}

		levelClause := fmt.Sprintf(" AND (q.level = ANY($%d) OR q.level = 'Any')", argIdx)
		levelQuery := query + whereStr + levelClause + " ORDER BY RANDOM() LIMIT 1"
		levelArgs := append(args, targetLevels)

		var checkID uuid.UUID
		err := r.db.QueryRowContext(ctx, levelQuery, levelArgs...).Scan(&checkID)
		if err == nil {
			return checkID, nil
		}
		// If not found, fall through to fallback
	}

	// Fallback: Ignore level
	finalQuery := query + whereStr + " ORDER BY RANDOM() LIMIT 1"
	finalArgs := args

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

func (r *PracticeRepository) GetQuestionSampleCache(ctx context.Context, questionID uuid.UUID) (string, string, []string, string, error) {
	query := `
		SELECT COALESCE(sample_answer, ''), COALESCE(sample_feedback, ''), sample_suggestions, COALESCE(sample_source, '')
		FROM questions
		WHERE id = $1
	`

	var sampleAnswer, sampleFeedback, sampleSource string
	var suggestionsRaw []byte

	row := r.db.QueryRowContext(ctx, query, questionID)
	if err := row.Scan(&sampleAnswer, &sampleFeedback, &suggestionsRaw, &sampleSource); err != nil {
		return "", "", nil, "", fmt.Errorf("failed to get question sample cache: %w", err)
	}

	var sampleSuggestions []string
	if len(suggestionsRaw) > 0 {
		_ = json.Unmarshal(suggestionsRaw, &sampleSuggestions)
	}

	return sampleAnswer, sampleFeedback, sampleSuggestions, sampleSource, nil
}

func (r *PracticeRepository) UpsertQuestionSampleCache(ctx context.Context, questionID uuid.UUID, sampleAnswer, sampleFeedback string, sampleSuggestions []string, sampleSource string) error {
	var suggestionsJSON []byte
	var err error
	if sampleSuggestions != nil {
		suggestionsJSON, err = json.Marshal(sampleSuggestions)
		if err != nil {
			return fmt.Errorf("failed to marshal sample suggestions: %w", err)
		}
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE questions
		SET sample_answer = $2,
			sample_feedback = $3,
			sample_suggestions = $4,
			sample_source = $5,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, questionID, sampleAnswer, sampleFeedback, suggestionsJSON, sampleSource)
	if err != nil {
		return fmt.Errorf("failed to upsert question sample cache: %w", err)
	}
	return nil
}

func (r *PracticeRepository) GetQuestionContent(ctx context.Context, questionID uuid.UUID) (string, string, string, string, string, error) {
	// Join with topics table to get topic name if needed, but for now assuming we just need question fields
	// But wait, topic name is in topics table.
	// Let's assume questions table has content and level.
	// And we need topic name.

	query := `
		SELECT q.content, t.name, q.level, COALESCE(q.correct_answer, ''), COALESCE(q.hint, '')
		FROM questions q
		LEFT JOIN topics t ON q.topic_id = t.id
		WHERE q.id = $1
	`
	var content, topic, level, correctAnswer, hint string
	// Handle potential NULLs if topic is missing
	var topicName sql.NullString

	row := r.db.QueryRowContext(ctx, query, questionID)
	if err := row.Scan(&content, &topicName, &level, &correctAnswer, &hint); err != nil {
		return "", "", "", "", "", fmt.Errorf("failed to get question content: %w", err)
	}

	if topicName.Valid {
		topic = topicName.String
	} else {
		topic = "General"
	}

	return content, topic, level, correctAnswer, hint, nil
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
		`INSERT INTO questions (id, topic_id, content, level, correct_answer, sample_answer, sample_source, hint, title) 
		 VALUES ($1, $2, $3, $4, $5, $6, 'user', $7, $8)`,
		q.ID, topicID, q.Content, q.Level, q.CorrectAnswer, q.CorrectAnswer, q.Hint, "Generated Question")

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
