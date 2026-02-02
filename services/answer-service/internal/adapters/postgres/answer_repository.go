package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/question-interviewer/answer-service/internal/domain"
	"github.com/question-interviewer/answer-service/internal/ports"
)

type AnswerRepository struct {
	db *sql.DB
}

func NewAnswerRepository(db *sql.DB) ports.AnswerRepository {
	return &AnswerRepository{
		db: db,
	}
}

func (r *AnswerRepository) Create(ctx context.Context, answer *domain.Answer) error {
	query := `
		INSERT INTO answers (id, question_id, content, author_id, answer_type, vote_count, is_accepted, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		answer.ID,
		answer.QuestionID,
		answer.Content,
		answer.AuthorID,
		answer.AnswerType,
		answer.VoteCount,
		answer.IsAccepted,
		answer.CreatedAt,
		answer.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create answer: %w", err)
	}
	return nil
}

func (r *AnswerRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Answer, error) {
	query := `
		SELECT id, question_id, content, author_id, answer_type, vote_count, is_accepted, created_at, updated_at
		FROM answers
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var a domain.Answer
	err := row.Scan(
		&a.ID,
		&a.QuestionID,
		&a.Content,
		&a.AuthorID,
		&a.AnswerType,
		&a.VoteCount,
		&a.IsAccepted,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("answer not found")
		}
		return nil, fmt.Errorf("failed to get answer: %w", err)
	}
	return &a, nil
}

func (r *AnswerRepository) ListByQuestionID(ctx context.Context, questionID uuid.UUID, limit, offset int) ([]*domain.Answer, error) {
	query := `
		SELECT id, question_id, content, author_id, answer_type, vote_count, is_accepted, created_at, updated_at
		FROM answers
		WHERE question_id = $1
		ORDER BY vote_count DESC, created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, questionID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list answers: %w", err)
	}
	defer rows.Close()

	var answers []*domain.Answer
	for rows.Next() {
		var a domain.Answer
		err := rows.Scan(
			&a.ID,
			&a.QuestionID,
			&a.Content,
			&a.AuthorID,
			&a.AnswerType,
			&a.VoteCount,
			&a.IsAccepted,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan answer: %w", err)
		}
		answers = append(answers, &a)
	}
	return answers, nil
}

func (r *AnswerRepository) Update(ctx context.Context, answer *domain.Answer) error {
	query := `
		UPDATE answers
		SET content = $1, answer_type = $2, is_accepted = $3, updated_at = $4
		WHERE id = $5
	`
	answer.UpdatedAt = time.Now()
	result, err := r.db.ExecContext(ctx, query,
		answer.Content,
		answer.AnswerType,
		answer.IsAccepted,
		answer.UpdatedAt,
		answer.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update answer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("answer not found")
	}
	return nil
}

func (r *AnswerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM answers WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete answer: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("answer not found")
	}
	return nil
}
