package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/question-interviewer/practice-service/internal/domain"
	"github.com/question-interviewer/practice-service/internal/ports"
)

type practiceService struct {
	repo      ports.PracticeRepository
	ai        ports.AIService
	aiEnabled bool
}

func NewPracticeService(repo ports.PracticeRepository, ai ports.AIService, aiEnabled bool) ports.PracticeService {
	return &practiceService{
		repo:      repo,
		ai:        ai,
		aiEnabled: aiEnabled,
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
	session.Config["session_id"] = session.ID.String()

	// 1. Initialize Rounds if in Interview Mode
	if mode, ok := session.Config["mode"].(string); ok && mode == "interview" {
		role := "BackEnd"
		if r, ok := session.Config["role"].(string); ok && r != "" {
			role = r
		}

		ensureRounds := func() []map[string]interface{} {
			var out []map[string]interface{}

			if raw, ok := session.Config["rounds"]; ok && raw != nil {
				if list, ok := raw.([]interface{}); ok {
					for i, v := range list {
						if m, ok := v.(map[string]interface{}); ok {
							out = append(out, m)
							continue
						}
						name := strings.TrimSpace(fmt.Sprint(v))
						out = append(out, map[string]interface{}{"name": name, "topic": name, "count": 2})
						_ = i
					}
				}
			}

			if len(out) == 0 {
				roundNames := getRoundsForRole(role)
				for _, name := range roundNames {
					out = append(out, map[string]interface{}{"name": name, "topic": name, "count": 2})
				}
			}

			for i := range out {
				if _, ok := out[i]["id"]; !ok {
					out[i]["id"] = fmt.Sprintf("%d", i)
				}
				if _, ok := out[i]["topic"]; !ok {
					out[i]["topic"] = out[i]["name"]
				}
				status := "pending"
				if i == 0 {
					status = "in_progress"
				}
				out[i]["status"] = status
			}

			return out
		}

		if _, ok := session.Config["questions_per_round"]; !ok {
			session.Config["questions_per_round"] = 2
		}

		rounds := ensureRounds()
		session.Config["rounds"] = rounds
		session.Config["current_round_index"] = 0
		session.Config["current_round_question_index"] = 0

		firstTopic := strings.TrimSpace(fmt.Sprint(rounds[0]["topic"]))
		if firstTopic == "" {
			firstTopic = strings.TrimSpace(fmt.Sprint(rounds[0]["name"]))
		}
		if firstTopic != "" {
			tID, err := s.repo.GetTopicIDByName(ctx, firstTopic)
			if err == nil {
				topicID = &tID
			}
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

	_, _, qLevel, _, _, err := s.repo.GetQuestionContent(ctx, questionID)
	if err == nil && strings.TrimSpace(qLevel) != "" {
		session.Config["difficulty_floor"] = qLevel
		_ = s.repo.UpdateSession(ctx, session)
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

	if aiEnabled && s.aiEnabled {
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
	attempt.QuestionContent = &qContent
	attempt.Score = score
	attempt.Feedback = feedbackText
	attempt.Suggestions = suggestions
	attempt.ImprovedAnswer = improvedAnswer

	if session.Config == nil {
		session.Config = make(map[string]interface{})
	}
	session.Config["session_id"] = session.ID.String()
	if strings.TrimSpace(qLevel) != "" {
		session.Config["difficulty_floor"] = qLevel
	}

	if mode, ok := session.Config["mode"].(string); ok && mode == "interview" {
		currentIdx := 0
		switch idx := session.Config["current_round_index"].(type) {
		case float64:
			currentIdx = int(idx)
		case int:
			currentIdx = idx
		case int64:
			currentIdx = int(idx)
		}
		attempt.RoundIndex = &currentIdx

		roundsAny := []interface{}{}
		switch v := session.Config["rounds"].(type) {
		case []interface{}:
			roundsAny = v
		case []map[string]interface{}:
			for _, m := range v {
				roundsAny = append(roundsAny, m)
			}
		}

		if currentIdx >= 0 && currentIdx < len(roundsAny) {
			if roundObj, ok := roundsAny[currentIdx].(map[string]interface{}); ok {
				if name, ok := roundObj["name"].(string); ok && name != "" {
					attempt.RoundName = &name
				}
			}
		}
	}

	if err := s.repo.CreateAttempt(ctx, attempt); err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to save attempt: %w", err)
	}

	// 5. Update Session Score
	session.Score += score

	// 6. Get Next Question
	var nextQuestionID uuid.UUID

	if mode, ok := session.Config["mode"].(string); ok && mode == "interview" {
		rounds := []interface{}{}
		switch v := session.Config["rounds"].(type) {
		case []interface{}:
			rounds = v
		case []map[string]interface{}:
			for _, m := range v {
				rounds = append(rounds, m)
			}
		}

		currentRoundIdx := 0
		switch idx := session.Config["current_round_index"].(type) {
		case float64:
			currentRoundIdx = int(idx)
		case int:
			currentRoundIdx = idx
		case int64:
			currentRoundIdx = int(idx)
		}

		answeredInRound := 0
		switch idx := session.Config["current_round_question_index"].(type) {
		case float64:
			answeredInRound = int(idx)
		case int:
			answeredInRound = idx
		case int64:
			answeredInRound = int(idx)
		}
		answeredInRound++

		questionsPerRound := 2
		if currentRoundIdx >= 0 && currentRoundIdx < len(rounds) {
			if roundObj, ok := rounds[currentRoundIdx].(map[string]interface{}); ok {
				switch v := roundObj["count"].(type) {
				case float64:
					if int(v) > 0 {
						questionsPerRound = int(v)
					}
				case int:
					if v > 0 {
						questionsPerRound = v
					}
				case int64:
					if int(v) > 0 {
						questionsPerRound = int(v)
					}
				}
			}
		}
		if questionsPerRound <= 0 {
			questionsPerRound = 2
		}
		switch v := session.Config["questions_per_round"].(type) {
		case float64:
			if questionsPerRound == 2 && int(v) > 0 {
				questionsPerRound = int(v)
			}
		case int:
			if questionsPerRound == 2 && v > 0 {
				questionsPerRound = v
			}
		case int64:
			if questionsPerRound == 2 && int(v) > 0 {
				questionsPerRound = int(v)
			}
		}

		advanceRound := answeredInRound >= questionsPerRound
		if advanceRound {
			session.Config["current_round_question_index"] = 0
		} else {
			session.Config["current_round_question_index"] = answeredInRound
		}

		nextRoundIdx := currentRoundIdx
		if advanceRound {
			if currentRoundIdx >= 0 && currentRoundIdx < len(rounds) {
				if roundObj, ok := rounds[currentRoundIdx].(map[string]interface{}); ok {
					roundObj["status"] = "completed"
				}
			}
			nextRoundIdx = currentRoundIdx + 1
			session.Config["current_round_index"] = nextRoundIdx
		}

		if nextRoundIdx < len(rounds) {
			if roundObj, ok := rounds[nextRoundIdx].(map[string]interface{}); ok {
				roundObj["status"] = "in_progress"
				topic := strings.TrimSpace(fmt.Sprint(roundObj["topic"]))
				if topic == "" {
					topic = strings.TrimSpace(fmt.Sprint(roundObj["name"]))
				}
				if topic != "" {
					tID, err := s.repo.GetTopicIDByName(ctx, topic)
					if err == nil {
						nextQuestionID, err = s.repo.GetRandomQuestionID(ctx, &tID, session.Level, session.Language, session.Config)
						if err != nil {
							nextQuestionID = uuid.Nil
						}
					} else {
						nextQuestionID, _ = s.repo.GetRandomQuestionID(ctx, nil, session.Level, session.Language, session.Config)
					}
				} else {
					nextQuestionID, _ = s.repo.GetRandomQuestionID(ctx, nil, session.Level, session.Language, session.Config)
				}
			} else {
				nextQuestionID, _ = s.repo.GetRandomQuestionID(ctx, nil, session.Level, session.Language, session.Config)
			}
		} else {
			now := time.Now()
			session.Status = "completed"
			session.EndedAt = &now
			nextQuestionID = uuid.Nil
		}
	} else {
		// Normal Practice Mode
		nextQuestionID, err = s.repo.GetRandomQuestionID(ctx, session.TopicID, session.Level, session.Language, session.Config)
		if err != nil {
			nextQuestionID = uuid.Nil
		}
	}

	if nextQuestionID != uuid.Nil {
		_, _, nextLevel, _, _, err := s.repo.GetQuestionContent(ctx, nextQuestionID)
		if err == nil && strings.TrimSpace(nextLevel) != "" {
			session.Config["difficulty_floor"] = nextLevel
		}
	}

	if err := s.repo.UpdateSession(ctx, session); err != nil {
		fmt.Printf("Failed to update session: %v\n", err)
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
	requestingSample := userAnswer == ""

	if requestingSample {
		sampleAnswer, sampleFeedback, sampleSuggestions, sampleSource, err := s.repo.GetQuestionSampleCache(ctx, questionID)
		if err == nil && strings.TrimSpace(sampleAnswer) != "" && strings.TrimSpace(sampleSource) == "ai" {
			return 0, sampleFeedback, sampleSuggestions, sampleAnswer, nil
		}

		if !s.aiEnabled {
			if err == nil && strings.TrimSpace(sampleAnswer) != "" {
				return 0, sampleFeedback, sampleSuggestions, sampleAnswer, nil
			}
			return 0, "", nil, qCorrectAnswer, nil
		}

		if evalLanguage == "vi" {
			userAnswer = "N/A (Ứng viên chưa trả lời. Hãy đưa ra câu trả lời mẫu hoàn chỉnh.)"
		} else {
			userAnswer = "N/A (No candidate answer. Provide a complete sample answer.)"
		}
	}

	if !s.aiEnabled {
		return 0, "", nil, qCorrectAnswer, nil
	}

	score, feedback, suggestions, improvedAnswer, err := s.ai.EvaluateAnswer(ctx, qContent, userAnswer, qCorrectAnswer, qTopic, qLevel, evalLanguage)
	if err != nil {
		return 0, "", nil, qCorrectAnswer, nil
	}

	if strings.TrimSpace(improvedAnswer) == "" {
		improvedAnswer = qCorrectAnswer
	}

	if requestingSample {
		_ = s.repo.UpsertQuestionSampleCache(ctx, questionID, improvedAnswer, feedback, suggestions, "ai")
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
		session.Config["session_id"] = session.ID.String()

		rounds := []interface{}{}
		switch v := session.Config["rounds"].(type) {
		case []interface{}:
			rounds = v
		case []map[string]interface{}:
			for _, m := range v {
				rounds = append(rounds, m)
			}
		}

		currentIdx := 0
		switch idx := session.Config["current_round_index"].(type) {
		case float64:
			currentIdx = int(idx)
		case int:
			currentIdx = idx
		case int64:
			currentIdx = int(idx)
		}

		if currentIdx >= 0 && currentIdx < len(rounds) {
			if roundObj, ok := rounds[currentIdx].(map[string]interface{}); ok {
				roundObj["status"] = "skipped"
			}
		}

		nextIdx := currentIdx + 1
		session.Config["current_round_index"] = nextIdx
		session.Config["current_round_question_index"] = 0

		if nextIdx < len(rounds) {
			if roundObj, ok := rounds[nextIdx].(map[string]interface{}); ok {
				roundObj["status"] = "in_progress"
				topic := strings.TrimSpace(fmt.Sprint(roundObj["topic"]))
				if topic == "" {
					topic = strings.TrimSpace(fmt.Sprint(roundObj["name"]))
				}
				if topic != "" {
					tID, err := s.repo.GetTopicIDByName(ctx, topic)
					if err == nil {
						nextQuestionID, err = s.repo.GetRandomQuestionID(ctx, &tID, session.Level, session.Language, session.Config)
						if err != nil {
							nextQuestionID = uuid.Nil
						}
					} else {
						nextQuestionID, _ = s.repo.GetRandomQuestionID(ctx, nil, session.Level, session.Language, session.Config)
					}
				} else {
					nextQuestionID, _ = s.repo.GetRandomQuestionID(ctx, nil, session.Level, session.Language, session.Config)
				}
			} else {
				nextQuestionID, _ = s.repo.GetRandomQuestionID(ctx, nil, session.Level, session.Language, session.Config)
			}
		} else {
			now := time.Now()
			session.Status = "completed"
			session.EndedAt = &now
			nextQuestionID = uuid.Nil
		}

		if nextQuestionID != uuid.Nil {
			_, _, nextLevel, _, _, err := s.repo.GetQuestionContent(ctx, nextQuestionID)
			if err == nil && strings.TrimSpace(nextLevel) != "" {
				session.Config["difficulty_floor"] = nextLevel
			}
		}

		if err := s.repo.UpdateSession(ctx, session); err != nil {
			return uuid.Nil, fmt.Errorf("failed to update session config: %w", err)
		}
	} else {
		// Normal Practice Mode: Just get another question
		session.Config["session_id"] = session.ID.String()
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
	if session.Config == nil {
		session.Config = make(map[string]interface{})
	}
	session.Config["session_id"] = session.ID.String()
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

func (s *practiceService) ListAttempts(ctx context.Context, sessionID uuid.UUID) ([]*domain.PracticeAttempt, error) {
	return s.repo.ListAttemptsBySession(ctx, sessionID)
}

func (s *practiceService) GetSessionSummary(ctx context.Context, sessionID uuid.UUID, language string) (map[string]interface{}, error) {
	session, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	attempts, err := s.repo.ListAttemptsBySession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list attempts: %w", err)
	}

	type roundAgg struct {
		sum   int
		count int
	}

	aggByName := map[string]*roundAgg{}
	for _, a := range attempts {
		name := "General"
		if a.RoundName != nil && *a.RoundName != "" {
			name = *a.RoundName
		}
		if _, ok := aggByName[name]; !ok {
			aggByName[name] = &roundAgg{}
		}
		aggByName[name].sum += a.Score
		aggByName[name].count++
	}

	roundsOut := make([]map[string]interface{}, 0)
	if mode, ok := session.Config["mode"].(string); ok && mode == "interview" {
		roundsAny := []interface{}{}
		switch v := session.Config["rounds"].(type) {
		case []interface{}:
			roundsAny = v
		case []map[string]interface{}:
			for _, m := range v {
				roundsAny = append(roundsAny, m)
			}
		}

		for i, r := range roundsAny {
			roundObj, ok := r.(map[string]interface{})
			if !ok {
				continue
			}
			name, _ := roundObj["name"].(string)
			id, _ := roundObj["id"].(string)
			status, _ := roundObj["status"].(string)
			a := aggByName[name]
			avg := 0.0
			count := 0
			if a != nil && a.count > 0 {
				avg = float64(a.sum) / float64(a.count)
				count = a.count
			}
			roundsOut = append(roundsOut, map[string]interface{}{
				"index":         i,
				"id":            id,
				"name":          name,
				"status":        status,
				"attempt_count": count,
				"avg_score":     avg,
			})
		}
	}

	overallAvg := 0.0
	if len(attempts) > 0 {
		sum := 0
		for _, a := range attempts {
			sum += a.Score
		}
		overallAvg = float64(sum) / float64(len(attempts))
	}

	role := ""
	if v, ok := session.Config["role"].(string); ok {
		role = v
	}
	if language == "" {
		language = session.Language
	}
	if language == "" {
		language = "en"
	}

	summary := map[string]interface{}{
		"session_id":    session.ID,
		"status":        session.Status,
		"overall_score": overallAvg,
		"rounds":        roundsOut,
	}

	if s.aiEnabled && s.ai != nil {
		attemptVals := make([]domain.PracticeAttempt, 0, len(attempts))
		for _, a := range attempts {
			attemptVals = append(attemptVals, *a)
		}
		strengths, weaknesses, readiness, overallScore, err := s.ai.SummarizeInterview(ctx, role, language, attemptVals)
		if err == nil {
			summary["ai_summary"] = map[string]interface{}{
				"strengths":     strengths,
				"weaknesses":    weaknesses,
				"readiness":     readiness,
				"overall_score": overallScore,
			}
		}
	}

	return summary, nil
}
