package domain

import (
	"time"

	"github.com/google/uuid"
)

type PracticeSession struct {
	ID        uuid.UUID              `json:"id"`
	UserID    uuid.UUID              `json:"user_id"`
	Score     int                    `json:"score"`
	StartedAt time.Time              `json:"started_at"`
	EndedAt   *time.Time             `json:"ended_at"` // Pointer to allow null
	Status    string                 `json:"status"`   // e.g., "in_progress", "completed"
	TopicID   *uuid.UUID             `json:"topic_id"`
	Level     *string                `json:"level"`
	Language  string                 `json:"language"`
	Config    map[string]interface{} `json:"config"`
}

type PracticeAttempt struct {
	ID             uuid.UUID `json:"id"`
	SessionID      uuid.UUID `json:"session_id"`
	QuestionID     uuid.UUID `json:"question_id"`
	UserAnswer     string    `json:"user_answer"`
	Score          int       `json:"score"`                     // 0-100 from AI
	Feedback       string    `json:"feedback"`                  // Short feedback text (no suggestions)
	Suggestions    []string  `json:"suggestions,omitempty"`     // Returned in API; not persisted
	ImprovedAnswer string    `json:"improved_answer,omitempty"` // Returned in API; not persisted
	CreatedAt      time.Time `json:"created_at"`
}

type Question struct {
	ID            uuid.UUID  `json:"id"`
	Content       string     `json:"content"`
	TopicID       *uuid.UUID `json:"topic_id"` // Simplified for now
	Level         string     `json:"level"`
	CorrectAnswer string     `json:"correct_answer"`
	Hint          string     `json:"hint"`
	CreatedAt     time.Time  `json:"created_at"`
	// We might need to handle Topic string vs ID, but for now let's assume simple string mapping or null
	TopicName string `json:"topic"` // Helper for now
}

func NewQuestion(content, topic, level string) *Question {
	return &Question{
		ID:        uuid.New(),
		Content:   content,
		TopicName: topic,
		Level:     level,
		CreatedAt: time.Now(),
	}
}

func NewPracticeSession(userID uuid.UUID) *PracticeSession {
	return &PracticeSession{
		ID:        uuid.New(),
		UserID:    userID,
		Score:     0,
		StartedAt: time.Now(),
		Status:    "in_progress",
		Config:    make(map[string]interface{}),
	}
}

func NewPracticeAttempt(sessionID, questionID uuid.UUID, userAnswer string) *PracticeAttempt {
	return &PracticeAttempt{
		ID:         uuid.New(),
		SessionID:  sessionID,
		QuestionID: questionID,
		UserAnswer: userAnswer,
		CreatedAt:  time.Now(),
	}
}
