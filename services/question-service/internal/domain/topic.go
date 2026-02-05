package domain

import (
	"github.com/google/uuid"
)

type Topic struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func NewTopic(name, description string) *Topic {
	return &Topic{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
	}
}
