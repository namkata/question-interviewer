package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/question-interviewer/question-service/internal/domain"
	"github.com/question-interviewer/question-service/internal/ports"
)

type QuestionRepository struct {
	db *sql.DB
}

func NewQuestionRepository(db *sql.DB) ports.QuestionRepository {
	return &QuestionRepository{
		db: db,
	}
}

func (r *QuestionRepository) Create(ctx context.Context, question *domain.Question) error {
	query := `
		INSERT INTO questions (id, title, content, level, topic_id, created_by, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		question.ID,
		question.Title,
		question.Content,
		question.Level,
		question.TopicID,
		question.CreatedBy,
		question.Status,
		question.CreatedAt,
		question.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}
	return nil
}

func (r *QuestionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Question, error) {
	query := `
		SELECT id, title, content, level, topic_id, created_by, status, created_at, updated_at
		FROM questions
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var q domain.Question
	err := row.Scan(
		&q.ID,
		&q.Title,
		&q.Content,
		&q.Level,
		&q.TopicID,
		&q.CreatedBy,
		&q.Status,
		&q.CreatedAt,
		&q.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("question not found")
		}
		return nil, fmt.Errorf("failed to get question: %w", err)
	}
	return &q, nil
}

func (r *QuestionRepository) List(ctx context.Context, limit, offset int) ([]*domain.Question, error) {
	query := `
		SELECT id, title, content, level, topic_id, created_by, status, created_at, updated_at
		FROM questions
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list questions: %w", err)
	}
	defer rows.Close()

	var questions []*domain.Question
	for rows.Next() {
		var q domain.Question
		err := rows.Scan(
			&q.ID,
			&q.Title,
			&q.Content,
			&q.Level,
			&q.TopicID,
			&q.CreatedBy,
			&q.Status,
			&q.CreatedAt,
			&q.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan question: %w", err)
		}
		questions = append(questions, &q)
	}
	return questions, nil
}

func (r *QuestionRepository) Update(ctx context.Context, question *domain.Question) error {
	query := `
		UPDATE questions
		SET title = $1, content = $2, level = $3, topic_id = $4, status = $5, updated_at = $6
		WHERE id = $7
	`
	question.UpdatedAt = time.Now()
	result, err := r.db.ExecContext(ctx, query,
		question.Title,
		question.Content,
		question.Level,
		question.TopicID,
		question.Status,
		question.UpdatedAt,
		question.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update question: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("question not found")
	}
	return nil
}

func (r *QuestionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM questions WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete question: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("question not found")
	}
	return nil
}
