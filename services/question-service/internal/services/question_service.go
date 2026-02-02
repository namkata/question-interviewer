package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/question-interviewer/question-service/internal/domain"
	"github.com/question-interviewer/question-service/internal/ports"
)

type questionService struct {
	repo ports.QuestionRepository
}

// NewQuestionService creates a new instance of QuestionService
func NewQuestionService(repo ports.QuestionRepository) ports.QuestionService {
	return &questionService{
		repo: repo,
	}
}

func (s *questionService) CreateQuestion(ctx context.Context, title, content, level string, topicID, createdBy uuid.UUID) (*domain.Question, error) {
	question := domain.NewQuestion(title, content, level, topicID, createdBy)
	if err := s.repo.Create(ctx, question); err != nil {
		return nil, err
	}
	return question, nil
}

func (s *questionService) GetQuestion(ctx context.Context, id uuid.UUID) (*domain.Question, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *questionService) ListQuestions(ctx context.Context, limit, offset int) ([]*domain.Question, error) {
	return s.repo.List(ctx, limit, offset)
}
