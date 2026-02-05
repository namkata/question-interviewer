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

// TopicRepository defines the interface for topic data access
type TopicRepository interface {
	Create(ctx context.Context, topic *domain.Topic) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Topic, error)
	GetByName(ctx context.Context, name string) (*domain.Topic, error)
	List(ctx context.Context) ([]*domain.Topic, error)
}

// QuestionService defines the interface for business logic
type QuestionService interface {
	CreateQuestion(ctx context.Context, title, content, level, language, role, hint, correctAnswer string, topicID, createdBy uuid.UUID) (*domain.Question, error)
	GetQuestion(ctx context.Context, id uuid.UUID) (*domain.Question, error)
	ListQuestions(ctx context.Context, limit, offset int) ([]*domain.Question, error)

	// Topic methods
	CreateTopic(ctx context.Context, name, description string) (*domain.Topic, error)
	GetTopicByName(ctx context.Context, name string) (*domain.Topic, error)
	ListTopics(ctx context.Context) ([]*domain.Topic, error)
}
