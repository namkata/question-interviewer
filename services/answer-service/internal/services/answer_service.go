package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/question-interviewer/answer-service/internal/domain"
	"github.com/question-interviewer/answer-service/internal/ports"
)

type answerService struct {
	repo ports.AnswerRepository
}

func NewAnswerService(repo ports.AnswerRepository) ports.AnswerService {
	return &answerService{
		repo: repo,
	}
}

func (s *answerService) CreateAnswer(ctx context.Context, questionID, authorID uuid.UUID, content string, answerType domain.AnswerType) (*domain.Answer, error) {
	// TODO: Validate questionID exists via QuestionService (gRPC/HTTP call)
	answer := domain.NewAnswer(questionID, authorID, content, answerType)
	if err := s.repo.Create(ctx, answer); err != nil {
		return nil, err
	}
	return answer, nil
}

func (s *answerService) GetAnswer(ctx context.Context, id uuid.UUID) (*domain.Answer, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *answerService) ListAnswersForQuestion(ctx context.Context, questionID uuid.UUID, limit, offset int) ([]*domain.Answer, error) {
	return s.repo.ListByQuestionID(ctx, questionID, limit, offset)
}
