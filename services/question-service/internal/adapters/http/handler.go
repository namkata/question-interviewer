package http_adapter

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/question-interviewer/question-service/internal/ports"
)

type QuestionHandler struct {
	service ports.QuestionService
}

func NewQuestionHandler(service ports.QuestionService) *QuestionHandler {
	return &QuestionHandler{
		service: service,
	}
}

type CreateQuestionRequest struct {
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Level     string `json:"level" binding:"required"`
	TopicID   string `json:"topic_id" binding:"required,uuid"`
	CreatedBy string `json:"created_by" binding:"required,uuid"` // TODO: Get from context/auth
}

// CreateQuestion godoc
// @Summary Create a new question
// @Description Create a new question
// @Tags questions
// @Accept json
// @Produce json
// @Param question body CreateQuestionRequest true "Question Data"
// @Success 201 {object} domain.Question
// @Router /questions [post]
func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	var req CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	topicID, err := uuid.Parse(req.TopicID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic_id"})
		return
	}

	createdBy, err := uuid.Parse(req.CreatedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid created_by"})
		return
	}

	question, err := h.service.CreateQuestion(c.Request.Context(), req.Title, req.Content, req.Level, topicID, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, question)
}

// GetQuestion godoc
// @Summary Get a question by ID
// @Description Get a question by ID
// @Tags questions
// @Accept json
// @Produce json
// @Param id path string true "Question ID"
// @Success 200 {object} domain.Question
// @Router /questions/{id} [get]
func (h *QuestionHandler) GetQuestion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	question, err := h.service.GetQuestion(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, question)
}

// ListQuestions godoc
// @Summary List questions
// @Description List questions with pagination
// @Tags questions
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} domain.Question
// @Router /questions [get]
func (h *QuestionHandler) ListQuestions(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	questions, err := h.service.ListQuestions(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questions)
}

func (h *QuestionHandler) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/questions", h.CreateQuestion)
		v1.GET("/questions/:id", h.GetQuestion)
		v1.GET("/questions", h.ListQuestions)
	}
}
