# Environment Setup Guide

This project is a microservices-based application using Docker, Go, Python, and React.

## Prerequisites

- [Docker](https://www.docker.com/get-started) and Docker Compose
- [Go](https://go.dev/dl/) (1.23+)
- [Node.js](https://nodejs.org/) (18+)
- [Python](https://www.python.org/) (3.9+)

## Configuration

### Environment Variables

The services rely on environment variables for configuration. You can set these in your shell or in a `.env` file (if supported by the service loader, but Docker Compose handles most of them).

**Key Variables for AI Service:**

You must provide API keys for the AI service to generate feedback and questions.

- `OPENAI_API_KEY`: API key for OpenAI (GPT-4/3.5).
- `GOOGLE_API_KEY`: API key for Google Gemini.
- `GROQ_API_KEY`: API key for Groq (Llama models).
- `LLM_PROVIDER`: Select the provider to use (`openai`, `gemini`, or `groq`). Default is `gemini`.

### Setting up for Docker Compose

1.  Create a `.env` file in the `infra` directory or export variables in your shell before running `docker-compose up`.
    ```bash
    export OPENAI_API_KEY="sk-..."
    export GOOGLE_API_KEY="AIza..."
    export GROQ_API_KEY="gsk_..."
    ```

## Running the Application

### 1. Start Services

Use Docker Compose to start all services (Postgres, Backend, AI Service, Frontend).

```bash
# Navigate to the project root
cd infra
docker-compose up -d --build
```

### 2. Database Migration and Seeding

The database needs to be initialized with tables and initial data.

**Migrations:**
Migrations run automatically via the `migration` service in Docker Compose.

**Seeding Data (Questions):**
If you want to populate the database with sample questions (especially if you don't have AI keys):

```bash
# Run the seed script manually
cd services/practice-service
go mod tidy
go run cmd/seed/main.go
```

### 3. Accessing the Application

- **Frontend:** [http://localhost:3000](http://localhost:3000)
- **Practice Service API:** [http://localhost:8080](http://localhost:8080)
- **BFF Service API:** [http://localhost:8081](http://localhost:8081)
- **AI Service API:** [http://localhost:8000](http://localhost:8000)

## Features

- **AI Evaluation:** Enable "AI Evaluation" in the setup screen to get feedback from the AI.
- **Standard Mode:** Disable "AI Evaluation" to use pre-defined standard answers from the database (requires seeded data).
- **Interview Rounds:** Select from Recruiter, Technical, Algorithms, System Design, and Leadership rounds.
- **Tech Stacks:** Support for Golang, Python, NodeJS, NestJS.
- **Bilingual Support:** Switch between English and Vietnamese.

## Troubleshooting

- **No Questions Found:** Run the seed script (`go run cmd/seed/main.go`) or ensure your AI API keys are valid and the crawler script (`scripts/crawl_data.py`) has run.
- **Database Connection Errors:** Ensure the `postgres` container is healthy (`docker-compose ps`).
