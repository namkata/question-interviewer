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

func (s *practiceService) StartSession(ctx context.Context, userID uuid.UUID, topicID *uuid.UUID, level *string, language string, config map[string]interface{}) (*domain.PracticeSession, uuid.UUID, error) {
	// Create session
	session := domain.NewPracticeSession(userID)
	session.TopicID = topicID
	session.Level = level
	session.Language = language

	// Set default config if nil
	if config != nil {
		session.Config = config
	} else {
		session.Config = make(map[string]interface{})
	}

	// 1. Initialize Rounds if in Interview Mode
	if mode, ok := config["mode"].(string); ok && mode == "interview" {
		role := "BackEnd" // default
		if r, ok := config["role"].(string); ok && r != "" {
			role = r
		}

		rounds := getRoundsForRole(role)
		session.Config["rounds"] = rounds
		session.Config["current_round_index"] = 0

		// Override topicID for the first round
		if len(rounds) > 0 {
			firstRound := rounds[0]
			tID, err := s.repo.GetTopicIDByName(ctx, firstRound)
			if err == nil {
				topicID = &tID
			}
			// If error (topic not found), we might fall back to nil topicID and let GetRandomQuestionID handle it based on role/stack
		}
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Get first question
	questionID, err := s.repo.GetRandomQuestionID(ctx, topicID, level, language, session.Config)
	if err != nil {
		// Non-blocking error? No, we need a question to start.
		// But maybe we return session and empty question ID if none found?
		// Let's return error for now.
		return nil, uuid.Nil, fmt.Errorf("failed to get initial question: %w", err)
	}

	return session, questionID, nil
}

func getRoundsForRole(role string) []string {
	// Define standard 8 rounds for each role
	// This could be moved to DB or Config file later
	rounds := []string{}

	switch role {
	case "FrontEnd":
		rounds = append(rounds, "CV Screening", "Behavioral", "Frontend Basic", "CSS", "JavaScript", "React", "System Design", "Algorithms")

	case "BackEnd":
		rounds = append(rounds, "CV Screening", "Behavioral", "Network", "Database", "Golang", "System Design", "Algorithms", "Leadership")

	case "DevOps":
		rounds = append(rounds, "CV Screening", "Behavioral", "Network", "Docker", "Kubernetes", "CI/CD", "Terraform", "System Design")

	case "Data Engineer":
		rounds = append(rounds, "CV Screening", "Behavioral", "SQL", "Python", "Data Warehousing", "Spark", "Data Architecture", "Behavioral")

	default:
		rounds = []string{"CV Screening", "Behavioral", "Algorithms", "System Design", "Database", "Network", "Behavioral", "Leadership"}
	}

	// Ensure we have 8 rounds if possible, or truncate/pad?
	// User said "8 rounds". The above logic mostly produces 8.
	// FrontEnd: CV, Behav, Basic, CSS, JS, React, Sys, Algo = 8.
	// BackEnd: CV, Behav, Net, DB, Go, Sys, Algo, Behav = 8.
	// DevOps: CV, Behav, Net, Docker, K8s, CI/CD, TF, Sys = 8.
	// Data: CV, Behav, SQL, Py, DW, Spark, Arch, Behav = 8. // Fixed SQL, Python

	return rounds
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
	qContent, qTopic, qLevel, qCorrectAnswer, _, err := s.repo.GetQuestionContent(ctx, questionID)
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to get question content: %w", err)
	}

	var score int
	var feedbackText string
	var suggestions []string
	var improvedAnswer string

	if aiEnabled {
		// 3. Call AI Service
		// Use session language if available, otherwise fallback to request language or default
		evalLanguage := language
		if session.Language != "" {
			evalLanguage = session.Language
		}

		score, feedbackText, suggestions, improvedAnswer, err = s.ai.EvaluateAnswer(ctx, qContent, answerContent, qCorrectAnswer, qTopic, qLevel, evalLanguage)
		if err != nil {
			score = 0
			feedbackText = "AI unavailable."
			improvedAnswer = qCorrectAnswer
			suggestions = nil
		}
	} else {
		// No AI: Use database answer
		score = 0 // Not graded
		feedbackText = "Standard answer provided (AI disabled)."
		improvedAnswer = qCorrectAnswer
		suggestions = nil
	}

	// 4. Create Attempt
	attempt := domain.NewPracticeAttempt(sessionID, questionID, answerContent)
	attempt.Score = score
	attempt.Feedback = feedbackText
	attempt.Suggestions = suggestions
	attempt.ImprovedAnswer = improvedAnswer

	if err := s.repo.CreateAttempt(ctx, attempt); err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to save attempt: %w", err)
	}

	// 5. Update Session Score
	session.Score += score
	if err := s.repo.UpdateSession(ctx, session); err != nil {
		// Non-critical error
		fmt.Printf("Failed to update session score: %v\n", err)
	}

	// 6. Get Next Question
	var nextQuestionID uuid.UUID

	if mode, ok := session.Config["mode"].(string); ok && mode == "interview" {
		// Advance round logic
		rounds := []string{}
		if r, ok := session.Config["rounds"].([]interface{}); ok {
			for _, v := range r {
				rounds = append(rounds, fmt.Sprint(v))
			}
		}

		currentIdx := 0
		if idx, ok := session.Config["current_round_index"].(float64); ok {
			currentIdx = int(idx)
		}

		nextIdx := currentIdx + 1
		if nextIdx < len(rounds) {
			session.Config["current_round_index"] = nextIdx

			// Update session config in DB (persist progress)
			if err := s.repo.UpdateSession(ctx, session); err != nil {
				fmt.Printf("Failed to update session config for next round: %v\n", err)
			}

			nextTopicName := rounds[nextIdx]
			tID, err := s.repo.GetTopicIDByName(ctx, nextTopicName)
			if err == nil {
				// Use the specific topic ID for this round
				nextQuestionID, err = s.repo.GetRandomQuestionID(ctx, &tID, session.Level, session.Language, session.Config)
				if err != nil {
					nextQuestionID = uuid.Nil
				}
			} else {
				// Topic not found? Fallback to random without specific topic
				nextQuestionID, _ = s.repo.GetRandomQuestionID(ctx, nil, session.Level, session.Language, session.Config)
			}
		} else {
			// Finished all rounds
			nextQuestionID = uuid.Nil
		}
	} else {
		// Normal Practice Mode
		nextQuestionID, err = s.repo.GetRandomQuestionID(ctx, session.TopicID, session.Level, session.Language, session.Config)
		if err != nil {
			nextQuestionID = uuid.Nil
		}
	}

	return attempt, nextQuestionID, nil
}

func (s *practiceService) SuggestAnswer(ctx context.Context, questionID uuid.UUID, answerContent, language string) (int, string, []string, string, error) {
	qContent, qTopic, qLevel, qCorrectAnswer, _, err := s.repo.GetQuestionContent(ctx, questionID)
	if err != nil {
		return 0, "", nil, "", fmt.Errorf("failed to get question content: %w", err)
	}

	evalLanguage := language
	if evalLanguage == "" {
		evalLanguage = "vi"
	}

	userAnswer := strings.TrimSpace(answerContent)
	if userAnswer == "" {
		if evalLanguage == "vi" {
			userAnswer = "N/A (Ứng viên chưa trả lời. Hãy đưa ra câu trả lời mẫu hoàn chỉnh.)"
		} else {
			userAnswer = "N/A (No candidate answer. Provide a complete sample answer.)"
		}
	}

	score, feedback, suggestions, improvedAnswer, err := s.ai.EvaluateAnswer(ctx, qContent, userAnswer, qCorrectAnswer, qTopic, qLevel, evalLanguage)
	if err != nil {
		return 0, "**AI unavailable**\n\n**Standard Answer:**\n" + qCorrectAnswer, nil, qCorrectAnswer, nil
	}

	if strings.TrimSpace(improvedAnswer) == "" {
		improvedAnswer = qCorrectAnswer
	}

	return score, feedback, suggestions, improvedAnswer, nil
}

func (s *practiceService) SkipCurrentRound(ctx context.Context, sessionID uuid.UUID) (uuid.UUID, error) {
	// 1. Verify session exists
	session, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("session not found: %w", err)
	}

	if session.Status != "in_progress" {
		return uuid.Nil, fmt.Errorf("session is not in progress")
	}

	// 2. Advance Round Logic (Similar to SubmitAnswer but no score update)
	var nextQuestionID uuid.UUID

	if mode, ok := session.Config["mode"].(string); ok && mode == "interview" {
		rounds := []string{}
		if r, ok := session.Config["rounds"].([]interface{}); ok {
			for _, v := range r {
				rounds = append(rounds, fmt.Sprint(v))
			}
		}

		currentIdx := 0
		if idx, ok := session.Config["current_round_index"].(float64); ok {
			currentIdx = int(idx)
		}

		nextIdx := currentIdx + 1
		if nextIdx < len(rounds) {
			session.Config["current_round_index"] = nextIdx

			// Update session config in DB (persist progress)
			if err := s.repo.UpdateSession(ctx, session); err != nil {
				return uuid.Nil, fmt.Errorf("failed to update session config: %w", err)
			}

			nextTopicName := rounds[nextIdx]
			tID, err := s.repo.GetTopicIDByName(ctx, nextTopicName)
			if err == nil {
				// Use the specific topic ID for this round
				nextQuestionID, err = s.repo.GetRandomQuestionID(ctx, &tID, session.Level, session.Language, session.Config)
				if err != nil {
					nextQuestionID = uuid.Nil
				}
			} else {
				// Topic not found? Fallback to random without specific topic
				nextQuestionID, _ = s.repo.GetRandomQuestionID(ctx, nil, session.Level, session.Language, session.Config)
			}
		} else {
			// Finished all rounds
			nextQuestionID = uuid.Nil
		}
	} else {
		// Normal Practice Mode: Just get another question
		nextQuestionID, err = s.repo.GetRandomQuestionID(ctx, session.TopicID, session.Level, session.Language, session.Config)
		if err != nil {
			nextQuestionID = uuid.Nil
		}
	}

	return nextQuestionID, nil
}

func (s *practiceService) GetRandomQuestion(ctx context.Context, sessionID uuid.UUID, topicName *string) (uuid.UUID, error) {
	// 1. Verify session exists
	session, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("session not found: %w", err)
	}

	var topicID *uuid.UUID

	// 2. Resolve TopicName if provided
	if topicName != nil && *topicName != "" {
		tID, err := s.repo.GetTopicIDByName(ctx, *topicName)
		if err != nil {
			return uuid.Nil, fmt.Errorf("topic not found: %w", err)
		}
		topicID = &tID
	} else {
		// Use session topic if set
		topicID = session.TopicID
	}

	// 3. Get Random Question
	id, err := s.repo.GetRandomQuestionID(ctx, topicID, session.Level, session.Language, session.Config)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get random question: %w", err)
	}

	return id, nil
}

func (s *practiceService) GetSession(ctx context.Context, id uuid.UUID) (*domain.PracticeSession, error) {
	return s.repo.GetSession(ctx, id)
}

func (s *practiceService) GetQuestion(ctx context.Context, questionID uuid.UUID) (string, string, string, string, string, error) {
	return s.repo.GetQuestionContent(ctx, questionID)
}

func (s *practiceService) GetTopicIDByName(ctx context.Context, name string) (uuid.UUID, error) {
	return s.repo.GetTopicIDByName(ctx, name)
}

func (s *practiceService) CreateQuestion(ctx context.Context, content, topic, level, correctAnswer, hint string) (*domain.Question, error) {
	question := domain.NewQuestion(content, topic, level)
	question.CorrectAnswer = correctAnswer
	question.Hint = hint

	// We need to resolve Topic ID if possible, but for now NewQuestion just sets TopicName.
	// The repo implementation of CreateQuestion should handle the TopicName -> TopicID mapping or insertion if needed.
	// However, looking at repo interface, it takes *domain.Question.

	// Let's check if we can resolve topic ID here
	if topic != "" {
		topicID, err := s.repo.GetTopicIDByName(ctx, topic)
		if err == nil {
			question.TopicID = &topicID
		}
		// If error (topic not found), we might want to create it or just leave it nil/fail?
		// For now, let's assume the repo handles it or we proceed with just TopicName if the repo supports it.
		// But domain.Question struct has TopicID *uuid.UUID.
	}

	if err := s.repo.CreateQuestion(ctx, question); err != nil {
		return nil, fmt.Errorf("failed to create question: %w", err)
	}
	return question, nil
}
