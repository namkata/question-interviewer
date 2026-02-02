#!/bin/bash

echo "Testing AI Service Health..."
curl -s http://localhost:8000/health | jq .

echo -e "\n\nTesting AI Service Evaluation (using Groq)..."
curl -s -X POST http://localhost:8000/api/v1/evaluate \
  -H "Content-Type: application/json" \
  -d '{
    "question_content": "Explain the difference between concurrency and parallelism.",
    "user_answer": "Concurrency is about dealing with lots of things at once. Parallelism is about doing lots of things at once.",
    "correct_answer": "Concurrency is the composition of independently executing processes, while parallelism is the simultaneous execution of (possibly related) computations. Concurrency is about dealing with many things at once. Parallelism is about doing many things at once.",
    "topic": "Computer Science",
    "level": "Medium"
  }' | jq .
