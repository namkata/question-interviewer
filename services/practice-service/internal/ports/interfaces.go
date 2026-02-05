package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/question-interviewer/practice-service/internal/domain"
)

type PracticeRepository interface {
	CreateSession(ctx context.Context, session *domain.PracticeSession) error
	GetSession(ctx context.Context, id uuid.UUID) (*domain.PracticeSession, error)
	UpdateSession(ctx context.Context, session *domain.PracticeSession) error

	// Attempts
	CreateAttempt(ctx context.Context, attempt *domain.PracticeAttempt) error

	// Helper method to get a random question ID for the session
	GetRandomQuestionID(ctx context.Context, topicID *uuid.UUID, level *string, language string, config map[string]interface{}) (uuid.UUID, error)

	// Helper to get question content (needed for AI) - in real microservices this might come from Question Service,
	// but here we have direct DB access for now.
	GetQuestionContent(ctx context.Context, questionID uuid.UUID) (string, string, string, string, string, error) // returns content, topic, level, correctAnswer, hint
	CreateQuestion(ctx context.Context, question *domain.Question) error
	GetTopicIDByName(ctx context.Context, name string) (uuid.UUID, error)
}

type AIService interface {
	EvaluateAnswer(ctx context.Context, question, userAnswer, correctAnswer, topic, level, language string) (int, string, []string, string, error) // score, feedback, suggestions, improvedAnswer
}

type PracticeService interface {
	StartSession(ctx context.Context, userID uuid.UUID, topicID *uuid.UUID, level *string, language string, config map[string]interface{}) (*domain.PracticeSession, uuid.UUID, error)
	SubmitAnswer(ctx context.Context, sessionID, questionID uuid.UUID, answerContent, language string, aiEnabled bool) (*domain.PracticeAttempt, uuid.UUID, error)
	SuggestAnswer(ctx context.Context, questionID uuid.UUID, answerContent, language string) (int, string, []string, string, error)
	SkipCurrentRound(ctx context.Context, sessionID uuid.UUID) (uuid.UUID, error)
	GetSession(ctx context.Context, id uuid.UUID) (*domain.PracticeSession, error)
	GetQuestion(ctx context.Context, questionID uuid.UUID) (string, string, string, string, string, error) // returns content, topic, level, correctAnswer, hint
	GetRandomQuestion(ctx context.Context, sessionID uuid.UUID, topicName *string) (uuid.UUID, error)
	CreateQuestion(ctx context.Context, content, topic, level, correctAnswer, hint string) (*domain.Question, error)
	GetTopicIDByName(ctx context.Context, name string) (uuid.UUID, error)
}
