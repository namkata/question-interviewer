package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const apiBaseURL = "http://localhost:8080/api/v1/practice"

type SessionResponse struct {
	Session struct {
		ID     string                 `json:"id"`
		Config map[string]interface{} `json:"config"`
	} `json:"session"`
	FirstQuestionID string `json:"first_question_id"`
}

type QuestionResponse struct {
	ID            string `json:"id"`
	Content       string `json:"content"`
	Topic         string `json:"topic"`
	Hint          string `json:"hint"`
	CorrectAnswer string `json:"correct_answer"`
}

type AnswerResponse struct {
	Attempt struct {
		Score    int    `json:"score"`
		Feedback string `json:"feedback"`
	} `json:"attempt"`
	NextQuestionID string `json:"next_question_id"`
}

func main() {
	// 1. Start Session
	fmt.Println("Starting FrontEnd Session...")
	startReq := map[string]interface{}{
		"user_id":  "123e4567-e89b-12d3-a456-426614174000",
		"language": "en",
		"config": map[string]interface{}{
			"mode":        "interview",
			"role":        "FrontEnd",
			"tech_stacks": []string{"React"},
		},
	}
	jsonData, _ := json.Marshal(startReq)

	resp, err := http.Post(apiBaseURL+"/sessions", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error starting session: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Failed to start session. Status: %d, Body: %s\n", resp.StatusCode, string(body))
		os.Exit(1)
	}

	var sessionResp SessionResponse
	json.NewDecoder(resp.Body).Decode(&sessionResp)
	sessionID := sessionResp.Session.ID
	currentQuestionID := sessionResp.FirstQuestionID

	fmt.Printf("Session Started: %s\n", sessionID)
	fmt.Printf("First Question: %s\n", currentQuestionID)

	questionCount := 0

	// 2. Loop through questions
	for currentQuestionID != "" && currentQuestionID != "00000000-0000-0000-0000-000000000000" && questionCount < 2 {
		questionCount++

		// Get Question Details
		resp, err := http.Get(apiBaseURL + "/questions/" + currentQuestionID)
		if err != nil {
			fmt.Printf("Error getting question: %v\n", err)
			break
		}

		var qResp QuestionResponse
		json.NewDecoder(resp.Body).Decode(&qResp)
		resp.Body.Close()

		fmt.Printf("\n--- Question %d ---\n", questionCount)
		fmt.Printf("Topic: %s\n", qResp.Topic)
		fmt.Printf("Content: %s\n", qResp.Content)
		if qResp.Hint != "" {
			fmt.Printf("Hint: %s\n", qResp.Hint)
		} else {
			fmt.Printf("Hint: [MISSING]\n")
		}

		// Submit Answer
		answerReq := map[string]interface{}{
			"question_id": currentQuestionID,
			"content":     "I have experience with this. Specifically...",
			"language":    "en",
			"ai_enabled":  true,
		}
		jsonData, _ = json.Marshal(answerReq)
		resp, err = http.Post(fmt.Sprintf("%s/sessions/%s/answers", apiBaseURL, sessionID), "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error submitting answer: %v\n", err)
			break
		}

		var ansResp AnswerResponse
		json.NewDecoder(resp.Body).Decode(&ansResp)
		resp.Body.Close()

		fmt.Printf("Score: %d\n", ansResp.Attempt.Score)
		fmt.Printf("Feedback: %s\n", ansResp.Attempt.Feedback)
		fmt.Printf("Next Question ID: %s\n", ansResp.NextQuestionID)

		currentQuestionID = ansResp.NextQuestionID
	}

	fmt.Println("\nSession Finished!")
}
