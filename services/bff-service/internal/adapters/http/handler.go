package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BFFHandler struct {
	practiceServiceURL string
	client             *http.Client
}

func NewBFFHandler(practiceServiceURL string) *BFFHandler {
	return &BFFHandler{
		practiceServiceURL: practiceServiceURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (h *BFFHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.POST("/sessions", h.StartSession)
		api.GET("/sessions/:id", h.GetSession)
		api.POST("/sessions/:id/answers", h.SubmitAnswer)
		api.GET("/questions/:id", h.GetQuestion)
	}
}

func (h *BFFHandler) GetQuestion(c *gin.Context) {
	questionID := c.Param("id")
	url := fmt.Sprintf("%s/api/v1/practice/questions/%s", h.practiceServiceURL, questionID)
	h.proxyRequest(c, "GET", url, nil)
}

func (h *BFFHandler) proxyRequest(c *gin.Context, method, url string, body []byte) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copy headers
	req.Header.Set("Content-Type", "application/json")
	for k, v := range c.Request.Header {
		if k != "Host" && k != "Content-Length" {
			req.Header[k] = v
		}
	}

	resp, err := h.client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("Failed to contact service: %v", err)})
		return
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Set status code and content type
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}

func (h *BFFHandler) StartSession(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	url := fmt.Sprintf("%s/api/v1/practice/sessions", h.practiceServiceURL)
	h.proxyRequest(c, "POST", url, body)
}

func (h *BFFHandler) GetSession(c *gin.Context) {
	sessionID := c.Param("id")
	url := fmt.Sprintf("%s/api/v1/practice/sessions/%s", h.practiceServiceURL, sessionID)
	h.proxyRequest(c, "GET", url, nil)
}

func (h *BFFHandler) SubmitAnswer(c *gin.Context) {
	sessionID := c.Param("id")
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	url := fmt.Sprintf("%s/api/v1/practice/sessions/%s/answers", h.practiceServiceURL, sessionID)
	h.proxyRequest(c, "POST", url, body)
}
