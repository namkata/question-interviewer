package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/question-interviewer/question-service/internal/domain"
)

// QuestionRepository defines the interface for data access
type QuestionRepository interface {
	Create(ctx context.Context, question *domain.Question) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Question, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Question, error)
	Update(ctx context.Context, question *domain.Question) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// QuestionService defines the interface for business logic
type QuestionService interface {
	CreateQuestion(ctx context.Context, title, content, level string, topicID, createdBy uuid.UUID) (*domain.Question, error)
	GetQuestion(ctx context.Context, id uuid.UUID) (*domain.Question, error)
	ListQuestions(ctx context.Context, limit, offset int) ([]*domain.Question, error)
}
