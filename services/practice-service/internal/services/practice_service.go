package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/question-interviewer/practice-service/internal/domain"
	"github.com/question-interviewer/practice-service/internal/ports"
)

type practiceService struct {
	repo ports.PracticeRepository
	ai   ports.AIService
}

func NewPracticeService(repo ports.PracticeRepository, ai ports.AIService) ports.PracticeService {
	return &practiceService{
		repo: repo,
		ai:   ai,
	}
}

func (s *practiceService) CreateQuestion(ctx context.Context, content, topic, level, correctAnswer string) (*domain.Question, error) {
	q := domain.NewQuestion(content, topic, level)
	q.CorrectAnswer = correctAnswer
	if err := s.repo.CreateQuestion(ctx, q); err != nil {
		return nil, fmt.Errorf("failed to create question: %w", err)
	}
	return q, nil
}

func (s *practiceService) StartSession(ctx context.Context, userID uuid.UUID, topicID *uuid.UUID, level *string) (*domain.PracticeSession, uuid.UUID, error) {
	// 1. Create a new session
	session := domain.NewPracticeSession(userID)
	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to create session: %w", err)
	}

	// 2. Get the first random question
	firstQuestionID, err := s.repo.GetRandomQuestionID(ctx, topicID, level)
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to get initial question: %w", err)
	}

	return session, firstQuestionID, nil
}

func (s *practiceService) SubmitAnswer(ctx context.Context, sessionID, questionID uuid.UUID, answerContent, language string, aiEnabled bool) (*domain.PracticeAttempt, uuid.UUID, error) {
	// 1. Verify session exists
	session, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("session not found: %w", err)
	}

	if session.Status != "in_progress" {
		return nil, uuid.Nil, fmt.Errorf("session is not in progress")
	}

	// 2. Get Question Data (Content, Topic, Level, CorrectAnswer)
	qContent, qTopic, qLevel, qCorrectAnswer, err := s.repo.GetQuestionContent(ctx, questionID)
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to get question content: %w", err)
	}

	var score int
	var fullFeedback string

	if aiEnabled {
		// 3. Call AI Service
		var feedback, improvedAnswer string
		var suggestions []string
		score, feedback, suggestions, improvedAnswer, err = s.ai.EvaluateAnswer(ctx, qContent, answerContent, qTopic, qLevel, language)
		if err != nil {
			// Fallback or error? For now, log and error, or return attempt with error status.
			// Let's return error for now to keep it simple.
			return nil, uuid.Nil, fmt.Errorf("AI evaluation failed: %w", err)
		}

		// Format feedback with suggestions and improved answer
		fullFeedback = fmt.Sprintf("%s\n\n**Suggestions:**\n%s\n\n**Improved Answer:**\n%s",
			feedback,
			strings.Join(suggestions, "\n- "),
			improvedAnswer)
	} else {
		// No AI: Use database answer
		score = 0 // Not graded
		fullFeedback = fmt.Sprintf("**Standard Answer (No AI):**\n%s", qCorrectAnswer)
	}

	// 4. Create Attempt
	attempt := domain.NewPracticeAttempt(sessionID, questionID, answerContent)
	attempt.Score = score
	attempt.Feedback = fullFeedback

	if err := s.repo.CreateAttempt(ctx, attempt); err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to save attempt: %w", err)
	}

	// 5. Update Session Score
	session.Score += score
	if err := s.repo.UpdateSession(ctx, session); err != nil {
		// Non-critical error
		fmt.Printf("Failed to update session score: %v\n", err)
	}

	// 6. Get Next Question (Random for now)
	nextQuestionID, err := s.repo.GetRandomQuestionID(ctx, nil, nil) // Keep same topic/level? For now random
	if err != nil {
		// If we can't get a next question, just return nil UUID, frontend handles it (e.g. "Session Complete")
		nextQuestionID = uuid.Nil
	}

	return attempt, nextQuestionID, nil
}

func (s *practiceService) GetSession(ctx context.Context, id uuid.UUID) (*domain.PracticeSession, error) {
	return s.repo.GetSession(ctx, id)
}

func (s *practiceService) GetQuestion(ctx context.Context, questionID uuid.UUID) (string, string, string, string, error) {
	return s.repo.GetQuestionContent(ctx, questionID)
}

func (s *practiceService) GetTopicIDByName(ctx context.Context, name string) (uuid.UUID, error) {
	return s.repo.GetTopicIDByName(ctx, name)
}
