package domain

import (
	"time"

	"github.com/google/uuid"
)

type AnswerType string

const (
	AnswerTypeCanonical AnswerType = "canonical"
	AnswerTypeCommunity AnswerType = "community"
	AnswerTypeSuggested AnswerType = "suggested"
)

type Answer struct {
	ID          uuid.UUID  `json:"id"`
	QuestionID  uuid.UUID  `json:"question_id"`
	Content     string     `json:"content"`
	AuthorID    uuid.UUID  `json:"author_id"`
	AnswerType  AnswerType `json:"answer_type"`
	VoteCount   int        `json:"vote_count"`
	IsAccepted  bool       `json:"is_accepted"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func NewAnswer(questionID, authorID uuid.UUID, content string, answerType AnswerType) *Answer {
	return &Answer{
		ID:         uuid.New(),
		QuestionID: questionID,
		Content:    content,
		AuthorID:   authorID,
		AnswerType: answerType,
		VoteCount:  0,
		IsAccepted: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
