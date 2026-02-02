#!/bin/bash

# Base URL
BASE_URL="http://localhost:8080/api/v1/practice"

# 1. Start Session
echo "Starting session..."
SESSION_RES=$(curl -s -X POST "${BASE_URL}/sessions" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "topic_id": null,
    "level": "Junior"
  }')

echo "Session Response: $SESSION_RES"
SESSION_ID=$(echo $SESSION_RES | jq -r '.session.id')
FIRST_QUESTION_ID=$(echo $SESSION_RES | jq -r '.first_question_id')

if [ "$SESSION_ID" == "null" ]; then
  echo "Failed to start session"
  exit 1
fi

echo "Session ID: $SESSION_ID"
echo "First Question ID: $FIRST_QUESTION_ID"

# 2. Submit Answer
echo "Submitting answer..."
ANSWER_RES=$(curl -s -X POST "${BASE_URL}/sessions/${SESSION_ID}/answers" \
  -H "Content-Type: application/json" \
  -d "{
    \"question_id\": \"$FIRST_QUESTION_ID\",
    \"content\": \"This is a test answer.\"
  }")

echo "Answer Response: $ANSWER_RES"
NEXT_QUESTION_ID=$(echo $ANSWER_RES | jq -r '.next_question_id')

echo "Next Question ID: $NEXT_QUESTION_ID"

if [ "$NEXT_QUESTION_ID" != "null" ] && [ -n "$NEXT_QUESTION_ID" ]; then
  echo "SUCCESS: Received next question ID"
else
  echo "FAILURE: Did not receive next question ID"
  exit 1
fi
