# Question Interviewer System

A microservices-based application to help developers practice interview questions with real-time AI feedback.

## 🏗 System Architecture

The system consists of the following microservices:

```mermaid
graph TD
    User[User / Browser] -->|HTTP| Frontend[Frontend (React)]
    Frontend -->|HTTP| BFF[BFF Service (Go)]
    BFF -->|HTTP| Practice[Practice Service (Go)]
    Practice -->|HTTP| AI[AI Service (Python)]
    Practice -->|SQL| DB[(PostgreSQL)]
    AI -->|API| LLM[Google Gemini / OpenAI]
```

- **Frontend**: React + Vite + Tailwind CSS application.
- **BFF Service**: Backend-for-Frontend in Go, handles proxying and aggregation.
- **Practice Service**: Core domain logic for interview sessions and questions.
- **AI Service**: Python FastAPI service using LangChain to evaluate answers via LLMs.
- **PostgreSQL**: Persistent storage for questions, users, and attempts.

## 🚀 Getting Started

### Prerequisites

- [Docker](https://www.docker.com/get-started) installed on your machine.
- [Docker Compose](https://docs.docker.com/compose/install/) installed.

### Environment Setup

1. Create a `.env` file in the `infra` directory (or export variables in your shell) for AI API keys:

```bash
export OPENAI_API_KEY="your-openai-key"
# OR
export GOOGLE_API_KEY="your-gemini-key"
```

### 🛠 Installation & Run

1. Navigate to the infrastructure directory:

```bash
cd infra
```

2. Build and start all services:

```bash
docker-compose up --build -d
```

3. Check status of services:

```bash
docker-compose ps
```

### 🗄 Database Migrations

After starting the containers, you need to apply database migrations to set up the schema and seed data.

You can use the [golang-migrate](https://github.com/golang-migrate/migrate) tool or run the SQL files manually.

**Using migrate tool:**
```bash
migrate -path ../migrations -database "postgresql://user:password@localhost:5432/question_db?sslmode=disable" up
```

**Manual Setup (if no tool):**
You can execute the SQL files in `migrations/` folder against the postgres container.

### 🖥 Usage

- **Frontend**: Open [http://localhost:3000](http://localhost:3000)
- **BFF API**: [http://localhost:8081](http://localhost:8081)
- **Practice Service API**: [http://localhost:8080](http://localhost:8080)
- **AI Service Docs**: [http://localhost:8000/docs](http://localhost:8000/docs)

## 📸 Screenshots

### Practice Session
![Practice Session](docs/images/practice-session.png)
*Interface where users answer questions and receive feedback.*

### AI Feedback
![AI Feedback](docs/images/ai-feedback.png)
*Detailed markdown feedback provided by the AI Service.*

## 🎥 Demo Video

[Watch the Demo Video](docs/videos/demo.mp4)

---
**Note**: Please replace the placeholder images/videos in `docs/` with actual assets after running the application.
