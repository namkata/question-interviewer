package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/question-interviewer/question-service/internal/domain"
	"github.com/question-interviewer/question-service/internal/ports"
)

type questionService struct {
	repo      ports.QuestionRepository
	topicRepo ports.TopicRepository
}

// NewQuestionService creates a new instance of QuestionService
func NewQuestionService(repo ports.QuestionRepository, topicRepo ports.TopicRepository) ports.QuestionService {
	return &questionService{
		repo:      repo,
		topicRepo: topicRepo,
	}
}

func (s *questionService) CreateQuestion(ctx context.Context, title, content, level, language, role, hint, correctAnswer string, topicID, createdBy uuid.UUID) (*domain.Question, error) {
	question := domain.NewQuestion(title, content, level, language, role, hint, correctAnswer, topicID, createdBy)
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

// Topic methods
func (s *questionService) CreateTopic(ctx context.Context, name, description string) (*domain.Topic, error) {
	topic := domain.NewTopic(name, description)
	if err := s.topicRepo.Create(ctx, topic); err != nil {
		return nil, err
	}
	return topic, nil
}

func (s *questionService) GetTopicByName(ctx context.Context, name string) (*domain.Topic, error) {
	return s.topicRepo.GetByName(ctx, name)
}

func (s *questionService) ListTopics(ctx context.Context) ([]*domain.Topic, error) {
	return s.topicRepo.List(ctx)
}
