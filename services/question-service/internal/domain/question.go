package domain

import (
	"time"

	"github.com/google/uuid"
)

// Question represents the core domain entity for a question
type Question struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Level         string    `json:"level"`
	Language      string    `json:"language"`
	Role          string    `json:"role"`
	Hint          string    `json:"hint"`
	CorrectAnswer string    `json:"correct_answer"`
	TopicID       uuid.UUID `json:"topic_id"`
	CreatedBy     uuid.UUID `json:"created_by"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// NewQuestion creates a new question instance
func NewQuestion(title, content, level, language, role, hint, correctAnswer string, topicID, createdBy uuid.UUID) *Question {
	return &Question{
		ID:            uuid.New(),
		Title:         title,
		Content:       content,
		Level:         level,
		Language:      language,
		Role:          role,
		Hint:          hint,
		CorrectAnswer: correctAnswer,
		TopicID:       topicID,
		CreatedBy:     createdBy,
		Status:        "published",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}
