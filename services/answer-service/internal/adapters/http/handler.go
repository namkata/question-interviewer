package http_adapter

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/question-interviewer/answer-service/internal/domain"
	"github.com/question-interviewer/answer-service/internal/ports"
)

type AnswerHandler struct {
	service ports.AnswerService
}

func NewAnswerHandler(service ports.AnswerService) *AnswerHandler {
	return &AnswerHandler{
		service: service,
	}
}

type CreateAnswerRequest struct {
	QuestionID string `json:"question_id" binding:"required,uuid"`
	Content    string `json:"content" binding:"required"`
	AnswerType string `json:"answer_type" binding:"required"`
	AuthorID   string `json:"author_id" binding:"required,uuid"` // TODO: Get from context/auth
}

// CreateAnswer godoc
// @Summary Create a new answer
// @Description Create a new answer for a question
// @Tags answers
// @Accept json
// @Produce json
// @Param answer body CreateAnswerRequest true "Answer Data"
// @Success 201 {object} domain.Answer
// @Router /answers [post]
func (h *AnswerHandler) CreateAnswer(c *gin.Context) {
	var req CreateAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	questionID, err := uuid.Parse(req.QuestionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question_id"})
		return
	}

	authorID, err := uuid.Parse(req.AuthorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author_id"})
		return
	}

	answerType := domain.AnswerType(req.AnswerType)
	if answerType != domain.AnswerTypeCanonical && answerType != domain.AnswerTypeCommunity && answerType != domain.AnswerTypeSuggested {
		// Default to community if invalid
		answerType = domain.AnswerTypeCommunity
	}

	answer, err := h.service.CreateAnswer(c.Request.Context(), questionID, authorID, req.Content, answerType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, answer)
}

// GetAnswer godoc
// @Summary Get an answer by ID
// @Description Get an answer by ID
// @Tags answers
// @Accept json
// @Produce json
// @Param id path string true "Answer ID"
// @Success 200 {object} domain.Answer
// @Router /answers/{id} [get]
func (h *AnswerHandler) GetAnswer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	answer, err := h.service.GetAnswer(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, answer)
}

// ListAnswers godoc
// @Summary List answers for a question
// @Description List answers for a question with pagination
// @Tags answers
// @Accept json
// @Produce json
// @Param question_id query string true "Question ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} domain.Answer
// @Router /answers [get]
func (h *AnswerHandler) ListAnswers(c *gin.Context) {
	questionIDStr := c.Query("question_id")
	if questionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "question_id is required"})
		return
	}

	questionID, err := uuid.Parse(questionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question_id"})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	answers, err := h.service.ListAnswersForQuestion(c.Request.Context(), questionID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, answers)
}

func (h *AnswerHandler) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/answers", h.CreateAnswer)
		v1.GET("/answers/:id", h.GetAnswer)
		v1.GET("/answers", h.ListAnswers)
	}
}
