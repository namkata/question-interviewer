package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/question-interviewer/practice-service/internal/ports"
)

type PracticeHandler struct {
	service ports.PracticeService
}

func NewPracticeHandler(service ports.PracticeService) *PracticeHandler {
	return &PracticeHandler{
		service: service,
	}
}

type StartSessionRequest struct {
	UserID   string                 `json:"user_id" binding:"required"`
	TopicID  *string                `json:"topic_id"`
	Level    *string                `json:"level"`
	Language string                 `json:"language"`
	Config   map[string]interface{} `json:"config"`
}

type SubmitAnswerRequest struct {
	QuestionID string `json:"question_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Language   string `json:"language"`
	AIEnabled  *bool  `json:"ai_enabled"`
}

type SuggestAnswerRequest struct {
	Content  string `json:"content"`
	Language string `json:"language"`
}

func (h *PracticeHandler) StartSession(c *gin.Context) {
	var req StartSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var topicID *uuid.UUID
	if req.TopicID != nil {
		id, err := uuid.Parse(*req.TopicID)
		if err != nil {
			// If not a valid UUID, try to lookup by name
			id, err = h.service.GetTopicIDByName(c.Request.Context(), *req.TopicID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID or name not found: " + err.Error()})
				return
			}
		}
		topicID = &id
	}

	session, firstQuestionID, err := h.service.StartSession(c.Request.Context(), userID, topicID, req.Level, req.Language, req.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"session":           session,
		"first_question_id": firstQuestionID,
	})
}

func (h *PracticeHandler) SubmitAnswer(c *gin.Context) {
	sessionIDStr := c.Param("id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID format"})
		return
	}

	var req SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	questionID, err := uuid.Parse(req.QuestionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID format"})
		return
	}

	aiEnabled := true
	if req.AIEnabled != nil {
		aiEnabled = *req.AIEnabled
	}

	attempt, nextQuestionID, err := h.service.SubmitAnswer(c.Request.Context(), sessionID, questionID, req.Content, req.Language, aiEnabled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"attempt":          attempt,
		"next_question_id": nextQuestionID,
	})
}

func (h *PracticeHandler) GetSession(c *gin.Context) {
	sessionIDStr := c.Param("id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID format"})
		return
	}

	session, err := h.service.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *PracticeHandler) GetQuestion(c *gin.Context) {
	questionIDStr := c.Param("id")
	questionID, err := uuid.Parse(questionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID format"})
		return
	}

	content, topic, level, correctAnswer, hint, err := h.service.GetQuestion(c.Request.Context(), questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":             questionID,
		"content":        content,
		"topic":          topic,
		"level":          level,
		"correct_answer": correctAnswer,
		"hint":           hint,
	})
}

func (h *PracticeHandler) SuggestAnswer(c *gin.Context) {
	questionIDStr := c.Param("id")
	questionID, err := uuid.Parse(questionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID format"})
		return
	}

	var req SuggestAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	score, feedback, suggestions, improvedAnswer, err := h.service.SuggestAnswer(c.Request.Context(), questionID, req.Content, req.Language)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"score":           score,
		"feedback":        feedback,
		"suggestions":     suggestions,
		"improved_answer": improvedAnswer,
	})
}

func (h *PracticeHandler) SkipRound(c *gin.Context) {
	sessionIDStr := c.Param("id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID format"})
		return
	}

	nextQuestionID, err := h.service.SkipCurrentRound(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"next_question_id": nextQuestionID,
	})
}

func (h *PracticeHandler) GetRandomQuestionForSession(c *gin.Context) {
	sessionIDStr := c.Param("id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID format"})
		return
	}

	topic := c.Query("topic")
	var topicPtr *string
	if topic != "" {
		topicPtr = &topic
	}

	questionID, err := h.service.GetRandomQuestion(c.Request.Context(), sessionID, topicPtr)
	if err != nil {
		if strings.Contains(err.Error(), "topic not found") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"question_id": questionID,
	})
}

func (h *PracticeHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1/practice")
	{
		api.POST("/sessions", h.StartSession)
		api.GET("/sessions/:id", h.GetSession)
		api.POST("/sessions/:id/answers", h.SubmitAnswer)
		api.POST("/sessions/:id/skip", h.SkipRound)
		api.GET("/sessions/:id/questions/random", h.GetRandomQuestionForSession)
		api.GET("/questions/:id", h.GetQuestion)
		api.POST("/questions/:id/suggest", h.SuggestAnswer)
		api.POST("/questions", h.CreateQuestion)
	}
}

// CreateQuestionRequest defines the payload for creating a question
type CreateQuestionRequest struct {
	Content       string `json:"content" binding:"required"`
	Topic         string `json:"topic" binding:"required"`
	Level         string `json:"level" binding:"required"`
	CorrectAnswer string `json:"correct_answer"`
	Hint          string `json:"hint"`
}

func (h *PracticeHandler) CreateQuestion(c *gin.Context) {
	var req CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question, err := h.service.CreateQuestion(c.Request.Context(), req.Content, req.Topic, req.Level, req.CorrectAnswer, req.Hint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, question)
}
