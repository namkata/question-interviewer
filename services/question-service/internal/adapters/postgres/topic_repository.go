package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/question-interviewer/question-service/internal/domain"
	"github.com/question-interviewer/question-service/internal/ports"
)

type TopicRepository struct {
	db *sql.DB
}

func NewTopicRepository(db *sql.DB) ports.TopicRepository {
	return &TopicRepository{
		db: db,
	}
}

func (r *TopicRepository) Create(ctx context.Context, topic *domain.Topic) error {
	query := `
		INSERT INTO topics (id, name, description)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query,
		topic.ID,
		topic.Name,
		topic.Description,
	)
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}
	return nil
}

func (r *TopicRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Topic, error) {
	query := `
		SELECT id, name, description
		FROM topics
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var t domain.Topic
	err := row.Scan(
		&t.ID,
		&t.Name,
		&t.Description,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("topic not found")
		}
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}
	return &t, nil
}

func (r *TopicRepository) GetByName(ctx context.Context, name string) (*domain.Topic, error) {
	query := `
		SELECT id, name, description
		FROM topics
		WHERE name = $1
	`
	row := r.db.QueryRowContext(ctx, query, name)

	var t domain.Topic
	err := row.Scan(
		&t.ID,
		&t.Name,
		&t.Description,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("topic not found")
		}
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}
	return &t, nil
}

func (r *TopicRepository) List(ctx context.Context) ([]*domain.Topic, error) {
	query := `
		SELECT id, name, description
		FROM topics
		ORDER BY name ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list topics: %w", err)
	}
	defer rows.Close()

	var topics []*domain.Topic
	for rows.Next() {
		var t domain.Topic
		err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan topic: %w", err)
		}
		topics = append(topics, &t)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return topics, nil
}
