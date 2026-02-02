package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/question-interviewer/answer-service/internal/domain"
)

type AnswerRepository interface {
	Create(ctx context.Context, answer *domain.Answer) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Answer, error)
	ListByQuestionID(ctx context.Context, questionID uuid.UUID, limit, offset int) ([]*domain.Answer, error)
	Update(ctx context.Context, answer *domain.Answer) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type AnswerService interface {
	CreateAnswer(ctx context.Context, questionID, authorID uuid.UUID, content string, answerType domain.AnswerType) (*domain.Answer, error)
	GetAnswer(ctx context.Context, id uuid.UUID) (*domain.Answer, error)
	ListAnswersForQuestion(ctx context.Context, questionID uuid.UUID, limit, offset int) ([]*domain.Answer, error)
}
