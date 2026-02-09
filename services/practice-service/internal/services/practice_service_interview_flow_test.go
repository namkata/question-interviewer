package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/question-interviewer/practice-service/internal/domain"
)

type inMemoryRepo struct {
	sessions    map[uuid.UUID]*domain.PracticeSession
	questions   map[uuid.UUID]domain.Question
	topicsByID  map[uuid.UUID]string
	idsByTopic  map[string]uuid.UUID
	attempts    map[uuid.UUID][]*domain.PracticeAttempt
	questionSeq []uuid.UUID
	nextIdx     int
}

func newInMemoryRepo(questionSeq []uuid.UUID, questions map[uuid.UUID]domain.Question, topicNames []string) *inMemoryRepo {
	idsByTopic := map[string]uuid.UUID{}
	topicsByID := map[uuid.UUID]string{}
	for _, name := range topicNames {
		id := uuid.New()
		idsByTopic[name] = id
		topicsByID[id] = name
	}
	return &inMemoryRepo{
		sessions:    map[uuid.UUID]*domain.PracticeSession{},
		questions:   questions,
		topicsByID:  topicsByID,
		idsByTopic:  idsByTopic,
		attempts:    map[uuid.UUID][]*domain.PracticeAttempt{},
		questionSeq: questionSeq,
	}
}

func (r *inMemoryRepo) CreateSession(ctx context.Context, session *domain.PracticeSession) error {
	cp := *session
	r.sessions[session.ID] = &cp
	return nil
}

func (r *inMemoryRepo) GetSession(ctx context.Context, id uuid.UUID) (*domain.PracticeSession, error) {
	s, ok := r.sessions[id]
	if !ok {
		return nil, context.Canceled
	}
	cp := *s
	if s.Config != nil {
		cfg := make(map[string]interface{}, len(s.Config))
		for k, v := range s.Config {
			cfg[k] = v
		}
		cp.Config = cfg
	}
	return &cp, nil
}

func (r *inMemoryRepo) UpdateSession(ctx context.Context, session *domain.PracticeSession) error {
	cp := *session
	r.sessions[session.ID] = &cp
	return nil
}

func (r *inMemoryRepo) CreateAttempt(ctx context.Context, attempt *domain.PracticeAttempt) error {
	cp := *attempt
	r.attempts[attempt.SessionID] = append(r.attempts[attempt.SessionID], &cp)
	return nil
}

func (r *inMemoryRepo) ListAttemptsBySession(ctx context.Context, sessionID uuid.UUID) ([]*domain.PracticeAttempt, error) {
	list := r.attempts[sessionID]
	out := make([]*domain.PracticeAttempt, 0, len(list))
	for _, a := range list {
		cp := *a
		out = append(out, &cp)
	}
	return out, nil
}

func (r *inMemoryRepo) GetQuestionSampleCache(ctx context.Context, questionID uuid.UUID) (string, string, []string, string, error) {
	return "", "", nil, "", nil
}

func (r *inMemoryRepo) UpsertQuestionSampleCache(ctx context.Context, questionID uuid.UUID, sampleAnswer, sampleFeedback string, sampleSuggestions []string, sampleSource string) error {
	return nil
}

func (r *inMemoryRepo) GetRandomQuestionID(ctx context.Context, topicID *uuid.UUID, level *string, language string, config map[string]interface{}) (uuid.UUID, error) {
	if r.nextIdx >= len(r.questionSeq) {
		return uuid.Nil, context.Canceled
	}
	id := r.questionSeq[r.nextIdx]
	r.nextIdx++
	return id, nil
}

func (r *inMemoryRepo) GetQuestionContent(ctx context.Context, questionID uuid.UUID) (string, string, string, string, string, error) {
	q := r.questions[questionID]
	topic := "General"
	if q.TopicID != nil {
		if name, ok := r.topicsByID[*q.TopicID]; ok {
			topic = name
		}
	}
	return q.Content, topic, q.Level, q.CorrectAnswer, q.Hint, nil
}

func (r *inMemoryRepo) CreateQuestion(ctx context.Context, question *domain.Question) error {
	return nil
}

func (r *inMemoryRepo) GetTopicIDByName(ctx context.Context, name string) (uuid.UUID, error) {
	id, ok := r.idsByTopic[name]
	if !ok {
		return uuid.Nil, context.Canceled
	}
	return id, nil
}

func TestInterviewMode_8Rounds_2QuestionsPerRound(t *testing.T) {
	roundNames := getRoundsForRole("BackEnd")

	questions := map[uuid.UUID]domain.Question{}
	seq := make([]uuid.UUID, 0, 16)

	levels := []string{
		"Fresher", "Fresher",
		"Junior", "Junior",
		"Junior", "Junior",
		"Mid", "Mid",
		"Mid", "Mid",
		"Senior", "Senior",
		"Senior", "Senior",
		"Senior", "Senior",
	}

	for i := 0; i < 16; i++ {
		id := uuid.New()
		seq = append(seq, id)
		topicName := roundNames[i/2]
		topicID := uuid.New()
		questions[id] = domain.Question{
			ID:            id,
			Content:       "Q",
			TopicID:       &topicID,
			TopicName:     topicName,
			Level:         levels[i],
			CorrectAnswer: "A",
			Hint:          "",
		}
	}

	repo := newInMemoryRepo(seq, questions, roundNames)
	for _, q := range repo.questions {
		if q.TopicID != nil {
			repo.topicsByID[*q.TopicID] = q.TopicName
		}
		if _, ok := repo.idsByTopic[q.TopicName]; !ok {
			repo.idsByTopic[q.TopicName] = uuid.New()
		}
	}

	svc := NewPracticeService(repo, nil, false)

	userID := uuid.New()
	level := "Fresher"
	session, firstQ, err := svc.StartSession(context.Background(), userID, nil, &level, "en", map[string]interface{}{
		"mode": "interview",
		"role": "BackEnd",
	})
	if err != nil {
		t.Fatalf("StartSession error: %v", err)
	}
	if firstQ == uuid.Nil {
		t.Fatalf("expected first question id")
	}

	currentQ := firstQ
	answers := 0
	for currentQ != uuid.Nil {
		answers++
		attempt, nextQ, err := svc.SubmitAnswer(context.Background(), session.ID, currentQ, "answer", "en", false)
		if err != nil {
			t.Fatalf("SubmitAnswer error: %v", err)
		}
		if attempt.RoundIndex == nil || attempt.RoundName == nil || *attempt.RoundName == "" {
			t.Fatalf("expected round metadata on attempt")
		}
		currentQ = nextQ
		if answers > 20 {
			t.Fatalf("expected to finish within 16 answers")
		}
	}
	if answers != 16 {
		t.Fatalf("expected 16 answers, got %d", answers)
	}

	gotSession, err := svc.GetSession(context.Background(), session.ID)
	if err != nil {
		t.Fatalf("GetSession error: %v", err)
	}
	if gotSession.Status != "completed" {
		t.Fatalf("expected session completed, got %s", gotSession.Status)
	}

	summary, err := svc.GetSessionSummary(context.Background(), session.ID, "en")
	if err != nil {
		t.Fatalf("GetSessionSummary error: %v", err)
	}
	rounds, ok := summary["rounds"].([]map[string]interface{})
	if !ok {
		t.Fatalf("expected rounds in summary")
	}
	if len(rounds) != 8 {
		t.Fatalf("expected 8 rounds in summary, got %d", len(rounds))
	}
}
