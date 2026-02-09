package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/question-interviewer/practice-service/internal/domain"
	"github.com/question-interviewer/practice-service/internal/ports"
)

type AIClient struct {
	baseURL string
	client  *http.Client
}

func NewAIClient(baseURL string) ports.AIService {
	return &AIClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type EvaluationRequest struct {
	QuestionContent string `json:"question_content"`
	UserAnswer      string `json:"user_answer"`
	CorrectAnswer   string `json:"correct_answer,omitempty"`
	Topic           string `json:"topic"`
	Level           string `json:"level"`
	Language        string `json:"language"`
}

type EvaluationResponse struct {
	Score          int      `json:"score"`
	Feedback       string   `json:"feedback"`
	Suggestions    []string `json:"suggestions"`
	ImprovedAnswer string   `json:"improved_answer"`
}

func (c *AIClient) EvaluateAnswer(ctx context.Context, question, userAnswer, correctAnswer, topic, level, language string) (int, string, []string, string, error) {
	reqBody := EvaluationRequest{
		QuestionContent: question,
		UserAnswer:      userAnswer,
		CorrectAnswer:   correctAnswer,
		Topic:           topic,
		Level:           level,
		Language:        language,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return 0, "", nil, "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/v1/evaluate", bytes.NewBuffer(jsonBody))
	if err != nil {
		return 0, "", nil, "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, "", nil, "", fmt.Errorf("failed to call AI service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, "", nil, "", fmt.Errorf("AI service returned status: %d", resp.StatusCode)
	}

	var evalResp EvaluationResponse
	if err := json.NewDecoder(resp.Body).Decode(&evalResp); err != nil {
		return 0, "", nil, "", fmt.Errorf("failed to decode response: %w", err)
	}

	return evalResp.Score, evalResp.Feedback, evalResp.Suggestions, evalResp.ImprovedAnswer, nil
}

type InterviewSummaryAttempt struct {
	QuestionContent string  `json:"question_content"`
	UserAnswer      string  `json:"user_answer"`
	Score           int     `json:"score"`
	Feedback        string  `json:"feedback"`
	RoundName       *string `json:"round_name,omitempty"`
	RoundIndex      *int    `json:"round_index,omitempty"`
}

type InterviewSummaryRequest struct {
	Role     string                    `json:"role"`
	Language string                    `json:"language"`
	Attempts []InterviewSummaryAttempt `json:"attempts"`
}

type InterviewSummaryResponse struct {
	Strengths    string `json:"strengths"`
	Weaknesses   string `json:"weaknesses"`
	Readiness    string `json:"readiness"`
	OverallScore int    `json:"overall_score"`
}

func (c *AIClient) SummarizeInterview(ctx context.Context, role, language string, attempts []domain.PracticeAttempt) (string, string, string, int, error) {
	summaryAttempts := make([]InterviewSummaryAttempt, 0, len(attempts))
	for _, a := range attempts {
		q := ""
		if a.QuestionContent != nil {
			q = *a.QuestionContent
		}
		summaryAttempts = append(summaryAttempts, InterviewSummaryAttempt{
			QuestionContent: q,
			UserAnswer:      a.UserAnswer,
			Score:           a.Score,
			Feedback:        a.Feedback,
			RoundName:       a.RoundName,
			RoundIndex:      a.RoundIndex,
		})
	}

	reqBody := InterviewSummaryRequest{
		Role:     role,
		Language: language,
		Attempts: summaryAttempts,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", "", 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/v1/summarize", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", "", "", 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", "", "", 0, fmt.Errorf("failed to call AI service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", "", 0, fmt.Errorf("AI service returned status: %d", resp.StatusCode)
	}

	var out InterviewSummaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", "", "", 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return out.Strengths, out.Weaknesses, out.Readiness, out.OverallScore, nil
}
