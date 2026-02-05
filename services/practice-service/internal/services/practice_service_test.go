package services

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/question-interviewer/practice-service/internal/domain"
	"github.com/question-interviewer/practice-service/internal/ports"
)

type fakeRepo struct {
	questionContent string
	questionTopic   string
	questionLevel   string
	correctAnswer   string
	hint            string
}

func (r *fakeRepo) CreateSession(ctx context.Context, session *domain.PracticeSession) error {
	return nil
}
func (r *fakeRepo) GetSession(ctx context.Context, id uuid.UUID) (*domain.PracticeSession, error) {
	return nil, errors.New("not implemented")
}
func (r *fakeRepo) UpdateSession(ctx context.Context, session *domain.PracticeSession) error {
	return nil
}
func (r *fakeRepo) CreateAttempt(ctx context.Context, attempt *domain.PracticeAttempt) error {
	return nil
}
func (r *fakeRepo) GetQuestionSampleCache(ctx context.Context, questionID uuid.UUID) (string, string, []string, string, error) {
	return "", "", nil, "", nil
}
func (r *fakeRepo) UpsertQuestionSampleCache(ctx context.Context, questionID uuid.UUID, sampleAnswer, sampleFeedback string, sampleSuggestions []string, sampleSource string) error {
	return nil
}
func (r *fakeRepo) GetRandomQuestionID(ctx context.Context, topicID *uuid.UUID, level *string, language string, config map[string]interface{}) (uuid.UUID, error) {
	return uuid.Nil, errors.New("not implemented")
}
func (r *fakeRepo) GetQuestionContent(ctx context.Context, questionID uuid.UUID) (string, string, string, string, string, error) {
	return r.questionContent, r.questionTopic, r.questionLevel, r.correctAnswer, r.hint, nil
}
func (r *fakeRepo) CreateQuestion(ctx context.Context, question *domain.Question) error {
	return errors.New("not implemented")
}
func (r *fakeRepo) GetTopicIDByName(ctx context.Context, name string) (uuid.UUID, error) {
	return uuid.Nil, errors.New("not implemented")
}

type fakeAI struct {
	score          int
	feedback       string
	suggestions    []string
	improvedAnswer string
	err            error
}

func (a *fakeAI) EvaluateAnswer(ctx context.Context, question, userAnswer, correctAnswer, topic, level, language string) (int, string, []string, string, error) {
	if a.err != nil {
		return 0, "", nil, "", a.err
	}
	return a.score, a.feedback, a.suggestions, a.improvedAnswer, nil
}

var _ ports.PracticeRepository = (*fakeRepo)(nil)
var _ ports.AIService = (*fakeAI)(nil)

func TestSuggestAnswer_FallbackWhenAIUnavailable(t *testing.T) {
	repo := &fakeRepo{
		questionContent: "Explain Scrum ceremonies.",
		questionTopic:   "Behavioral",
		questionLevel:   "Mid",
		correctAnswer:   "Explain purpose of planning, daily, review, retro, refinement.",
		hint:            "Focus on purpose.",
	}
	ai := &fakeAI{err: errors.New("ai down")}
	svc := NewPracticeService(repo, ai, true)

	score, feedback, suggestions, improved, err := svc.SuggestAnswer(context.Background(), uuid.New(), "", "vi")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if score != 0 {
		t.Fatalf("expected score=0, got %d", score)
	}
	if improved != repo.correctAnswer {
		t.Fatalf("expected improved answer fallback to correct answer")
	}
	if suggestions != nil {
		t.Fatalf("expected nil suggestions on fallback")
	}
	if feedback != "" {
		t.Fatalf("expected empty feedback on fallback")
	}
}
