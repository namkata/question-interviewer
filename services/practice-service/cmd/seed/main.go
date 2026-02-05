package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type TopicSeed struct {
	Name        string
	Description string
}

type QuestionSeed struct {
	Topic         string
	Title         string
	Content       string
	Level         string
	CorrectAnswer string
	Language      string
	Role          string
	Hint          string
}

type BilingualText struct {
	En string
	Vi string
}

type QuestionTemplate struct {
	Topic         string
	Title         BilingualText
	Content       BilingualText
	Level         string
	Role          string
	CorrectAnswer BilingualText
	Hint          BilingualText
}

func main() {
	dsn := "host=localhost port=5432 user=user password=password dbname=question_db sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open db: %v", err)
	}
	defer db.Close()

	seedUserID := "123e4567-e89b-12d3-a456-426614174000"

	CleanTables(db)
	EnsureSeedUser(db, seedUserID)
	EnsureTopics(db, BuildTopics())

	topicIDByName := GetTopicIDByName(db)
	questions := BuildQuestions()
	ValidateQuestions(questions)
	SeedQuestions(db, seedUserID, topicIDByName, questions)
}

func CleanTables(db *sql.DB) {
	tables := []string{
		"votes",
		"bookmarks",
		"practice_attempts",
		"practice_sessions",
		"answers",
		"questions",
		"topics",
		"users",
		"crawled_questions",
		"sessions",
	}
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			log.Printf("Note: Failed to delete from %s (might not exist or other error): %v", table, err)
			continue
		}
		fmt.Printf("Cleaned up table: %s\n", table)
	}
}

func EnsureSeedUser(db *sql.DB, seedUserID string) {
	_, err := db.Exec(
		"INSERT INTO users (id, email, username, password_hash, role) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING",
		seedUserID,
		"demo@example.com",
		"demo_user",
		"hashed_password",
		"user",
	)
	if err != nil {
		log.Fatalf("Failed to insert seed user: %v", err)
	}
}

func EnsureTopics(db *sql.DB, topics []TopicSeed) {
	for _, t := range topics {
		_, err := db.Exec(
			"INSERT INTO topics (id, name, description) VALUES (uuid_generate_v4(), $1, $2) ON CONFLICT (name) DO NOTHING",
			t.Name,
			t.Description,
		)
		if err != nil {
			log.Fatalf("Failed to ensure topic %s: %v", t.Name, err)
		}
	}
}

func GetTopicIDByName(db *sql.DB) map[string]string {
	rows, err := db.Query("SELECT id, name FROM topics")
	if err != nil {
		log.Fatalf("Failed to query topics: %v", err)
	}
	defer rows.Close()

	result := map[string]string{}
	for rows.Next() {
		var id string
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatalf("Failed to scan topic: %v", err)
		}
		result[name] = id
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Failed to iterate topics: %v", err)
	}

	return result
}

func ValidateQuestions(questions []QuestionSeed) {
	seen := map[string]struct{}{}
	for _, q := range questions {
		key := strings.ToLower(strings.TrimSpace(q.Language)) + "::" + strings.TrimSpace(q.Content)
		if _, ok := seen[key]; ok {
			log.Fatalf("Duplicate question content detected (language=%s): %s", q.Language, q.Content)
		}
		seen[key] = struct{}{}
	}
}

func SeedQuestions(db *sql.DB, seedUserID string, topicIDByName map[string]string, questions []QuestionSeed) {
	for _, q := range questions {
		topicID, ok := topicIDByName[q.Topic]
		if !ok {
			log.Fatalf("Topic not found for question title=%s topic=%s", q.Title, q.Topic)
		}

		_, err := db.Exec(
			`INSERT INTO questions (id, title, content, level, topic_id, created_by, status, correct_answer, language, role, hint)
			 VALUES (uuid_generate_v4(), $1, $2, $3, $4, $5, 'published', $6, $7, $8, $9)`,
			q.Title,
			q.Content,
			q.Level,
			topicID,
			seedUserID,
			q.CorrectAnswer,
			q.Language,
			q.Role,
			q.Hint,
		)
		if err != nil {
			log.Fatalf("Failed to insert question %s: %v", q.Title, err)
		}
	}
}

func BuildTopics() []TopicSeed {
	return []TopicSeed{
		{Name: "CV Screening", Description: "Round 1: CV screening and project background"},
		{Name: "Golang", Description: "Round 2: Core language (Go)"},
		{Name: "NodeJS", Description: "Round 2: Core backend runtime (Node.js)"},
		{Name: "Python", Description: "Round 2: Core language (Python)"},
		{Name: "Java", Description: "Round 2: Core language (Java)"},
		{Name: "C#", Description: "Round 2: Core language (C#/.NET)"},
		{Name: "JavaScript", Description: "Round 2: Core language (JavaScript/TypeScript)"},
		{Name: "React", Description: "Round 2: Frontend framework (React)"},
		{Name: "Vue", Description: "Round 2: Frontend framework (Vue)"},
		{Name: "Angular", Description: "Round 2: Frontend framework (Angular)"},
		{Name: "Django", Description: "Round 2: Backend framework (Django)"},
		{Name: "Spring Boot", Description: "Round 2: Backend framework (Spring Boot)"},
		{Name: "Data Layer", Description: "Round 2: Data layer fundamentals (PostgreSQL/Redis/MongoDB)"},
		{Name: "Docker", Description: "Round 2: Containerization fundamentals (Docker)"},
		{Name: "Kubernetes", Description: "Round 2: Orchestration fundamentals (Kubernetes)"},
		{Name: "AWS", Description: "Round 2: Cloud fundamentals (AWS)"},
		{Name: "Database", Description: "Round 3: Database modeling, SQL, performance"},
		{Name: "System Design", Description: "Round 4: System design and architecture"},
		{Name: "Algorithms", Description: "Round 5: Coding and algorithms"},
		{Name: "Testing", Description: "Round 6: Testing and quality"},
		{Name: "DevOps", Description: "Round 7: DevOps and infrastructure"},
		{Name: "Behavioral", Description: "Round 8: Behavioral and communication"},
	}
}

func BuildQuestions() []QuestionSeed {
	templates := BuildQuestionTemplates()
	questions := ExpandTemplatesToQuestions(templates)
	fmt.Printf("Prepared %d questions for seeding\n", len(questions))
	return questions
}

type Concept struct {
	Name      BilingualText
	KeyPoints BilingualText
	Hint      BilingualText
}

type Format struct {
	Title         BilingualText
	Content       BilingualText
	CorrectAnswer BilingualText
	Hint          BilingualText
	Level         string
}

func ExpandTemplatesToQuestions(templates []QuestionTemplate) []QuestionSeed {
	questions := make([]QuestionSeed, 0, len(templates)*2)
	for _, t := range templates {
		questions = append(questions, QuestionSeed{
			Topic:         t.Topic,
			Title:         t.Title.En,
			Content:       t.Content.En,
			Level:         t.Level,
			CorrectAnswer: t.CorrectAnswer.En,
			Language:      "en",
			Role:          t.Role,
			Hint:          t.Hint.En,
		})
		questions = append(questions, QuestionSeed{
			Topic:         t.Topic,
			Title:         t.Title.Vi,
			Content:       t.Content.Vi,
			Level:         t.Level,
			CorrectAnswer: t.CorrectAnswer.Vi,
			Language:      "vi",
			Role:          t.Role,
			Hint:          t.Hint.Vi,
		})
	}
	return questions
}

func BuildQuestionTemplates() []QuestionTemplate {
	templates := make([]QuestionTemplate, 0, 700)
	templates = append(templates, BuildCvTemplates()...)
	templates = append(templates, BuildBehavioralTemplates()...)
	templates = append(templates, BuildAlgorithmTemplates()...)
	templates = append(templates, BuildSystemDesignTemplates()...)
	templates = append(templates, BuildDatabaseTemplates()...)
	templates = append(templates, BuildTestingTemplates()...)
	templates = append(templates, BuildDevOpsTemplates()...)
	templates = append(templates, BuildDataLayerTemplates()...)
	templates = append(templates, BuildGolangTemplates()...)
	templates = append(templates, BuildNodeJSTemplates()...)
	templates = append(templates, BuildJavaScriptTemplates()...)
	templates = append(templates, BuildReactTemplates()...)
	templates = append(templates, BuildVueTemplates()...)
	templates = append(templates, BuildPythonTemplates()...)
	templates = append(templates, BuildJavaTemplates()...)
	templates = append(templates, BuildCSharpTemplates()...)
	templates = append(templates, BuildAngularTemplates()...)
	templates = append(templates, BuildDjangoTemplates()...)
	templates = append(templates, BuildSpringBootTemplates()...)
	templates = append(templates, BuildDockerTemplates()...)
	templates = append(templates, BuildKubernetesTemplates()...)
	templates = append(templates, BuildAwsTemplates()...)

	totalQuestions := len(templates) * 2
	if totalQuestions < 1000 || totalQuestions > 1500 {
		log.Fatalf("Generated question count out of range: %d (expected 1000-1500)", totalQuestions)
	}

	return templates
}

func Render(text string, data map[string]string) string {
	out := text
	for k, v := range data {
		out = strings.ReplaceAll(out, "{"+k+"}", v)
	}
	return out
}

func RenderBilingual(text BilingualText, dataEn map[string]string, dataVi map[string]string) BilingualText {
	return BilingualText{
		En: Render(text.En, dataEn),
		Vi: Render(text.Vi, dataVi),
	}
}

func BuildFromConcepts(topic string, role string, concepts []Concept, formats []Format) []QuestionTemplate {
	templates := make([]QuestionTemplate, 0, len(concepts)*len(formats))
	for _, c := range concepts {
		for _, f := range formats {
			dataEn := map[string]string{
				"Concept":   c.Name.En,
				"KeyPoints": c.KeyPoints.En,
			}
			dataVi := map[string]string{
				"Concept":   c.Name.Vi,
				"KeyPoints": c.KeyPoints.Vi,
			}

			templates = append(templates, QuestionTemplate{
				Topic:         topic,
				Title:         RenderBilingual(f.Title, dataEn, dataVi),
				Content:       RenderBilingual(f.Content, dataEn, dataVi),
				Level:         f.Level,
				Role:          role,
				CorrectAnswer: RenderBilingual(f.CorrectAnswer, dataEn, dataVi),
				Hint: RenderBilingual(BilingualText{
					En: f.Hint.En + " " + c.Hint.En,
					Vi: f.Hint.Vi + " " + c.Hint.Vi,
				}, dataEn, dataVi),
			})
		}
	}
	return templates
}

func BuildCvTemplates() []QuestionTemplate {
	concepts := []Concept{
		{Name: BilingualText{En: "Project impact", Vi: "Tác động dự án"}, KeyPoints: BilingualText{En: "scope, ownership, metrics, trade-offs", Vi: "phạm vi, ownership, số liệu, trade-off"}, Hint: BilingualText{En: "Use STAR and numbers.", Vi: "Dùng STAR và số liệu."}},
		{Name: BilingualText{En: "Trade-off decision", Vi: "Quyết định đánh đổi"}, KeyPoints: BilingualText{En: "constraints, options, decision criteria, outcome", Vi: "ràng buộc, phương án, tiêu chí, kết quả"}, Hint: BilingualText{En: "Compare at least two options.", Vi: "So sánh ít nhất hai phương án."}},
		{Name: BilingualText{En: "Hardest bug", Vi: "Bug khó nhất"}, KeyPoints: BilingualText{En: "repro, hypotheses, data, fix, prevention", Vi: "tái hiện, giả thuyết, dữ liệu, fix, phòng ngừa"}, Hint: BilingualText{En: "Explain how you narrowed scope.", Vi: "Nêu cách khoanh vùng."}},
		{Name: BilingualText{En: "Ownership", Vi: "Ownership"}, KeyPoints: BilingualText{En: "accountability, communication, delivery", Vi: "trách nhiệm, giao tiếp, bàn giao"}, Hint: BilingualText{En: "Avoid blame; focus on improvements.", Vi: "Tránh đổ lỗi; tập trung cải tiến."}},
		{Name: BilingualText{En: "Production incident", Vi: "Sự cố production"}, KeyPoints: BilingualText{En: "timeline, mitigation, postmortem, follow-up", Vi: "timeline, giảm thiểu, postmortem, follow-up"}, Hint: BilingualText{En: "Show prevention actions.", Vi: "Nêu hành động phòng ngừa."}},
		{Name: BilingualText{En: "Estimations", Vi: "Ước lượng"}, KeyPoints: BilingualText{En: "breakdown, risks, buffers, communication", Vi: "chia nhỏ, rủi ro, buffer, giao tiếp"}, Hint: BilingualText{En: "Explain confidence and risks.", Vi: "Nêu độ tin cậy và rủi ro."}},
		{Name: BilingualText{En: "Learning process", Vi: "Quy trình học"}, KeyPoints: BilingualText{En: "sources, experiments, adoption criteria", Vi: "nguồn, thử nghiệm, tiêu chí áp dụng"}, Hint: BilingualText{En: "Show how you evaluate ROI.", Vi: "Nêu cách đánh giá ROI."}},
		{Name: BilingualText{En: "Team collaboration", Vi: "Hợp tác nhóm"}, KeyPoints: BilingualText{En: "alignment, handoffs, conflict handling", Vi: "đồng thuận, bàn giao, xử lý xung đột"}, Hint: BilingualText{En: "Mention rituals and tools.", Vi: "Nêu thói quen và công cụ."}},
		{Name: BilingualText{En: "Code quality", Vi: "Chất lượng code"}, KeyPoints: BilingualText{En: "reviews, standards, tests, refactoring", Vi: "review, tiêu chuẩn, test, refactor"}, Hint: BilingualText{En: "Use concrete examples.", Vi: "Dùng ví dụ cụ thể."}},
		{Name: BilingualText{En: "Communication with stakeholders", Vi: "Giao tiếp với stakeholder"}, KeyPoints: BilingualText{En: "status, risks, options, escalation", Vi: "tiến độ, rủi ro, phương án, escalate"}, Hint: BilingualText{En: "Communicate early with options.", Vi: "Thông báo sớm kèm phương án."}},
		{Name: BilingualText{En: "Agile/Scrum experience", Vi: "Kinh nghiệm Agile/Scrum"}, KeyPoints: BilingualText{En: "roles, ceremonies, artifacts, outcomes", Vi: "vai trò, nghi thức, artifact, kết quả"}, Hint: BilingualText{En: "Be concrete: team size, sprint length, outcomes.", Vi: "Cụ thể: size team, độ dài sprint, kết quả."}},
		{Name: BilingualText{En: "Scrum ceremonies", Vi: "Các buổi Scrum"}, KeyPoints: BilingualText{En: "planning, daily, review, retrospective, refinement", Vi: "planning, daily, review, retrospective, refinement"}, Hint: BilingualText{En: "Explain purpose, not just definitions.", Vi: "Nêu mục đích, không chỉ định nghĩa."}},
		{Name: BilingualText{En: "Definition of Done", Vi: "Definition of Done"}, KeyPoints: BilingualText{En: "quality bar, testing, reviews, deployment readiness", Vi: "chuẩn chất lượng, test, review, sẵn sàng deploy"}, Hint: BilingualText{En: "Link DoD to risk reduction.", Vi: "Gắn DoD với giảm rủi ro."}},
		{Name: BilingualText{En: "Sprint planning and estimation", Vi: "Sprint planning và ước lượng"}, KeyPoints: BilingualText{En: "story points, capacity, uncertainties, slicing", Vi: "story points, capacity, bất định, chia nhỏ"}, Hint: BilingualText{En: "Explain how you avoid over-commit.", Vi: "Nêu cách tránh over-commit."}},
		{Name: BilingualText{En: "Backlog refinement", Vi: "Backlog refinement"}, KeyPoints: BilingualText{En: "acceptance criteria, dependencies, readiness, ambiguity", Vi: "acceptance criteria, phụ thuộc, sẵn sàng, mơ hồ"}, Hint: BilingualText{En: "Talk about how you reduce ambiguity.", Vi: "Nêu cách giảm mơ hồ."}},
		{Name: BilingualText{En: "Working with Product Owner", Vi: "Làm việc với Product Owner"}, KeyPoints: BilingualText{En: "priorities, trade-offs, stakeholder alignment", Vi: "ưu tiên, trade-off, align stakeholder"}, Hint: BilingualText{En: "Explain how you handle changing priorities.", Vi: "Nêu cách xử lý ưu tiên thay đổi."}},
		{Name: BilingualText{En: "Kanban and flow", Vi: "Kanban và flow"}, KeyPoints: BilingualText{En: "WIP limits, cycle time, throughput, pull system", Vi: "giới hạn WIP, cycle time, throughput, pull"}, Hint: BilingualText{En: "Use metrics (cycle time/throughput).", Vi: "Dùng metrics (cycle time/throughput)."}},
		{Name: BilingualText{En: "Agile values and principles", Vi: "Giá trị và nguyên tắc Agile"}, KeyPoints: BilingualText{En: "Agile manifesto, customer value, feedback loops", Vi: "Agile manifesto, giá trị khách hàng, vòng phản hồi"}, Hint: BilingualText{En: "Explain how it changes daily work.", Vi: "Nêu ảnh hưởng đến công việc hằng ngày."}},
		{Name: BilingualText{En: "User stories (INVEST)", Vi: "User story (INVEST)"}, KeyPoints: BilingualText{En: "INVEST, slicing, acceptance criteria", Vi: "INVEST, chia nhỏ, acceptance criteria"}, Hint: BilingualText{En: "Give an example story and AC.", Vi: "Cho ví dụ story và AC."}},
		{Name: BilingualText{En: "Definition of Ready", Vi: "Definition of Ready"}, KeyPoints: BilingualText{En: "ready criteria, dependencies, testability", Vi: "tiêu chí ready, phụ thuộc, testability"}, Hint: BilingualText{En: "Explain why DoR helps planning.", Vi: "Nêu vì sao DoR giúp planning."}},
		{Name: BilingualText{En: "Scrum anti-patterns", Vi: "Anti-pattern trong Scrum"}, KeyPoints: BilingualText{En: "fake agility, meeting-only Scrum, hidden work", Vi: "agile giả, Scrum chỉ họp, hidden work"}, Hint: BilingualText{En: "Describe a symptom and a fix.", Vi: "Nêu triệu chứng và cách sửa."}},
		{Name: BilingualText{En: "Handling sprint interruptions", Vi: "Xử lý gián đoạn sprint"}, KeyPoints: BilingualText{En: "interrupt policy, on-call, re-planning, transparency", Vi: "policy gián đoạn, on-call, re-plan, minh bạch"}, Hint: BilingualText{En: "Talk about trade-offs and communication.", Vi: "Nêu trade-off và giao tiếp."}},
	}

	formats := []Format{
		{
			Level: "Junior",
			Title: BilingualText{En: "{Concept}", Vi: "{Concept}"},
			Content: BilingualText{
				En: "Tell me about {Concept}. What did you do and what was the result?",
				Vi: "Kể về {Concept}. Bạn đã làm gì và kết quả ra sao?",
			},
			CorrectAnswer: BilingualText{
				En: "Explain: {KeyPoints}. Include measurable outcomes when possible.",
				Vi: "Trình bày: {KeyPoints}. Nếu có thể hãy có số liệu đo lường.",
			},
			Hint: BilingualText{En: "Be specific.", Vi: "Chọn ví dụ cụ thể."},
		},
		{
			Level: "Mid",
			Title: BilingualText{En: "{Concept} deep dive", Vi: "Đào sâu {Concept}"},
			Content: BilingualText{
				En: "For {Concept}, what options did you consider and how did you decide?",
				Vi: "Với {Concept}, bạn cân nhắc những phương án nào và quyết định ra sao?",
			},
			CorrectAnswer: BilingualText{
				En: "Compare options and trade-offs: {KeyPoints}. Explain your decision criteria.",
				Vi: "So sánh phương án và trade-off: {KeyPoints}. Nêu tiêu chí quyết định.",
			},
			Hint: BilingualText{En: "Show constraints and trade-offs.", Vi: "Nêu ràng buộc và trade-off."},
		},
	}

	return BuildFromConcepts("CV Screening", "Any", concepts, formats)
}

func BuildBehavioralTemplates() []QuestionTemplate {
	concepts := []Concept{
		{Name: BilingualText{En: "Conflict resolution", Vi: "Giải quyết xung đột"}, KeyPoints: BilingualText{En: "empathy, alignment, outcome", Vi: "đồng cảm, đồng thuận, kết quả"}, Hint: BilingualText{En: "Focus on process and alignment.", Vi: "Tập trung vào quy trình và đồng thuận."}},
		{Name: BilingualText{En: "Receiving feedback", Vi: "Nhận phản hồi"}, KeyPoints: BilingualText{En: "listen, clarify, act, follow up", Vi: "lắng nghe, làm rõ, hành động, follow up"}, Hint: BilingualText{En: "Show what changed afterward.", Vi: "Nêu bạn đã thay đổi gì."}},
		{Name: BilingualText{En: "Handling ambiguity", Vi: "Xử lý mơ hồ"}, KeyPoints: BilingualText{En: "clarify, assumptions, iterate, communicate", Vi: "làm rõ, giả định, lặp, giao tiếp"}, Hint: BilingualText{En: "State assumptions explicitly.", Vi: "Nêu rõ giả định."}},
		{Name: BilingualText{En: "Mentoring", Vi: "Mentoring"}, KeyPoints: BilingualText{En: "goals, pairing, review, growth plan", Vi: "mục tiêu, pairing, review, kế hoạch phát triển"}, Hint: BilingualText{En: "Use a growth plan and feedback loop.", Vi: "Có growth plan và vòng phản hồi."}},
		{Name: BilingualText{En: "Ownership under pressure", Vi: "Ownership khi áp lực"}, KeyPoints: BilingualText{En: "prioritization, communication, escalation", Vi: "ưu tiên, giao tiếp, escalate"}, Hint: BilingualText{En: "Communicate early.", Vi: "Thông báo sớm."}},
		{Name: BilingualText{En: "Working in Scrum teams", Vi: "Làm việc trong Scrum team"}, KeyPoints: BilingualText{En: "collaboration, transparency, commitments", Vi: "hợp tác, minh bạch, cam kết"}, Hint: BilingualText{En: "Describe your role and responsibilities.", Vi: "Nêu vai trò và trách nhiệm."}},
		{Name: BilingualText{En: "Retrospectives and improvements", Vi: "Retrospective và cải tiến"}, KeyPoints: BilingualText{En: "root cause, action items, follow-through", Vi: "nguyên nhân gốc, action item, follow-through"}, Hint: BilingualText{En: "Show an improvement that stuck.", Vi: "Nêu cải tiến thực sự duy trì được."}},
		{Name: BilingualText{En: "Handling changing requirements", Vi: "Xử lý thay đổi yêu cầu"}, KeyPoints: BilingualText{En: "communication, re-prioritization, scope control", Vi: "giao tiếp, ưu tiên lại, kiểm soát phạm vi"}, Hint: BilingualText{En: "Talk about trade-offs and impacts.", Vi: "Nêu trade-off và tác động."}},
		{Name: BilingualText{En: "Tech debt management", Vi: "Quản lý tech debt"}, KeyPoints: BilingualText{En: "prioritization, budgeting, refactoring strategy", Vi: "ưu tiên, phân bổ, chiến lược refactor"}, Hint: BilingualText{En: "Connect to risk and velocity.", Vi: "Gắn với rủi ro và velocity."}},
		{Name: BilingualText{En: "Cross-functional collaboration", Vi: "Phối hợp cross-functional"}, KeyPoints: BilingualText{En: "PM/QA/Design alignment, handoffs, shared goals", Vi: "align PM/QA/Design, bàn giao, mục tiêu chung"}, Hint: BilingualText{En: "Focus on outcomes and communication.", Vi: "Tập trung outcome và giao tiếp."}},
		{Name: BilingualText{En: "Facilitating meetings", Vi: "Facilitate cuộc họp"}, KeyPoints: BilingualText{En: "agenda, timeboxing, decisions, next steps", Vi: "agenda, timebox, quyết định, bước tiếp theo"}, Hint: BilingualText{En: "Show how you keep meetings effective.", Vi: "Nêu cách làm meeting hiệu quả."}},
		{Name: BilingualText{En: "Flow metrics", Vi: "Metric flow"}, KeyPoints: BilingualText{En: "lead time, cycle time, throughput, WIP", Vi: "lead time, cycle time, throughput, WIP"}, Hint: BilingualText{En: "Explain what you measure and why.", Vi: "Nêu bạn đo gì và vì sao."}},
		{Name: BilingualText{En: "Scrum roles and accountability", Vi: "Vai trò Scrum và trách nhiệm"}, KeyPoints: BilingualText{En: "PO/SM/Dev responsibilities, collaboration", Vi: "trách nhiệm PO/SM/Dev, phối hợp"}, Hint: BilingualText{En: "Avoid role confusion.", Vi: "Tránh nhầm vai trò."}},
		{Name: BilingualText{En: "Planning Poker and estimation bias", Vi: "Planning Poker và bias ước lượng"}, KeyPoints: BilingualText{En: "anchoring, consensus, uncertainty", Vi: "anchoring, đồng thuận, bất định"}, Hint: BilingualText{En: "Explain how you reduce bias.", Vi: "Nêu cách giảm bias."}},
		{Name: BilingualText{En: "Velocity misuse", Vi: "Lạm dụng velocity"}, KeyPoints: BilingualText{En: "gaming metrics, comparing teams, outcome vs output", Vi: "game metric, so sánh team, outcome vs output"}, Hint: BilingualText{En: "Explain what you use velocity for (and not).", Vi: "Nêu velocity dùng để làm gì (và không)."}},
		{Name: BilingualText{En: "Outcome vs output", Vi: "Outcome vs output"}, KeyPoints: BilingualText{En: "product value, impact metrics, learning loops", Vi: "giá trị sản phẩm, metric tác động, vòng học"}, Hint: BilingualText{En: "Give an example of an outcome metric.", Vi: "Cho ví dụ outcome metric."}},
		{Name: BilingualText{En: "ScrumBut / fake agility", Vi: "ScrumBut / Agile giả"}, KeyPoints: BilingualText{En: "symptoms, root causes, change strategy", Vi: "triệu chứng, nguyên nhân gốc, chiến lược đổi"}, Hint: BilingualText{En: "Describe how you influence change.", Vi: "Nêu cách bạn tác động thay đổi."}},
		{Name: BilingualText{En: "Backlog refinement facilitation", Vi: "Facilitate backlog refinement"}, KeyPoints: BilingualText{En: "questions, slicing, acceptance criteria, dependencies", Vi: "câu hỏi, chia nhỏ, acceptance criteria, phụ thuộc"}, Hint: BilingualText{En: "Show how you make items ready.", Vi: "Nêu cách làm item ready."}},
		{Name: BilingualText{En: "Managing WIP and bottlenecks", Vi: "Quản lý WIP và bottleneck"}, KeyPoints: BilingualText{En: "WIP limits, queues, swarm, pull policies", Vi: "giới hạn WIP, hàng đợi, swarm, pull"}, Hint: BilingualText{En: "Discuss how you improve flow.", Vi: "Nêu cách cải thiện flow."}},
		{Name: BilingualText{En: "Handling urgent production work", Vi: "Xử lý việc gấp production"}, KeyPoints: BilingualText{En: "triage, interrupt policy, transparency, learning", Vi: "triage, policy gián đoạn, minh bạch, học hỏi"}, Hint: BilingualText{En: "Explain how you protect focus.", Vi: "Nêu cách bảo vệ focus."}},
		{Name: BilingualText{En: "Balancing tech debt and features", Vi: "Cân bằng tech debt và feature"}, KeyPoints: BilingualText{En: "risk framing, roadmap, incremental refactors", Vi: "đóng khung rủi ro, roadmap, refactor tăng dần"}, Hint: BilingualText{En: "Tie to reliability and delivery.", Vi: "Gắn với reliability và delivery."}},
		{Name: BilingualText{En: "Improvement follow-through", Vi: "Theo dõi cải tiến"}, KeyPoints: BilingualText{En: "action items, ownership, review cadence", Vi: "action item, ownership, nhịp review"}, Hint: BilingualText{En: "Show how you ensure changes stick.", Vi: "Nêu cách đảm bảo cải tiến duy trì."}},
	}

	formats := []Format{
		{Level: "Junior", Title: BilingualText{En: "{Concept}", Vi: "{Concept}"}, Content: BilingualText{En: "Tell me about a time you faced {Concept}. What happened and what did you learn?", Vi: "Kể về một lần bạn gặp {Concept}. Chuyện gì xảy ra và bạn học được gì?"}, CorrectAnswer: BilingualText{En: "Describe context and actions. Highlight: {KeyPoints}.", Vi: "Nêu bối cảnh và hành động. Nhấn mạnh: {KeyPoints}."}, Hint: BilingualText{En: "Use a real story.", Vi: "Chọn câu chuyện thật."}},
		{Level: "Mid", Title: BilingualText{En: "{Concept} in teams", Vi: "{Concept} trong teamwork"}, Content: BilingualText{En: "How do you approach {Concept} when the team disagrees?", Vi: "Bạn tiếp cận {Concept} thế nào khi team bất đồng?"}, CorrectAnswer: BilingualText{En: "Explain approach and communication: {KeyPoints}.", Vi: "Giải thích cách làm và giao tiếp: {KeyPoints}."}, Hint: BilingualText{En: "Show alignment and trade-offs.", Vi: "Nêu cách align và trade-off."}},
		{Level: "Senior", Title: BilingualText{En: "Coaching on {Concept}", Vi: "Coaching về {Concept}"}, Content: BilingualText{En: "How would you coach a junior engineer on {Concept}?", Vi: "Bạn sẽ coach junior engineer về {Concept} như thế nào?"}, CorrectAnswer: BilingualText{En: "Outline a mentoring plan: {KeyPoints}.", Vi: "Phác thảo plan mentoring: {KeyPoints}."}, Hint: BilingualText{En: "Make it actionable.", Vi: "Hành động cụ thể."}},
		{Level: "Mid", Title: BilingualText{En: "Measuring {Concept}", Vi: "Đo lường {Concept}"}, Content: BilingualText{En: "How do you know your approach to {Concept} is working?", Vi: "Làm sao bạn biết cách làm {Concept} của mình hiệu quả?"}, CorrectAnswer: BilingualText{En: "Define signals and outcomes: {KeyPoints}.", Vi: "Nêu tín hiệu và kết quả: {KeyPoints}."}, Hint: BilingualText{En: "Use outcomes, not opinions.", Vi: "Dùng kết quả, không chỉ cảm tính."}},
		{Level: "Senior", Title: BilingualText{En: "Risks in {Concept}", Vi: "Rủi ro trong {Concept}"}, Content: BilingualText{En: "What risks do you watch for in {Concept} and how do you mitigate them?", Vi: "Bạn theo dõi rủi ro gì trong {Concept} và giảm thiểu ra sao?"}, CorrectAnswer: BilingualText{En: "List risks and mitigations: {KeyPoints}.", Vi: "Liệt kê rủi ro và giảm thiểu: {KeyPoints}."}, Hint: BilingualText{En: "Think people and process.", Vi: "Nghĩ về con người và quy trình."}},
	}

	return BuildFromConcepts("Behavioral", "Any", concepts, formats)
}

func BuildAlgorithmTemplates() []QuestionTemplate {
	templates := []QuestionTemplate{
		{Topic: "Algorithms", Level: "Fresher", Role: "Any", Title: BilingualText{En: "Big-O basics", Vi: "Big-O cơ bản"}, Content: BilingualText{En: "Explain time complexity and space complexity with examples.", Vi: "Giải thích độ phức tạp thời gian và bộ nhớ với ví dụ."}, CorrectAnswer: BilingualText{En: "Use Big-O; mention common patterns O(1), O(n), O(n log n), O(n^2).", Vi: "Dùng Big-O; nêu các dạng O(1), O(n), O(n log n), O(n^2)."}, Hint: BilingualText{En: "Start from loops and recursion.", Vi: "Bắt đầu từ vòng lặp và đệ quy."}},
		{Topic: "Algorithms", Level: "Junior", Role: "Any", Title: BilingualText{En: "Two pointers", Vi: "Two pointers"}, Content: BilingualText{En: "Given a sorted array, find two numbers that sum to a target. Explain approach and complexity.", Vi: "Với mảng đã sắp xếp, tìm hai số có tổng bằng target. Trình bày cách làm và độ phức tạp."}, CorrectAnswer: BilingualText{En: "Use two pointers; move left/right based on sum; O(n) time, O(1) space.", Vi: "Dùng two pointers; dịch con trỏ theo sum; O(n) thời gian, O(1) bộ nhớ."}, Hint: BilingualText{En: "Sorted array suggests pointers.", Vi: "Mảng sort gợi ý two pointers."}},
		{Topic: "Algorithms", Level: "Junior", Role: "Any", Title: BilingualText{En: "Two Sum", Vi: "Two Sum"}, Content: BilingualText{En: "How would you solve Two Sum on an unsorted array?", Vi: "Bạn giải bài Two Sum với mảng chưa sắp xếp như thế nào?"}, CorrectAnswer: BilingualText{En: "Use a hash map; track complements; O(n) average time.", Vi: "Dùng hash map; tra complement; O(n) thời gian trung bình."}, Hint: BilingualText{En: "Store value -> index.", Vi: "Lưu value -> index."}},
	}

	patterns := []struct {
		level   string
		titleEn string
		titleVi string
		qEn     string
		qVi     string
		hEn     string
		hVi     string
	}{
		{"Junior", "Binary search pitfalls", "Pitfall binary search", "Explain binary search and common off-by-one pitfalls.", "Giải thích binary search và lỗi off-by-one thường gặp.", "Define invariants for l/r.", "Đặt invariant cho l/r."},
		{"Junior", "Merge intervals", "Gộp intervals", "Given a list of intervals, merge overlaps. What is the approach?", "Cho danh sách interval, gộp các interval chồng lấp. Cách làm?", "Sort by start then merge.", "Sort theo start rồi merge."},
		{"Mid", "Sliding window", "Sliding window", "Explain how to find the longest substring without repeating characters.", "Giải thích cách tìm substring dài nhất không lặp ký tự.", "Window + set/map.", "Window + set/map."},
		{"Mid", "Minimum window substring", "Minimum window substring", "How would you find the minimum window substring that contains all characters of another string?", "Tìm minimum window substring chứa đủ ký tự của chuỗi khác như thế nào?", "Sliding window with counts.", "Sliding window + đếm tần suất."},
		{"Mid", "Top K frequent", "Top K xuất hiện nhiều", "How do you find top K frequent elements efficiently?", "Làm sao tìm top K phần tử xuất hiện nhiều nhất hiệu quả?", "Heap or bucket counting.", "Heap hoặc bucket."},
		{"Mid", "Kth largest element", "Phần tử lớn thứ K", "Find the Kth largest element in an array. What is the best approach?", "Tìm phần tử lớn thứ K trong mảng. Cách tối ưu?", "Quickselect or heap.", "Quickselect hoặc heap."},
		{"Mid", "Kadane algorithm", "Thuật toán Kadane", "Explain maximum subarray sum and Kadane’s algorithm.", "Giải thích bài max subarray sum và thuật toán Kadane.", "DP: best ending here.", "DP: best ending here."},
		{"Senior", "LRU cache design", "Thiết kế LRU cache", "Design an LRU cache. What data structures do you use?", "Thiết kế LRU cache. Dùng cấu trúc dữ liệu gì?", "Hash map + doubly linked list.", "Hash map + doubly linked list."},
		{"Senior", "Union-Find", "Union-Find", "What is Union-Find and where is it used?", "Union-Find là gì và dùng ở đâu?", "Connected components, cycle detection.", "Thành phần liên thông, phát hiện chu trình."},
		{"Senior", "Dijkstra vs BFS", "Dijkstra vs BFS", "When do you use Dijkstra vs BFS?", "Khi nào dùng Dijkstra vs BFS?", "Weighted edges -> Dijkstra.", "Có trọng số -> Dijkstra."},
		{"Mid", "BFS vs DFS", "BFS vs DFS", "When do you use BFS vs DFS? Give an example.", "Khi nào dùng BFS vs DFS? Cho ví dụ.", "Unweighted shortest path -> BFS.", "Unweighted shortest path -> BFS."},
		{"Mid", "Backtracking permutations", "Backtracking hoán vị", "Generate all permutations of a list. Explain backtracking approach.", "Sinh tất cả hoán vị của một list. Trình bày backtracking.", "Try/undo; prune when possible.", "Thử/hoàn tác; prune khi có thể."},
		{"Junior", "Valid parentheses", "Ngoặc hợp lệ", "Given a string of brackets, determine if it is valid.", "Cho chuỗi ngoặc, kiểm tra hợp lệ.", "Use a stack.", "Dùng stack."},
		{"Junior", "Reverse linked list", "Đảo linked list", "Reverse a singly linked list. What is the approach?", "Đảo ngược singly linked list. Cách làm?", "Iterative pointers or recursion.", "Pointer lặp hoặc đệ quy."},
		{"Mid", "Detect cycle in list", "Phát hiện cycle", "Detect a cycle in a linked list.", "Phát hiện cycle trong linked list.", "Floyd's slow/fast pointers.", "Slow/fast pointer (Floyd)."},
		{"Junior", "Merge sorted lists", "Gộp list đã sort", "Merge two sorted linked lists.", "Gộp 2 linked list đã sort.", "Two pointers.", "Two pointers."},
		{"Mid", "Binary tree inorder", "Inorder cây nhị phân", "Traverse a binary tree inorder iteratively.", "Duyệt inorder cây nhị phân theo kiểu iterative.", "Use a stack.", "Dùng stack."},
		{"Junior", "Level order traversal", "Duyệt theo tầng", "Traverse a binary tree level order.", "Duyệt cây nhị phân theo tầng.", "Queue (BFS).", "Queue (BFS)."},
		{"Mid", "Lowest common ancestor", "LCA", "Find the lowest common ancestor in a binary tree.", "Tìm LCA trong cây nhị phân.", "Recursive postorder reasoning.", "Suy luận đệ quy postorder."},
		{"Senior", "Serialize binary tree", "Serialize cây nhị phân", "Serialize and deserialize a binary tree.", "Serialize và deserialize cây nhị phân.", "Use BFS/DFS with null markers.", "Dùng BFS/DFS kèm null marker."},
		{"Mid", "Topological sort", "Topological sort", "Given prerequisites between courses, determine if you can finish all courses.", "Bài course schedule: xác định có thể học hết không.", "Topo sort (Kahn) or DFS cycle detection.", "Topo sort hoặc DFS phát hiện cycle."},
		{"Senior", "Cycle in directed graph", "Cycle graph có hướng", "Detect a cycle in a directed graph.", "Phát hiện cycle trong graph có hướng.", "DFS colors or Kahn.", "DFS colors hoặc Kahn."},
		{"Senior", "Trie design", "Thiết kế Trie", "Design a Trie (prefix tree) and discuss complexity.", "Thiết kế Trie (prefix tree) và độ phức tạp.", "Nodes with children map; O(L).", "Node + map children; O(L)."},
		{"Mid", "Group anagrams", "Nhóm anagram", "Group anagrams together. What key do you use?", "Nhóm các từ anagram. Dùng key gì?", "Sort string or frequency signature.", "Sort chuỗi hoặc signature tần suất."},
		{"Senior", "Longest increasing subsequence", "Dãy tăng dài nhất", "Explain approaches for Longest Increasing Subsequence.", "Trình bày cách giải Longest Increasing Subsequence.", "DP O(n^2) or patience O(n log n).", "DP O(n^2) hoặc patience O(n log n)."},
		{"Senior", "Edit distance", "Edit distance", "Explain dynamic programming for edit distance.", "Giải thích DP cho edit distance.", "DP table; transitions insert/delete/replace.", "Bảng DP; insert/delete/replace."},
		{"Mid", "Coin change", "Coin change", "Given coin denominations, find minimum coins to make amount.", "Cho các đồng xu, tìm số xu ít nhất để tạo amount.", "DP with min transitions.", "DP tìm min."},
		{"Senior", "0/1 knapsack", "Balo 0/1", "Explain 0/1 knapsack DP and optimization.", "Giải thích DP balo 0/1 và tối ưu.", "DP with capacity dimension.", "DP theo capacity."},
		{"Senior", "Median of two sorted arrays", "Median 2 mảng sort", "Find the median of two sorted arrays in O(log(min(n,m))).", "Tìm median của 2 mảng đã sort với O(log(min(n,m))).", "Binary search partition.", "Binary search partition."},
		{"Senior", "KMP substring search", "Tìm chuỗi con KMP", "Explain KMP substring search and prefix function.", "Giải thích KMP và hàm prefix.", "Prefix table avoids backtracking.", "Bảng prefix tránh backtracking."},
		{"Mid", "Count bits", "Đếm bit", "Compute number of 1-bits for all numbers 0..n.", "Tính số bit 1 cho tất cả số 0..n.", "DP: ans[i]=ans[i>>1]+(i&1).", "DP: ans[i]=ans[i>>1]+(i&1)."},
		{"Junior", "Climbing stairs", "Leo cầu thang", "How many ways to climb n stairs with 1 or 2 steps?", "Có bao nhiêu cách leo n bậc với bước 1 hoặc 2?", "DP / Fibonacci.", "DP / Fibonacci."},
	}

	for i := 0; len(templates) < 35; i++ {
		p := patterns[i]
		templates = append(templates, QuestionTemplate{
			Topic:         "Algorithms",
			Level:         p.level,
			Role:          "Any",
			Title:         BilingualText{En: p.titleEn, Vi: p.titleVi},
			Content:       BilingualText{En: p.qEn, Vi: p.qVi},
			CorrectAnswer: BilingualText{En: "Explain approach, complexity, and edge cases.", Vi: "Nêu cách làm, độ phức tạp và edge cases."},
			Hint:          BilingualText{En: p.hEn, Vi: p.hVi},
		})
	}

	return templates[:35]
}
func BuildSystemDesignTemplates() []QuestionTemplate {
	concepts := []Concept{
		{Name: BilingualText{En: "URL shortener", Vi: "Rút gọn URL"}, KeyPoints: BilingualText{En: "API, key generation, storage, redirects, collisions", Vi: "API, sinh key, lưu trữ, redirect, collision"}, Hint: BilingualText{En: "Start with APIs and data model.", Vi: "Bắt đầu từ API và data model."}},
		{Name: BilingualText{En: "Rate limiting", Vi: "Rate limiting"}, KeyPoints: BilingualText{En: "token bucket, sliding window, distributed counters", Vi: "token bucket, sliding window, counter phân tán"}, Hint: BilingualText{En: "Decide failure mode (fail-open/close).", Vi: "Chọn fail-open/close."}},
		{Name: BilingualText{En: "Caching strategy", Vi: "Chiến lược cache"}, KeyPoints: BilingualText{En: "cache-aside, TTL, invalidation", Vi: "cache-aside, TTL, invalidation"}, Hint: BilingualText{En: "Explain invalidation clearly.", Vi: "Nêu invalidation rõ ràng."}},
		{Name: BilingualText{En: "Event-driven processing", Vi: "Xử lý event-driven"}, KeyPoints: BilingualText{En: "queues, retries, idempotency, ordering, DLQ", Vi: "queue, retry, idempotency, ordering, DLQ"}, Hint: BilingualText{En: "Mention DLQ and retry policy.", Vi: "Nêu DLQ và policy retry."}},
		{Name: BilingualText{En: "Authentication & authorization", Vi: "Xác thực & phân quyền"}, KeyPoints: BilingualText{En: "JWT/OAuth, sessions, RBAC, token rotation", Vi: "JWT/OAuth, session, RBAC, rotate token"}, Hint: BilingualText{En: "Separate authn vs authz.", Vi: "Tách authn vs authz."}},
		{Name: BilingualText{En: "Search service", Vi: "Dịch vụ search"}, KeyPoints: BilingualText{En: "indexing, relevance, latency, consistency", Vi: "indexing, relevance, latency, nhất quán"}, Hint: BilingualText{En: "Define query patterns first.", Vi: "Xác định pattern query trước."}},
		{Name: BilingualText{En: "File upload service", Vi: "Dịch vụ upload file"}, KeyPoints: BilingualText{En: "chunk upload, scanning, CDN, metadata", Vi: "upload chunk, quét, CDN, metadata"}, Hint: BilingualText{En: "Think retries and idempotency.", Vi: "Nghĩ về retry và idempotency."}},
		{Name: BilingualText{En: "Notifications", Vi: "Thông báo"}, KeyPoints: BilingualText{En: "fan-out, preferences, delivery guarantees", Vi: "fan-out, preference, đảm bảo giao"}, Hint: BilingualText{En: "Discuss push vs pull.", Vi: "Nêu push vs pull."}},
		{Name: BilingualText{En: "API gateway", Vi: "API gateway"}, KeyPoints: BilingualText{En: "routing, auth, rate limit, observability", Vi: "routing, auth, rate limit, observability"}, Hint: BilingualText{En: "List cross-cutting concerns.", Vi: "Nêu các cross-cutting concern."}},
		{Name: BilingualText{En: "Observability", Vi: "Observability"}, KeyPoints: BilingualText{En: "metrics, logs, tracing, SLOs", Vi: "metrics, logs, tracing, SLO"}, Hint: BilingualText{En: "Use golden signals.", Vi: "Dùng golden signals."}},
	}

	formats := []Format{
		{Level: "Mid", Title: BilingualText{En: "Design: {Concept}", Vi: "Thiết kế: {Concept}"}, Content: BilingualText{En: "Design {Concept}. What components do you need?", Vi: "Thiết kế {Concept}. Bạn cần những thành phần nào?"}, CorrectAnswer: BilingualText{En: "Cover: {KeyPoints}. Explain trade-offs.", Vi: "Trình bày: {KeyPoints}. Nêu trade-off."}, Hint: BilingualText{En: "Start from requirements.", Vi: "Bắt đầu từ requirements."}},
		{Level: "Mid", Title: BilingualText{En: "{Concept} data model", Vi: "Data model {Concept}"}, Content: BilingualText{En: "What data model would you use for {Concept}?", Vi: "Bạn dùng data model nào cho {Concept}?"}, CorrectAnswer: BilingualText{En: "Define entities and keys based on access patterns: {KeyPoints}.", Vi: "Xác định entity và key theo access pattern: {KeyPoints}."}, Hint: BilingualText{En: "Work backward from queries.", Vi: "Đi từ query ngược lại."}},
		{Level: "Senior", Title: BilingualText{En: "{Concept} trade-offs", Vi: "Trade-off {Concept}"}, Content: BilingualText{En: "What are the key trade-offs in {Concept} at scale?", Vi: "Trade-off chính của {Concept} khi scale là gì?"}, CorrectAnswer: BilingualText{En: "Discuss latency/cost/reliability using: {KeyPoints}.", Vi: "Phân tích latency/cost/reliability theo: {KeyPoints}."}, Hint: BilingualText{En: "Think failure modes.", Vi: "Nghĩ về failure mode."}},
		{Level: "Senior", Title: BilingualText{En: "{Concept} reliability", Vi: "Độ tin cậy {Concept}"}, Content: BilingualText{En: "How do you make {Concept} reliable in production?", Vi: "Làm sao để {Concept} tin cậy trên production?"}, CorrectAnswer: BilingualText{En: "Include resilience patterns and operations: {KeyPoints}.", Vi: "Nêu resilience pattern và vận hành: {KeyPoints}."}, Hint: BilingualText{En: "Define SLOs and runbooks.", Vi: "Xác định SLO và runbook."}},
		{Level: "Senior", Title: BilingualText{En: "{Concept} rollout", Vi: "Rollout {Concept}"}, Content: BilingualText{En: "How would you roll out changes for {Concept} safely?", Vi: "Bạn rollout thay đổi cho {Concept} an toàn thế nào?"}, CorrectAnswer: BilingualText{En: "Use canary/flags/monitoring and rollback plans: {KeyPoints}.", Vi: "Dùng canary/flag/monitoring và plan rollback: {KeyPoints}."}, Hint: BilingualText{En: "Choose metrics for success.", Vi: "Chọn metrics thành công."}},
	}

	return BuildFromConcepts("System Design", "Any", concepts, formats)
}

func BuildDatabaseTemplates() []QuestionTemplate {
	concepts := []Concept{
		{Name: BilingualText{En: "Normalization vs denormalization", Vi: "Normalize vs denormalize"}, KeyPoints: BilingualText{En: "OLTP vs OLAP, anomalies, access patterns", Vi: "OLTP vs OLAP, anomaly, access pattern"}, Hint: BilingualText{En: "Tie it to workload.", Vi: "Gắn với workload."}},
		{Name: BilingualText{En: "Indexes and query plans", Vi: "Index và query plan"}, KeyPoints: BilingualText{En: "selectivity, composite indexes, EXPLAIN", Vi: "selectivity, index ghép, EXPLAIN"}, Hint: BilingualText{En: "Use EXPLAIN/ANALYZE.", Vi: "Dùng EXPLAIN/ANALYZE."}},
		{Name: BilingualText{En: "Transactions and isolation levels", Vi: "Transaction và isolation level"}, KeyPoints: BilingualText{En: "ACID, anomalies, locking, serializable", Vi: "ACID, anomaly, lock, serializable"}, Hint: BilingualText{En: "Mention phantom reads.", Vi: "Nêu phantom read."}},
		{Name: BilingualText{En: "Deadlocks", Vi: "Deadlock"}, KeyPoints: BilingualText{En: "lock ordering, detection, retries", Vi: "thứ tự lock, phát hiện, retry"}, Hint: BilingualText{En: "Explain prevention.", Vi: "Nêu cách phòng ngừa."}},
		{Name: BilingualText{En: "Schema constraints", Vi: "Constraint schema"}, KeyPoints: BilingualText{En: "PK/FK/unique/check, integrity", Vi: "PK/FK/unique/check, integrity"}, Hint: BilingualText{En: "Prefer DB constraints when possible.", Vi: "Ưu tiên constraint DB khi có thể."}},
		{Name: BilingualText{En: "Idempotency", Vi: "Idempotency"}, KeyPoints: BilingualText{En: "idempotency keys, retries, duplicates", Vi: "idempotency key, retry, trùng lặp"}, Hint: BilingualText{En: "Use payment/order examples.", Vi: "Dùng ví dụ payment/order."}},
		{Name: BilingualText{En: "Data migrations", Vi: "Migration dữ liệu"}, KeyPoints: BilingualText{En: "backward compatible, dual writes, rollback", Vi: "tương thích ngược, dual write, rollback"}, Hint: BilingualText{En: "Plan zero downtime.", Vi: "Lên kế hoạch zero downtime."}},
		{Name: BilingualText{En: "Sharding vs replication", Vi: "Sharding vs replication"}, KeyPoints: BilingualText{En: "scale reads/writes, consistency, operations", Vi: "scale đọc/ghi, nhất quán, vận hành"}, Hint: BilingualText{En: "Explain trade-offs.", Vi: "Nêu trade-off."}},
		{Name: BilingualText{En: "N+1 queries", Vi: "N+1 queries"}, KeyPoints: BilingualText{En: "joins, batching, caching", Vi: "join, batching, cache"}, Hint: BilingualText{En: "Measure and batch.", Vi: "Đo và batch."}},
		{Name: BilingualText{En: "Soft delete", Vi: "Soft delete"}, KeyPoints: BilingualText{En: "retention, uniqueness, restore, indexes", Vi: "retention, unique, restore, index"}, Hint: BilingualText{En: "Consider uniqueness constraints.", Vi: "Chú ý unique constraint."}},
	}

	formats := []Format{
		{Level: "Junior", Title: BilingualText{En: "{Concept}", Vi: "{Concept}"}, Content: BilingualText{En: "Explain {Concept}.", Vi: "Giải thích {Concept}."}, CorrectAnswer: BilingualText{En: "Cover: {KeyPoints}.", Vi: "Trình bày: {KeyPoints}."}, Hint: BilingualText{En: "Use a small example.", Vi: "Dùng ví dụ nhỏ."}},
		{Level: "Mid", Title: BilingualText{En: "{Concept} in production", Vi: "{Concept} trên production"}, Content: BilingualText{En: "How do you apply {Concept} in a production system?", Vi: "Bạn áp dụng {Concept} trong hệ thống production như thế nào?"}, CorrectAnswer: BilingualText{En: "Describe steps and trade-offs: {KeyPoints}.", Vi: "Nêu bước làm và trade-off: {KeyPoints}."}, Hint: BilingualText{En: "Mention monitoring and rollback.", Vi: "Nêu monitoring và rollback."}},
		{Level: "Senior", Title: BilingualText{En: "{Concept} pitfalls", Vi: "Pitfall {Concept}"}, Content: BilingualText{En: "What pitfalls do you watch for with {Concept}?", Vi: "Bạn theo dõi pitfall gì với {Concept}?"}, CorrectAnswer: BilingualText{En: "List pitfalls and mitigations: {KeyPoints}.", Vi: "Liệt kê pitfall và giảm thiểu: {KeyPoints}."}, Hint: BilingualText{En: "Think scale and concurrency.", Vi: "Nghĩ scale và concurrency."}},
		{Level: "Mid", Title: BilingualText{En: "{Concept} troubleshooting", Vi: "Troubleshoot {Concept}"}, Content: BilingualText{En: "How do you troubleshoot an issue related to {Concept}?", Vi: "Bạn troubleshoot vấn đề liên quan {Concept} thế nào?"}, CorrectAnswer: BilingualText{En: "Use evidence-driven steps: {KeyPoints}.", Vi: "Đi theo hướng dựa dữ liệu: {KeyPoints}."}, Hint: BilingualText{En: "Start from symptoms.", Vi: "Bắt đầu từ triệu chứng."}},
		{Level: "Senior", Title: BilingualText{En: "Design review: {Concept}", Vi: "Review: {Concept}"}, Content: BilingualText{En: "Review a design that uses {Concept}. What would you change?", Vi: "Review một thiết kế dùng {Concept}. Bạn sẽ thay đổi gì?"}, CorrectAnswer: BilingualText{En: "Propose improvements based on: {KeyPoints}.", Vi: "Đề xuất cải tiến dựa trên: {KeyPoints}."}, Hint: BilingualText{En: "Challenge assumptions.", Vi: "Thử thách giả định."}},
	}

	return BuildFromConcepts("Database", "Any", concepts, formats)
}

func BuildTestingTemplates() []QuestionTemplate {
	concepts := []Concept{
		{Name: BilingualText{En: "Test pyramid", Vi: "Test pyramid"}, KeyPoints: BilingualText{En: "unit/integration/e2e trade-offs", Vi: "trade-off unit/integration/e2e"}, Hint: BilingualText{En: "Prefer fast and stable tests.", Vi: "Ưu tiên test nhanh và ổn định."}},
		{Name: BilingualText{En: "Mocking strategy", Vi: "Chiến lược mock"}, KeyPoints: BilingualText{En: "mock boundaries, avoid over-mocking", Vi: "mock boundary, tránh over-mock"}, Hint: BilingualText{En: "Mock IO, not business logic.", Vi: "Mock IO, không mock logic lõi."}},
		{Name: BilingualText{En: "Flaky tests", Vi: "Flaky test"}, KeyPoints: BilingualText{En: "determinism, isolation, timing", Vi: "deterministic, cô lập, timing"}, Hint: BilingualText{En: "Remove shared state.", Vi: "Loại bỏ shared state."}},
		{Name: BilingualText{En: "Contract testing", Vi: "Contract test"}, KeyPoints: BilingualText{En: "API contracts, consumer-driven tests", Vi: "contract API, consumer-driven"}, Hint: BilingualText{En: "Keep contracts versioned.", Vi: "Version contract."}},
		{Name: BilingualText{En: "Test data management", Vi: "Quản lý test data"}, KeyPoints: BilingualText{En: "fixtures, factories, cleanup, seeding", Vi: "fixture, factory, cleanup, seed"}, Hint: BilingualText{En: "Keep test data predictable.", Vi: "Test data cần predict được."}},
	}

	formats := []Format{
		{Level: "Junior", Title: BilingualText{En: "{Concept}", Vi: "{Concept}"}, Content: BilingualText{En: "Explain {Concept} and why it matters.", Vi: "Giải thích {Concept} và vì sao quan trọng."}, CorrectAnswer: BilingualText{En: "Cover: {KeyPoints}.", Vi: "Trình bày: {KeyPoints}."}, Hint: BilingualText{En: "Use a simple example.", Vi: "Dùng ví dụ đơn giản."}},
		{Level: "Mid", Title: BilingualText{En: "{Concept} in CI", Vi: "{Concept} trong CI"}, Content: BilingualText{En: "How do you apply {Concept} in CI pipelines?", Vi: "Bạn áp dụng {Concept} trong CI pipeline thế nào?"}, CorrectAnswer: BilingualText{En: "Describe process and tooling: {KeyPoints}.", Vi: "Nêu quy trình và công cụ: {KeyPoints}."}, Hint: BilingualText{En: "Focus on trust and speed.", Vi: "Tập trung độ tin cậy và tốc độ."}},
		{Level: "Senior", Title: BilingualText{En: "Scaling tests: {Concept}", Vi: "Scale test: {Concept}"}, Content: BilingualText{En: "What breaks when tests scale and how does {Concept} help?", Vi: "Khi test tăng quy mô thì điều gì vỡ và {Concept} giúp gì?"}, CorrectAnswer: BilingualText{En: "Discuss maintenance, speed, flakiness: {KeyPoints}.", Vi: "Nêu bảo trì, tốc độ, flaky: {KeyPoints}."}, Hint: BilingualText{En: "Think determinism.", Vi: "Nghĩ deterministic."}},
		{Level: "Mid", Title: BilingualText{En: "{Concept} trade-offs", Vi: "Trade-off {Concept}"}, Content: BilingualText{En: "What trade-offs exist with {Concept}?", Vi: "Trade-off của {Concept} là gì?"}, CorrectAnswer: BilingualText{En: "Explain trade-offs and mitigations: {KeyPoints}.", Vi: "Giải thích trade-off và giảm thiểu: {KeyPoints}."}, Hint: BilingualText{En: "Be specific about costs.", Vi: "Nêu rõ chi phí."}},
		{Level: "Senior", Title: BilingualText{En: "Design: {Concept}", Vi: "Thiết kế: {Concept}"}, Content: BilingualText{En: "Design a testing approach using {Concept} for a service.", Vi: "Thiết kế hướng test dùng {Concept} cho một service."}, CorrectAnswer: BilingualText{En: "Outline layers and strategy: {KeyPoints}.", Vi: "Phác thảo các lớp và chiến lược: {KeyPoints}."}, Hint: BilingualText{En: "Start from risks.", Vi: "Bắt đầu từ rủi ro."}},
	}

	return BuildFromConcepts("Testing", "Any", concepts, formats)
}

func BuildDevOpsTemplates() []QuestionTemplate {
	concepts := []Concept{
		{Name: BilingualText{En: "CI/CD pipeline", Vi: "CI/CD pipeline"}, KeyPoints: BilingualText{En: "build, test, scan, deploy, verify", Vi: "build, test, scan, deploy, verify"}, Hint: BilingualText{En: "Mention gating and rollback.", Vi: "Nêu gating và rollback."}},
		{Name: BilingualText{En: "Monitoring and alerting", Vi: "Monitoring và alert"}, KeyPoints: BilingualText{En: "SLOs, burn rate, noise reduction", Vi: "SLO, burn rate, giảm nhiễu"}, Hint: BilingualText{En: "Tie to user impact.", Vi: "Gắn với user impact."}},
		{Name: BilingualText{En: "Incident response", Vi: "Ứng phó sự cố"}, KeyPoints: BilingualText{En: "triage, mitigation, postmortem", Vi: "triage, giảm thiểu, postmortem"}, Hint: BilingualText{En: "Show a clear timeline.", Vi: "Nêu timeline rõ."}},
		{Name: BilingualText{En: "Secrets management", Vi: "Quản lý secret"}, KeyPoints: BilingualText{En: "rotation, least privilege, vault", Vi: "rotate, least privilege, vault"}, Hint: BilingualText{En: "Never hardcode secrets.", Vi: "Không hardcode secret."}},
		{Name: BilingualText{En: "Deployment strategies", Vi: "Chiến lược deploy"}, KeyPoints: BilingualText{En: "blue/green, canary, feature flags", Vi: "blue/green, canary, feature flag"}, Hint: BilingualText{En: "Define rollback conditions.", Vi: "Nêu điều kiện rollback."}},
		{Name: BilingualText{En: "Networking basics", Vi: "Networking cơ bản"}, KeyPoints: BilingualText{En: "DNS, TLS, load balancers", Vi: "DNS, TLS, load balancer"}, Hint: BilingualText{En: "Explain request path.", Vi: "Nêu đường đi request."}},
		{Name: BilingualText{En: "Capacity planning", Vi: "Capacity planning"}, KeyPoints: BilingualText{En: "load testing, headroom, scaling rules", Vi: "load test, headroom, rule scale"}, Hint: BilingualText{En: "Use metrics not guesses.", Vi: "Dùng metrics, không đoán."}},
		{Name: BilingualText{En: "Config management", Vi: "Quản lý cấu hình"}, KeyPoints: BilingualText{En: "env vars, config validation, rollout", Vi: "env, validate config, rollout"}, Hint: BilingualText{En: "Separate config from code.", Vi: "Tách config khỏi code."}},
		{Name: BilingualText{En: "Security basics", Vi: "Security cơ bản"}, KeyPoints: BilingualText{En: "OWASP, patching, least privilege", Vi: "OWASP, vá lỗi, least privilege"}, Hint: BilingualText{En: "Threats and controls.", Vi: "Threat và control."}},
		{Name: BilingualText{En: "Logging & tracing", Vi: "Logging & tracing"}, KeyPoints: BilingualText{En: "structured logs, trace IDs, sampling", Vi: "log cấu trúc, trace ID, sampling"}, Hint: BilingualText{En: "Correlate logs to requests.", Vi: "Gắn log với request."}},
	}

	formats := []Format{
		{Level: "Junior", Title: BilingualText{En: "{Concept}", Vi: "{Concept}"}, Content: BilingualText{En: "Explain {Concept}.", Vi: "Giải thích {Concept}."}, CorrectAnswer: BilingualText{En: "Cover: {KeyPoints}.", Vi: "Trình bày: {KeyPoints}."}, Hint: BilingualText{En: "Use a real example.", Vi: "Dùng ví dụ thực tế."}},
		{Level: "Mid", Title: BilingualText{En: "{Concept} practice", Vi: "Thực hành {Concept}"}, Content: BilingualText{En: "How do you implement {Concept} in a real system?", Vi: "Bạn triển khai {Concept} trong hệ thống thực tế như thế nào?"}, CorrectAnswer: BilingualText{En: "Describe steps and tools: {KeyPoints}.", Vi: "Nêu bước làm và công cụ: {KeyPoints}."}, Hint: BilingualText{En: "Mention automation.", Vi: "Nêu automation."}},
		{Level: "Senior", Title: BilingualText{En: "{Concept} trade-offs", Vi: "Trade-off {Concept}"}, Content: BilingualText{En: "What trade-offs exist in {Concept}?", Vi: "Trade-off của {Concept} là gì?"}, CorrectAnswer: BilingualText{En: "Discuss trade-offs and mitigations: {KeyPoints}.", Vi: "Nêu trade-off và giảm thiểu: {KeyPoints}."}, Hint: BilingualText{En: "Consider cost and reliability.", Vi: "Chú ý cost và reliability."}},
		{Level: "Mid", Title: BilingualText{En: "{Concept} troubleshooting", Vi: "Troubleshoot {Concept}"}, Content: BilingualText{En: "How do you troubleshoot issues related to {Concept}?", Vi: "Bạn troubleshoot vấn đề liên quan {Concept} thế nào?"}, CorrectAnswer: BilingualText{En: "Use evidence and rollback plans: {KeyPoints}.", Vi: "Dựa trên dữ liệu và plan rollback: {KeyPoints}."}, Hint: BilingualText{En: "Start from symptoms.", Vi: "Bắt đầu từ triệu chứng."}},
		{Level: "Senior", Title: BilingualText{En: "Design: {Concept}", Vi: "Thiết kế: {Concept}"}, Content: BilingualText{En: "Design an approach for {Concept} for a critical service.", Vi: "Thiết kế cách tiếp cận {Concept} cho service quan trọng."}, CorrectAnswer: BilingualText{En: "Outline approach: {KeyPoints}.", Vi: "Phác thảo cách làm: {KeyPoints}."}, Hint: BilingualText{En: "Define SLOs and failure modes.", Vi: "Xác định SLO và failure mode."}},
	}

	return BuildFromConcepts("DevOps", "DevOps", concepts, formats)
}

func BuildDataLayerTemplates() []QuestionTemplate {
	dbs := []struct {
		name      BilingualText
		keyPoints BilingualText
		hint      BilingualText
	}{
		{name: BilingualText{En: "PostgreSQL", Vi: "PostgreSQL"}, keyPoints: BilingualText{En: "indexes, EXPLAIN, transactions, locks", Vi: "index, EXPLAIN, transaction, lock"}, hint: BilingualText{En: "Use EXPLAIN/ANALYZE.", Vi: "Dùng EXPLAIN/ANALYZE."}},
		{name: BilingualText{En: "Redis", Vi: "Redis"}, keyPoints: BilingualText{En: "cache patterns, TTL, eviction, persistence", Vi: "pattern cache, TTL, eviction, persistence"}, hint: BilingualText{En: "Discuss eviction policy.", Vi: "Nêu eviction policy."}},
		{name: BilingualText{En: "MongoDB", Vi: "MongoDB"}, keyPoints: BilingualText{En: "document model, indexes, aggregation", Vi: "document model, index, aggregation"}, hint: BilingualText{En: "Mention indexes and schema design.", Vi: "Nêu index và schema."}},
		{name: BilingualText{En: "SQL databases", Vi: "Database SQL"}, keyPoints: BilingualText{En: "ACID, joins, constraints", Vi: "ACID, join, constraint"}, hint: BilingualText{En: "Relate to OLTP.", Vi: "Liên hệ OLTP."}},
	}

	areas := []Concept{
		{Name: BilingualText{En: "indexing", Vi: "indexing"}, KeyPoints: BilingualText{En: "{KeyPoints}", Vi: "{KeyPoints}"}, Hint: BilingualText{En: "", Vi: ""}},
		{Name: BilingualText{En: "transactions", Vi: "transactions"}, KeyPoints: BilingualText{En: "{KeyPoints}", Vi: "{KeyPoints}"}, Hint: BilingualText{En: "", Vi: ""}},
		{Name: BilingualText{En: "caching", Vi: "caching"}, KeyPoints: BilingualText{En: "{KeyPoints}", Vi: "{KeyPoints}"}, Hint: BilingualText{En: "", Vi: ""}},
		{Name: BilingualText{En: "schema design", Vi: "thiết kế schema"}, KeyPoints: BilingualText{En: "{KeyPoints}", Vi: "{KeyPoints}"}, Hint: BilingualText{En: "", Vi: ""}},
		{Name: BilingualText{En: "consistency", Vi: "nhất quán"}, KeyPoints: BilingualText{En: "{KeyPoints}", Vi: "{KeyPoints}"}, Hint: BilingualText{En: "", Vi: ""}},
		{Name: BilingualText{En: "replication", Vi: "replication"}, KeyPoints: BilingualText{En: "{KeyPoints}", Vi: "{KeyPoints}"}, Hint: BilingualText{En: "", Vi: ""}},
		{Name: BilingualText{En: "backups", Vi: "backup"}, KeyPoints: BilingualText{En: "{KeyPoints}", Vi: "{KeyPoints}"}, Hint: BilingualText{En: "", Vi: ""}},
		{Name: BilingualText{En: "migration strategy", Vi: "chiến lược migration"}, KeyPoints: BilingualText{En: "{KeyPoints}", Vi: "{KeyPoints}"}, Hint: BilingualText{En: "", Vi: ""}},
		{Name: BilingualText{En: "performance tuning", Vi: "tối ưu hiệu năng"}, KeyPoints: BilingualText{En: "{KeyPoints}", Vi: "{KeyPoints}"}, Hint: BilingualText{En: "", Vi: ""}},
		{Name: BilingualText{En: "data modeling", Vi: "mô hình dữ liệu"}, KeyPoints: BilingualText{En: "{KeyPoints}", Vi: "{KeyPoints}"}, Hint: BilingualText{En: "", Vi: ""}},
	}

	templates := make([]QuestionTemplate, 0, 40)
	for _, db := range dbs {
		for i := 0; i < 10; i++ {
			a := areas[i]
			dataEn := map[string]string{"Db": db.name.En, "Concept": a.Name.En, "KeyPoints": db.keyPoints.En}
			dataVi := map[string]string{"Db": db.name.Vi, "Concept": a.Name.Vi, "KeyPoints": db.keyPoints.Vi}
			templates = append(templates, QuestionTemplate{
				Topic: "Data Layer",
				Level: "Junior",
				Role:  "Any",
				Title: RenderBilingual(BilingualText{En: "{Db} {Concept}", Vi: "{Db} {Concept}"}, dataEn, dataVi),
				Content: RenderBilingual(BilingualText{
					En: "In {Db}, explain {Concept} and common mistakes.",
					Vi: "Trong {Db}, giải thích {Concept} và lỗi thường gặp.",
				}, dataEn, dataVi),
				CorrectAnswer: RenderBilingual(BilingualText{
					En: "Cover key points: {KeyPoints}.",
					Vi: "Trình bày các ý chính: {KeyPoints}.",
				}, dataEn, dataVi),
				Hint: RenderBilingual(BilingualText{En: db.hint.En, Vi: db.hint.Vi}, dataEn, dataVi),
			})
		}
	}

	return templates
}
func BuildGolangTemplates() []QuestionTemplate {
	concepts := []Concept{
		{Name: BilingualText{En: "goroutines and scheduling", Vi: "goroutine và scheduling"}, KeyPoints: BilingualText{En: "M:N scheduler, stack growth, context switches", Vi: "scheduler M:N, stack growth, context switch"}, Hint: BilingualText{En: "Mention runtime scheduler and costs.", Vi: "Nêu scheduler runtime và chi phí."}},
		{Name: BilingualText{En: "channels and buffering", Vi: "channel và buffer"}, KeyPoints: BilingualText{En: "blocking semantics, buffer size, deadlocks", Vi: "blocking, kích thước buffer, deadlock"}, Hint: BilingualText{En: "Explain how buffering changes behavior.", Vi: "Nêu buffer làm đổi hành vi."}},
		{Name: BilingualText{En: "context cancellation", Vi: "context cancellation"}, KeyPoints: BilingualText{En: "deadlines, propagation, request scope", Vi: "deadline, lan truyền, scope request"}, Hint: BilingualText{En: "Use HTTP request example.", Vi: "Dùng ví dụ HTTP request."}},
		{Name: BilingualText{En: "error wrapping", Vi: "wrap error"}, KeyPoints: BilingualText{En: "%w wrapping, typed errors, preserving cause", Vi: "wrap %w, typed error, giữ cause"}, Hint: BilingualText{En: "Mention errors.Is/errors.As.", Vi: "Nêu errors.Is/errors.As."}},
		{Name: BilingualText{En: "panic and recover", Vi: "panic và recover"}, KeyPoints: BilingualText{En: "when to use, middleware recovery, avoid hiding bugs", Vi: "khi dùng, middleware recovery, tránh che bug"}, Hint: BilingualText{En: "Distinguish programmer vs runtime errors.", Vi: "Tách lỗi lập trình và lỗi runtime."}},
		{Name: BilingualText{En: "interfaces design", Vi: "thiết kế interface"}, KeyPoints: BilingualText{En: "small interfaces, consumer-defined, testability", Vi: "interface nhỏ, consumer-defined, dễ test"}, Hint: BilingualText{En: "Mention io.Reader style.", Vi: "Nêu io.Reader."}},
		{Name: BilingualText{En: "sync primitives", Vi: "sync primitive"}, KeyPoints: BilingualText{En: "mutex/rwmutex/once/waitgroup patterns", Vi: "pattern mutex/rwmutex/once/waitgroup"}, Hint: BilingualText{En: "Talk about contention.", Vi: "Nêu contention."}},
		{Name: BilingualText{En: "race conditions", Vi: "race condition"}, KeyPoints: BilingualText{En: "race detector, shared state, locks vs channels", Vi: "race detector, shared state, lock vs channel"}, Hint: BilingualText{En: "Mention go test -race.", Vi: "Nêu go test -race."}},
		{Name: BilingualText{En: "graceful shutdown", Vi: "graceful shutdown"}, KeyPoints: BilingualText{En: "signals, draining, timeouts, context", Vi: "signal, drain, timeout, context"}, Hint: BilingualText{En: "Explain shutdown sequence.", Vi: "Nêu thứ tự shutdown."}},
		{Name: BilingualText{En: "benchmarking and profiling", Vi: "benchmark và profiling"}, KeyPoints: BilingualText{En: "benchmarks, allocations, pprof hotspots", Vi: "benchmark, allocation, hotspot pprof"}, Hint: BilingualText{En: "Measure before optimizing.", Vi: "Đo trước khi tối ưu."}},
	}

	formats := []Format{
		{Level: "Junior", Title: BilingualText{En: "Go: {Concept}", Vi: "Go: {Concept}"}, Content: BilingualText{En: "Explain {Concept} in Go.", Vi: "Giải thích {Concept} trong Go."}, CorrectAnswer: BilingualText{En: "Include: {KeyPoints}.", Vi: "Bao gồm: {KeyPoints}."}, Hint: BilingualText{En: "Give a short example.", Vi: "Cho ví dụ ngắn."}},
		{Level: "Mid", Title: BilingualText{En: "Applying {Concept}", Vi: "Áp dụng {Concept}"}, Content: BilingualText{En: "How do you apply {Concept} in a production Go service?", Vi: "Bạn áp dụng {Concept} trong Go service production như thế nào?"}, CorrectAnswer: BilingualText{En: "Describe approach and pitfalls: {KeyPoints}.", Vi: "Nêu cách làm và pitfall: {KeyPoints}."}, Hint: BilingualText{En: "Mention edge cases.", Vi: "Nêu edge case."}},
		{Level: "Senior", Title: BilingualText{En: "{Concept} pitfalls", Vi: "Pitfall {Concept}"}, Content: BilingualText{En: "What are common pitfalls with {Concept} and how do you mitigate them?", Vi: "Pitfall thường gặp của {Concept} và cách giảm thiểu?"}, CorrectAnswer: BilingualText{En: "List pitfalls and mitigations: {KeyPoints}.", Vi: "Liệt kê pitfall và giảm thiểu: {KeyPoints}."}, Hint: BilingualText{En: "Think failure modes.", Vi: "Nghĩ về failure mode."}},
		{Level: "Senior", Title: BilingualText{En: "Design review: {Concept}", Vi: "Review: {Concept}"}, Content: BilingualText{En: "Review a Go design that relies on {Concept}. What would you change?", Vi: "Review thiết kế Go dùng {Concept}. Bạn sẽ đổi gì?"}, CorrectAnswer: BilingualText{En: "Propose improvements based on: {KeyPoints}.", Vi: "Đề xuất cải tiến dựa trên: {KeyPoints}."}, Hint: BilingualText{En: "Challenge assumptions.", Vi: "Thử thách giả định."}},
	}

	return BuildFromConcepts("Golang", "BackEnd", concepts, formats)
}

func BuildNodeJSTemplates() []QuestionTemplate {
	concepts := []Concept{
		{Name: BilingualText{En: "event loop and microtasks", Vi: "event loop và microtask"}, KeyPoints: BilingualText{En: "call stack, macrotasks, microtasks, IO", Vi: "call stack, macrotask, microtask, IO"}, Hint: BilingualText{En: "Explain order of logs for async code.", Vi: "Nêu thứ tự log cho async."}},
		{Name: BilingualText{En: "streams and backpressure", Vi: "stream và backpressure"}, KeyPoints: BilingualText{En: "stream types, buffering, flow control", Vi: "loại stream, buffer, điều khiển luồng"}, Hint: BilingualText{En: "Explain why backpressure matters.", Vi: "Nêu vì sao backpressure quan trọng."}},
		{Name: BilingualText{En: "memory leaks", Vi: "memory leak"}, KeyPoints: BilingualText{En: "heap snapshots, retained objects, listeners", Vi: "heap snapshot, object bị giữ, listener"}, Hint: BilingualText{En: "Mention profiling tools.", Vi: "Nêu tool profiling."}},
		{Name: BilingualText{En: "worker threads vs cluster", Vi: "worker threads vs cluster"}, KeyPoints: BilingualText{En: "CPU-bound tasks, scaling, isolation", Vi: "tác vụ CPU, scale, cô lập"}, Hint: BilingualText{En: "Explain when to use each.", Vi: "Nêu khi nào dùng."}},
		{Name: BilingualText{En: "error handling in APIs", Vi: "xử lý lỗi trong API"}, KeyPoints: BilingualText{En: "middleware, safe errors, correlation IDs", Vi: "middleware, lỗi an toàn, correlation ID"}, Hint: BilingualText{En: "Show a consistent error shape.", Vi: "Thống nhất schema lỗi."}},
		{Name: BilingualText{En: "security basics", Vi: "security cơ bản"}, KeyPoints: BilingualText{En: "input validation, SSRF, XSS, rate limiting", Vi: "validate input, SSRF, XSS, rate limit"}, Hint: BilingualText{En: "Think OWASP.", Vi: "Nghĩ OWASP."}},
		{Name: BilingualText{En: "async patterns", Vi: "pattern bất đồng bộ"}, KeyPoints: BilingualText{En: "promises, async/await, error propagation", Vi: "promise, async/await, lan truyền lỗi"}, Hint: BilingualText{En: "Mention Promise.all vs allSettled.", Vi: "Nêu all vs allSettled."}},
		{Name: BilingualText{En: "logging and tracing", Vi: "logging và tracing"}, KeyPoints: BilingualText{En: "structured logs, trace IDs, sampling", Vi: "log cấu trúc, trace ID, sampling"}, Hint: BilingualText{En: "Correlate logs to requests.", Vi: "Gắn log với request."}},
		{Name: BilingualText{En: "API performance", Vi: "hiệu năng API"}, KeyPoints: BilingualText{En: "p99 latency, caching, batching", Vi: "p99 latency, cache, batching"}, Hint: BilingualText{En: "Start from profiling.", Vi: "Bắt đầu từ profiling."}},
		{Name: BilingualText{En: "database access patterns", Vi: "pattern truy cập DB"}, KeyPoints: BilingualText{En: "pooling, N+1, transactions", Vi: "pooling, N+1, transaction"}, Hint: BilingualText{En: "Discuss connection pooling.", Vi: "Nêu connection pool."}},
	}

	formats := []Format{
		{Level: "Junior", Title: BilingualText{En: "Node.js: {Concept}", Vi: "Node.js: {Concept}"}, Content: BilingualText{En: "Explain {Concept} in Node.js.", Vi: "Giải thích {Concept} trong Node.js."}, CorrectAnswer: BilingualText{En: "Cover: {KeyPoints}.", Vi: "Trình bày: {KeyPoints}."}, Hint: BilingualText{En: "Use a concrete example.", Vi: "Dùng ví dụ cụ thể."}},
		{Level: "Mid", Title: BilingualText{En: "Production: {Concept}", Vi: "Production: {Concept}"}, Content: BilingualText{En: "How do you handle {Concept} in a production Node.js API?", Vi: "Bạn xử lý {Concept} trong Node.js API production như thế nào?"}, CorrectAnswer: BilingualText{En: "Describe approach and trade-offs: {KeyPoints}.", Vi: "Nêu cách làm và trade-off: {KeyPoints}."}, Hint: BilingualText{En: "Mention tooling and metrics.", Vi: "Nêu tool và metrics."}},
		{Level: "Senior", Title: BilingualText{En: "Pitfalls: {Concept}", Vi: "Pitfall: {Concept}"}, Content: BilingualText{En: "What pitfalls exist with {Concept} and how do you mitigate them?", Vi: "Pitfall của {Concept} và cách giảm thiểu?"}, CorrectAnswer: BilingualText{En: "List risks and mitigations: {KeyPoints}.", Vi: "Liệt kê rủi ro và giảm thiểu: {KeyPoints}."}, Hint: BilingualText{En: "Think scale and failures.", Vi: "Nghĩ scale và failure."}},
	}

	return BuildFromConcepts("NodeJS", "BackEnd", concepts, formats)
}

func BuildJavaScriptTemplates() []QuestionTemplate {
	concepts := []Concept{
		{Name: BilingualText{En: "scope and hoisting", Vi: "scope và hoisting"}, KeyPoints: BilingualText{En: "var vs let/const, block scope, TDZ", Vi: "var vs let/const, block scope, TDZ"}, Hint: BilingualText{En: "Explain TDZ briefly.", Vi: "Nêu TDZ ngắn gọn."}},
		{Name: BilingualText{En: "closures", Vi: "closure"}, KeyPoints: BilingualText{En: "lexical environment, callbacks, encapsulation", Vi: "lexical env, callback, đóng gói"}, Hint: BilingualText{En: "Explain captured variables.", Vi: "Giải thích biến được capture."}},
		{Name: BilingualText{En: "prototypes", Vi: "prototype"}, KeyPoints: BilingualText{En: "prototype chain, inheritance, lookup", Vi: "prototype chain, kế thừa, lookup"}, Hint: BilingualText{En: "Describe lookup along the chain.", Vi: "Nêu lookup theo chain."}},
		{Name: BilingualText{En: "promises and async/await", Vi: "promise và async/await"}, KeyPoints: BilingualText{En: "microtasks, error propagation, concurrency", Vi: "microtask, lan truyền lỗi, concurrency"}, Hint: BilingualText{En: "Compare all vs allSettled.", Vi: "So sánh all và allSettled."}},
		{Name: BilingualText{En: "event loop", Vi: "event loop"}, KeyPoints: BilingualText{En: "tasks vs microtasks, timers, IO", Vi: "task vs microtask, timer, IO"}, Hint: BilingualText{En: "Think ordering questions.", Vi: "Nghĩ về thứ tự."}},
		{Name: BilingualText{En: "immutability", Vi: "immutability"}, KeyPoints: BilingualText{En: "const semantics, shallow vs deep copy", Vi: "const, shallow vs deep copy"}, Hint: BilingualText{En: "Explain reference vs value.", Vi: "Nêu tham chiếu vs giá trị."}},
		{Name: BilingualText{En: "TypeScript narrowing", Vi: "TypeScript narrowing"}, KeyPoints: BilingualText{En: "union types, type guards, discriminated unions", Vi: "union type, type guard, discriminated union"}, Hint: BilingualText{En: "Use a union example.", Vi: "Dùng ví dụ union."}},
		{Name: BilingualText{En: "runtime debugging", Vi: "debug runtime"}, KeyPoints: BilingualText{En: "stack traces, source maps, logging", Vi: "stack trace, source map, logging"}, Hint: BilingualText{En: "Explain a debugging approach.", Vi: "Nêu cách debug."}},
		{Name: BilingualText{En: "web security", Vi: "bảo mật web"}, KeyPoints: BilingualText{En: "XSS, CSRF, CORS, cookies", Vi: "XSS, CSRF, CORS, cookie"}, Hint: BilingualText{En: "Threats and mitigations.", Vi: "Threat và mitigation."}},
		{Name: BilingualText{En: "performance basics", Vi: "hiệu năng cơ bản"}, KeyPoints: BilingualText{En: "bundle size, caching, rendering cost", Vi: "bundle size, cache, chi phí render"}, Hint: BilingualText{En: "Mention measurement tools.", Vi: "Nêu tool đo."}},
	}

	formats := []Format{
		{Level: "Junior", Title: BilingualText{En: "JS: {Concept}", Vi: "JS: {Concept}"}, Content: BilingualText{En: "Explain {Concept}.", Vi: "Giải thích {Concept}."}, CorrectAnswer: BilingualText{En: "Cover: {KeyPoints}.", Vi: "Trình bày: {KeyPoints}."}, Hint: BilingualText{En: "Give a quick example.", Vi: "Cho ví dụ nhanh."}},
		{Level: "Mid", Title: BilingualText{En: "{Concept} in practice", Vi: "{Concept} trong thực tế"}, Content: BilingualText{En: "How does {Concept} show up in real projects?", Vi: "{Concept} xuất hiện trong dự án thực tế như thế nào?"}, CorrectAnswer: BilingualText{En: "Explain scenarios and pitfalls: {KeyPoints}.", Vi: "Nêu tình huống và pitfall: {KeyPoints}."}, Hint: BilingualText{En: "Tie to production bugs.", Vi: "Gắn với bug production."}},
		{Level: "Senior", Title: BilingualText{En: "Pitfalls: {Concept}", Vi: "Pitfall: {Concept}"}, Content: BilingualText{En: "What pitfalls exist with {Concept} and how do you avoid them?", Vi: "Pitfall của {Concept} và cách tránh?"}, CorrectAnswer: BilingualText{En: "List pitfalls and mitigations: {KeyPoints}.", Vi: "Liệt kê pitfall và giảm thiểu: {KeyPoints}."}, Hint: BilingualText{En: "Think correctness and edge cases.", Vi: "Nghĩ đúng đắn và edge case."}},
	}

	return BuildFromConcepts("JavaScript", "Any", concepts, formats)
}
func BuildReactTemplates() []QuestionTemplate {
	concepts := []Concept{
		{Name: BilingualText{En: "rendering and re-renders", Vi: "render và re-render"}, KeyPoints: BilingualText{En: "state/props, memoization, referential equality", Vi: "state/props, memoization, referential equality"}, Hint: BilingualText{En: "Explain why props identity matters.", Vi: "Nêu vì sao identity của props quan trọng."}},
		{Name: BilingualText{En: "useEffect pitfalls", Vi: "pitfall useEffect"}, KeyPoints: BilingualText{En: "deps, stale closures, cleanup", Vi: "deps, stale closure, cleanup"}, Hint: BilingualText{En: "Focus on dependency arrays.", Vi: "Tập trung vào dependency array."}},
		{Name: BilingualText{En: "state management", Vi: "quản lý state"}, KeyPoints: BilingualText{En: "local vs global, reducers, server state", Vi: "local vs global, reducer, server state"}, Hint: BilingualText{En: "Separate UI state and server state.", Vi: "Tách UI state và server state."}},
		{Name: BilingualText{En: "forms and validation", Vi: "form và validation"}, KeyPoints: BilingualText{En: "controlled/uncontrolled, performance, UX", Vi: "controlled/uncontrolled, hiệu năng, UX"}, Hint: BilingualText{En: "Mention accessibility.", Vi: "Nêu accessibility."}},
		{Name: BilingualText{En: "SSR and hydration", Vi: "SSR và hydration"}, KeyPoints: BilingualText{En: "SEO, performance, hydration mismatches", Vi: "SEO, hiệu năng, lệch hydration"}, Hint: BilingualText{En: "Explain mismatch causes.", Vi: "Nêu nguyên nhân mismatch."}},
	}

	formats := []Format{
		{Level: "Junior", Title: BilingualText{En: "React: {Concept}", Vi: "React: {Concept}"}, Content: BilingualText{En: "Explain {Concept} in React.", Vi: "Giải thích {Concept} trong React."}, CorrectAnswer: BilingualText{En: "Cover: {KeyPoints}.", Vi: "Trình bày: {KeyPoints}."}, Hint: BilingualText{En: "Use a component example.", Vi: "Dùng ví dụ component."}},
		{Level: "Mid", Title: BilingualText{En: "{Concept} optimization", Vi: "Tối ưu {Concept}"}, Content: BilingualText{En: "How do you optimize issues related to {Concept}?", Vi: "Bạn tối ưu vấn đề liên quan {Concept} như thế nào?"}, CorrectAnswer: BilingualText{En: "Explain steps and trade-offs: {KeyPoints}.", Vi: "Nêu bước làm và trade-off: {KeyPoints}."}, Hint: BilingualText{En: "Measure before optimizing.", Vi: "Đo trước khi tối ưu."}},
		{Level: "Mid", Title: BilingualText{En: "{Concept} pitfalls", Vi: "Pitfall {Concept}"}, Content: BilingualText{En: "What pitfalls exist with {Concept}?", Vi: "Pitfall của {Concept} là gì?"}, CorrectAnswer: BilingualText{En: "List pitfalls and fixes: {KeyPoints}.", Vi: "Liệt kê pitfall và cách fix: {KeyPoints}."}, Hint: BilingualText{En: "Think stale state/effects.", Vi: "Nghĩ về stale state/effect."}},
		{Level: "Senior", Title: BilingualText{En: "Design: {Concept}", Vi: "Thiết kế: {Concept}"}, Content: BilingualText{En: "Design a feature considering {Concept}. What would you watch for?", Vi: "Thiết kế một feature với {Concept}. Bạn chú ý gì?"}, CorrectAnswer: BilingualText{En: "Cover design and trade-offs: {KeyPoints}.", Vi: "Nêu thiết kế và trade-off: {KeyPoints}."}, Hint: BilingualText{En: "Talk about DX and performance.", Vi: "Nêu DX và hiệu năng."}},
		{Level: "Senior", Title: BilingualText{En: "Debugging: {Concept}", Vi: "Debug: {Concept}"}, Content: BilingualText{En: "How do you debug problems related to {Concept}?", Vi: "Bạn debug vấn đề liên quan {Concept} thế nào?"}, CorrectAnswer: BilingualText{En: "Explain debugging approach: {KeyPoints}.", Vi: "Nêu cách debug: {KeyPoints}."}, Hint: BilingualText{En: "Use React DevTools.", Vi: "Dùng React DevTools."}},
	}

	return BuildFromConcepts("React", "FrontEnd", concepts, formats)
}

func BuildVueTemplates() []QuestionTemplate {
	concepts := []Concept{
		{Name: BilingualText{En: "reactivity system", Vi: "hệ thống reactivity"}, KeyPoints: BilingualText{En: "dependency tracking, refs, computed", Vi: "track dependency, ref, computed"}, Hint: BilingualText{En: "Explain how updates propagate.", Vi: "Nêu update lan truyền."}},
		{Name: BilingualText{En: "Composition API", Vi: "Composition API"}, KeyPoints: BilingualText{En: "setup, composables, reuse logic", Vi: "setup, composable, tái sử dụng logic"}, Hint: BilingualText{En: "Mention separation of concerns.", Vi: "Nêu separation of concerns."}},
		{Name: BilingualText{En: "state management", Vi: "quản lý state"}, KeyPoints: BilingualText{En: "Pinia/Vuex patterns, modules", Vi: "pattern Pinia/Vuex, module"}, Hint: BilingualText{En: "Separate server state.", Vi: "Tách server state."}},
		{Name: BilingualText{En: "performance", Vi: "hiệu năng"}, KeyPoints: BilingualText{En: "rendering, watchers, computed caching", Vi: "render, watcher, cache computed"}, Hint: BilingualText{En: "Measure and optimize.", Vi: "Đo và tối ưu."}},
		{Name: BilingualText{En: "routing", Vi: "routing"}, KeyPoints: BilingualText{En: "route guards, lazy loading, auth", Vi: "route guard, lazy load, auth"}, Hint: BilingualText{En: "Discuss auth guards.", Vi: "Nêu auth guard."}},
	}

	formats := []Format{
		{Level: "Junior", Title: BilingualText{En: "Vue: {Concept}", Vi: "Vue: {Concept}"}, Content: BilingualText{En: "Explain {Concept} in Vue.", Vi: "Giải thích {Concept} trong Vue."}, CorrectAnswer: BilingualText{En: "Cover: {KeyPoints}.", Vi: "Trình bày: {KeyPoints}."}, Hint: BilingualText{En: "Use a small example.", Vi: "Dùng ví dụ nhỏ."}},
		{Level: "Mid", Title: BilingualText{En: "{Concept} pitfalls", Vi: "Pitfall {Concept}"}, Content: BilingualText{En: "What pitfalls exist with {Concept} in Vue apps?", Vi: "Pitfall của {Concept} trong app Vue là gì?"}, CorrectAnswer: BilingualText{En: "List pitfalls and fixes: {KeyPoints}.", Vi: "Liệt kê pitfall và cách fix: {KeyPoints}."}, Hint: BilingualText{En: "Think watchers and reactivity.", Vi: "Nghĩ về watcher và reactivity."}},
		{Level: "Senior", Title: BilingualText{En: "Design: {Concept}", Vi: "Thiết kế: {Concept}"}, Content: BilingualText{En: "Design a Vue feature considering {Concept}.", Vi: "Thiết kế feature Vue với {Concept}."}, CorrectAnswer: BilingualText{En: "Describe design and trade-offs: {KeyPoints}.", Vi: "Nêu thiết kế và trade-off: {KeyPoints}."}, Hint: BilingualText{En: "Consider maintainability.", Vi: "Chú ý bảo trì."}},
	}

	return BuildFromConcepts("Vue", "FrontEnd", concepts, formats)
}

func BuildPythonTemplates() []QuestionTemplate {
	concepts := []Concept{
		{Name: BilingualText{En: "virtual environments", Vi: "virtualenv"}, KeyPoints: BilingualText{En: "dependency isolation, lock files, reproducibility", Vi: "cô lập dependency, lock file, tái hiện"}, Hint: BilingualText{En: "Explain reproducible builds.", Vi: "Nêu build tái hiện được."}},
		{Name: BilingualText{En: "asyncio", Vi: "asyncio"}, KeyPoints: BilingualText{En: "event loop, awaitables, concurrency", Vi: "event loop, awaitable, concurrency"}, Hint: BilingualText{En: "Contrast with threads.", Vi: "So sánh với thread."}},
		{Name: BilingualText{En: "GIL", Vi: "GIL"}, KeyPoints: BilingualText{En: "CPU vs IO bound, multiprocessing", Vi: "CPU vs IO bound, multiprocessing"}, Hint: BilingualText{En: "Mention when GIL matters.", Vi: "Nêu khi nào GIL ảnh hưởng."}},
		{Name: BilingualText{En: "typing and validation", Vi: "typing và validation"}, KeyPoints: BilingualText{En: "type hints, mypy, pydantic", Vi: "type hint, mypy, pydantic"}, Hint: BilingualText{En: "Separate static typing vs runtime validation.", Vi: "Tách static typing và runtime validation."}},
		{Name: BilingualText{En: "profiling performance", Vi: "profiling hiệu năng"}, KeyPoints: BilingualText{En: "cProfile, flamegraphs, hotspots", Vi: "cProfile, flamegraph, hotspot"}, Hint: BilingualText{En: "Measure before optimizing.", Vi: "Đo trước khi tối ưu."}},
	}

	formats := []Format{
		{Level: "Junior", Title: BilingualText{En: "Python: {Concept}", Vi: "Python: {Concept}"}, Content: BilingualText{En: "Explain {Concept} in Python services.", Vi: "Giải thích {Concept} trong Python service."}, CorrectAnswer: BilingualText{En: "Cover: {KeyPoints}.", Vi: "Trình bày: {KeyPoints}."}, Hint: BilingualText{En: "Use a small example.", Vi: "Dùng ví dụ nhỏ."}},
		{Level: "Mid", Title: BilingualText{En: "{Concept} in production", Vi: "{Concept} trên production"}, Content: BilingualText{En: "How do you apply {Concept} in production Python systems?", Vi: "Bạn áp dụng {Concept} trong hệ thống Python production thế nào?"}, CorrectAnswer: BilingualText{En: "Describe approach and trade-offs: {KeyPoints}.", Vi: "Nêu cách làm và trade-off: {KeyPoints}."}, Hint: BilingualText{En: "Mention testing and observability.", Vi: "Nêu testing và observability."}},
		{Level: "Senior", Title: BilingualText{En: "{Concept} pitfalls", Vi: "Pitfall {Concept}"}, Content: BilingualText{En: "What pitfalls exist with {Concept} and how do you mitigate them?", Vi: "Pitfall của {Concept} và cách giảm thiểu?"}, CorrectAnswer: BilingualText{En: "List pitfalls and mitigations: {KeyPoints}.", Vi: "Liệt kê pitfall và giảm thiểu: {KeyPoints}."}, Hint: BilingualText{En: "Think reliability.", Vi: "Nghĩ reliability."}},
		{Level: "Senior", Title: BilingualText{En: "Design review: {Concept}", Vi: "Review: {Concept}"}, Content: BilingualText{En: "Review a design using {Concept}. What would you improve?", Vi: "Review thiết kế dùng {Concept}. Bạn cải tiến gì?"}, CorrectAnswer: BilingualText{En: "Suggest improvements using: {KeyPoints}.", Vi: "Gợi ý cải tiến theo: {KeyPoints}."}, Hint: BilingualText{En: "Challenge assumptions.", Vi: "Thử thách giả định."}},
	}

	return BuildFromConcepts("Python", "BackEnd", concepts, formats)
}
func BuildJavaTemplates() []QuestionTemplate {
	concepts := []Concept{
		{Name: BilingualText{En: "JVM memory & GC", Vi: "Bộ nhớ JVM & GC"}, KeyPoints: BilingualText{En: "heap/stack, generations, GC trade-offs", Vi: "heap/stack, generation, trade-off GC"}, Hint: BilingualText{En: "Mention GC tuning and profiling.", Vi: "Nêu tuning GC và profiling."}},
		{Name: BilingualText{En: "concurrency", Vi: "concurrency"}, KeyPoints: BilingualText{En: "synchronization, thread pools, concurrent collections", Vi: "sync, thread pool, collection concurrent"}, Hint: BilingualText{En: "Talk about contention.", Vi: "Nêu contention."}},
		{Name: BilingualText{En: "exceptions", Vi: "exception"}, KeyPoints: BilingualText{En: "checked vs unchecked, error handling", Vi: "checked vs unchecked, xử lý lỗi"}, Hint: BilingualText{En: "Explain when to use each.", Vi: "Nêu khi nào dùng."}},
		{Name: BilingualText{En: "streams API", Vi: "Streams API"}, KeyPoints: BilingualText{En: "map/filter/reduce, laziness, performance", Vi: "map/filter/reduce, lazy, hiệu năng"}, Hint: BilingualText{En: "Mention readability vs performance.", Vi: "Nêu readability vs performance."}},
		{Name: BilingualText{En: "Spring patterns", Vi: "pattern Spring"}, KeyPoints: BilingualText{En: "DI, bean lifecycle, configuration", Vi: "DI, vòng đời bean, cấu hình"}, Hint: BilingualText{En: "Tie to testability.", Vi: "Gắn với testability."}},
	}

	formats := []Format{
		{Level: "Junior", Title: BilingualText{En: "Java: {Concept}", Vi: "Java: {Concept}"}, Content: BilingualText{En: "Explain {Concept}.", Vi: "Giải thích {Concept}."}, CorrectAnswer: BilingualText{En: "Cover: {KeyPoints}.", Vi: "Trình bày: {KeyPoints}."}, Hint: BilingualText{En: "Use a small example.", Vi: "Dùng ví dụ nhỏ."}},
		{Level: "Mid", Title: BilingualText{En: "{Concept} in production", Vi: "{Concept} trên production"}, Content: BilingualText{En: "How do you apply {Concept} in production Java systems?", Vi: "Bạn áp dụng {Concept} trong hệ thống Java production thế nào?"}, CorrectAnswer: BilingualText{En: "Describe approach and trade-offs: {KeyPoints}.", Vi: "Nêu cách làm và trade-off: {KeyPoints}."}, Hint: BilingualText{En: "Mention monitoring and profiling.", Vi: "Nêu monitoring và profiling."}},
		{Level: "Senior", Title: BilingualText{En: "{Concept} pitfalls", Vi: "Pitfall {Concept}"}, Content: BilingualText{En: "What pitfalls exist with {Concept} and how do you mitigate them?", Vi: "Pitfall của {Concept} và cách giảm thiểu?"}, CorrectAnswer: BilingualText{En: "List pitfalls and mitigations: {KeyPoints}.", Vi: "Liệt kê pitfall và giảm thiểu: {KeyPoints}."}, Hint: BilingualText{En: "Think scale and correctness.", Vi: "Nghĩ scale và đúng đắn."}},
		{Level: "Senior", Title: BilingualText{En: "Design review: {Concept}", Vi: "Review: {Concept}"}, Content: BilingualText{En: "Review a design using {Concept}. What would you improve?", Vi: "Review thiết kế dùng {Concept}. Bạn cải tiến gì?"}, CorrectAnswer: BilingualText{En: "Suggest improvements using: {KeyPoints}.", Vi: "Gợi ý cải tiến theo: {KeyPoints}."}, Hint: BilingualText{En: "Challenge assumptions.", Vi: "Thử thách giả định."}},
	}

	return BuildFromConcepts("Java", "BackEnd", concepts, formats)
}

func BuildCSharpTemplates() []QuestionTemplate {
	concepts := []Concept{
		{Name: BilingualText{En: "async/await", Vi: "async/await"}, KeyPoints: BilingualText{En: "Tasks, synchronization context, deadlocks", Vi: "Task, sync context, deadlock"}, Hint: BilingualText{En: "Mention ConfigureAwait and deadlocks.", Vi: "Nêu ConfigureAwait và deadlock."}},
		{Name: BilingualText{En: "dependency injection", Vi: "dependency injection"}, KeyPoints: BilingualText{En: "lifetimes, scopes, testability", Vi: "lifetime, scope, testability"}, Hint: BilingualText{En: "Explain singleton pitfalls.", Vi: "Nêu pitfall singleton."}},
		{Name: BilingualText{En: "LINQ", Vi: "LINQ"}, KeyPoints: BilingualText{En: "deferred execution, performance, allocations", Vi: "deferred execution, hiệu năng, allocation"}, Hint: BilingualText{En: "Mention query translation in ORMs.", Vi: "Nêu dịch query trong ORM."}},
		{Name: BilingualText{En: "GC and memory", Vi: "GC và memory"}, KeyPoints: BilingualText{En: "generations, LOH, allocations", Vi: "generation, LOH, allocation"}, Hint: BilingualText{En: "Profile allocations.", Vi: "Profile allocation."}},
		{Name: BilingualText{En: "ASP.NET middleware", Vi: "middleware ASP.NET"}, KeyPoints: BilingualText{En: "pipeline, ordering, cross-cutting concerns", Vi: "pipeline, thứ tự, cross-cutting"}, Hint: BilingualText{En: "Explain request/response flow.", Vi: "Nêu flow request/response."}},
	}

	formats := []Format{
		{Level: "Junior", Title: BilingualText{En: "C#: {Concept}", Vi: "C#: {Concept}"}, Content: BilingualText{En: "Explain {Concept}.", Vi: "Giải thích {Concept}."}, CorrectAnswer: BilingualText{En: "Cover: {KeyPoints}.", Vi: "Trình bày: {KeyPoints}."}, Hint: BilingualText{En: "Use a short example.", Vi: "Dùng ví dụ ngắn."}},
		{Level: "Mid", Title: BilingualText{En: "{Concept} in production", Vi: "{Concept} trên production"}, Content: BilingualText{En: "How do you apply {Concept} in production .NET services?", Vi: "Bạn áp dụng {Concept} trong .NET service production thế nào?"}, CorrectAnswer: BilingualText{En: "Describe approach and trade-offs: {KeyPoints}.", Vi: "Nêu cách làm và trade-off: {KeyPoints}."}, Hint: BilingualText{En: "Mention tooling.", Vi: "Nêu tooling."}},
		{Level: "Senior", Title: BilingualText{En: "{Concept} pitfalls", Vi: "Pitfall {Concept}"}, Content: BilingualText{En: "What pitfalls exist with {Concept} and how do you mitigate them?", Vi: "Pitfall của {Concept} và cách giảm thiểu?"}, CorrectAnswer: BilingualText{En: "List pitfalls and mitigations: {KeyPoints}.", Vi: "Liệt kê pitfall và giảm thiểu: {KeyPoints}."}, Hint: BilingualText{En: "Think scale and correctness.", Vi: "Nghĩ scale và đúng đắn."}},
	}

	return BuildFromConcepts("C#", "BackEnd", concepts, formats)
}

func BuildAngularTemplates() []QuestionTemplate {
	concepts := []Concept{
		{Name: BilingualText{En: "change detection", Vi: "change detection"}, KeyPoints: BilingualText{En: "default vs OnPush, performance", Vi: "default vs OnPush, hiệu năng"}, Hint: BilingualText{En: "Explain when OnPush helps.", Vi: "Nêu khi nào OnPush giúp."}},
		{Name: BilingualText{En: "RxJS", Vi: "RxJS"}, KeyPoints: BilingualText{En: "observables, operators, subscriptions", Vi: "observable, operator, subscription"}, Hint: BilingualText{En: "Mention unsubscribe patterns.", Vi: "Nêu pattern unsubscribe."}},
		{Name: BilingualText{En: "dependency injection", Vi: "dependency injection"}, KeyPoints: BilingualText{En: "providers, scopes, testing", Vi: "provider, scope, test"}, Hint: BilingualText{En: "Tie to testability.", Vi: "Gắn với testability."}},
		{Name: BilingualText{En: "routing and guards", Vi: "routing và guard"}, KeyPoints: BilingualText{En: "lazy loading, route guards, auth", Vi: "lazy load, guard, auth"}, Hint: BilingualText{En: "Explain auth guards.", Vi: "Nêu guard auth."}},
		{Name: BilingualText{En: "forms", Vi: "forms"}, KeyPoints: BilingualText{En: "reactive forms, validation, UX", Vi: "reactive form, validation, UX"}, Hint: BilingualText{En: "Mention validation strategy.", Vi: "Nêu chiến lược validation."}},
	}

	formats := []Format{
		{Level: "Junior", Title: BilingualText{En: "Angular: {Concept}", Vi: "Angular: {Concept}"}, Content: BilingualText{En: "Explain {Concept} in Angular.", Vi: "Giải thích {Concept} trong Angular."}, CorrectAnswer: BilingualText{En: "Cover: {KeyPoints}.", Vi: "Trình bày: {KeyPoints}."}, Hint: BilingualText{En: "Use a small example.", Vi: "Dùng ví dụ nhỏ."}},
		{Level: "Mid", Title: BilingualText{En: "{Concept} pitfalls", Vi: "Pitfall {Concept}"}, Content: BilingualText{En: "What pitfalls exist with {Concept} in Angular apps?", Vi: "Pitfall của {Concept} trong app Angular là gì?"}, CorrectAnswer: BilingualText{En: "List pitfalls and fixes: {KeyPoints}.", Vi: "Liệt kê pitfall và cách fix: {KeyPoints}."}, Hint: BilingualText{En: "Think performance and leaks.", Vi: "Nghĩ hiệu năng và leak."}},
	}

	return BuildFromConcepts("Angular", "FrontEnd", concepts, formats)
}

func BuildDjangoTemplates() []QuestionTemplate {
	concepts := []Concept{
		{Name: BilingualText{En: "ORM performance", Vi: "hiệu năng ORM"}, KeyPoints: BilingualText{En: "select_related/prefetch_related, N+1", Vi: "select_related/prefetch_related, N+1"}, Hint: BilingualText{En: "Explain query optimization.", Vi: "Nêu tối ưu query."}},
		{Name: BilingualText{En: "middleware", Vi: "middleware"}, KeyPoints: BilingualText{En: "request/response lifecycle, ordering", Vi: "vòng đời request/response, thứ tự"}, Hint: BilingualText{En: "Mention auth/logging use cases.", Vi: "Nêu auth/logging."}},
		{Name: BilingualText{En: "auth and permissions", Vi: "auth và permission"}, KeyPoints: BilingualText{En: "sessions, tokens, permission classes", Vi: "session, token, permission"}, Hint: BilingualText{En: "Separate authn and authz.", Vi: "Tách authn và authz."}},
		{Name: BilingualText{En: "migrations", Vi: "migration"}, KeyPoints: BilingualText{En: "schema changes, zero downtime, rollbacks", Vi: "đổi schema, zero downtime, rollback"}, Hint: BilingualText{En: "Plan backward compatibility.", Vi: "Lên kế hoạch tương thích ngược."}},
		{Name: BilingualText{En: "caching", Vi: "caching"}, KeyPoints: BilingualText{En: "cache backends, invalidation, TTL", Vi: "backend cache, invalidation, TTL"}, Hint: BilingualText{En: "Talk about cache keys.", Vi: "Nêu cache key."}},
	}

	formats := []Format{
		{Level: "Junior", Title: BilingualText{En: "Django: {Concept}", Vi: "Django: {Concept}"}, Content: BilingualText{En: "Explain {Concept} in Django.", Vi: "Giải thích {Concept} trong Django."}, CorrectAnswer: BilingualText{En: "Cover: {KeyPoints}.", Vi: "Trình bày: {KeyPoints}."}, Hint: BilingualText{En: "Use a production example.", Vi: "Dùng ví dụ production."}},
		{Level: "Mid", Title: BilingualText{En: "{Concept} pitfalls", Vi: "Pitfall {Concept}"}, Content: BilingualText{En: "What pitfalls exist with {Concept} in Django services?", Vi: "Pitfall của {Concept} trong Django service là gì?"}, CorrectAnswer: BilingualText{En: "List pitfalls and mitigations: {KeyPoints}.", Vi: "Liệt kê pitfall và giảm thiểu: {KeyPoints}."}, Hint: BilingualText{En: "Think performance and security.", Vi: "Nghĩ hiệu năng và bảo mật."}},
	}

	return BuildFromConcepts("Django", "BackEnd", concepts, formats)
}

func BuildSpringBootTemplates() []QuestionTemplate {
	concepts := []Concept{
		{Name: BilingualText{En: "dependency injection", Vi: "dependency injection"}, KeyPoints: BilingualText{En: "beans, scopes, configuration", Vi: "bean, scope, cấu hình"}, Hint: BilingualText{En: "Tie to testability.", Vi: "Gắn với testability."}},
		{Name: BilingualText{En: "transactions", Vi: "transaction"}, KeyPoints: BilingualText{En: "@Transactional, isolation, propagation", Vi: "@Transactional, isolation, propagation"}, Hint: BilingualText{En: "Explain propagation behaviors.", Vi: "Nêu propagation."}},
		{Name: BilingualText{En: "REST API design", Vi: "thiết kế REST API"}, KeyPoints: BilingualText{En: "validation, errors, pagination", Vi: "validation, lỗi, pagination"}, Hint: BilingualText{En: "Keep errors consistent.", Vi: "Thống nhất schema lỗi."}},
		{Name: BilingualText{En: "observability", Vi: "observability"}, KeyPoints: BilingualText{En: "metrics, logs, tracing", Vi: "metrics, log, tracing"}, Hint: BilingualText{En: "Mention actuator and tracing.", Vi: "Nêu actuator và tracing."}},
		{Name: BilingualText{En: "performance", Vi: "hiệu năng"}, KeyPoints: BilingualText{En: "thread pools, DB access, caching", Vi: "thread pool, DB, cache"}, Hint: BilingualText{En: "Measure p99 latency.", Vi: "Đo p99 latency."}},
	}

	formats := []Format{
		{Level: "Junior", Title: BilingualText{En: "Spring Boot: {Concept}", Vi: "Spring Boot: {Concept}"}, Content: BilingualText{En: "Explain {Concept} in Spring Boot.", Vi: "Giải thích {Concept} trong Spring Boot."}, CorrectAnswer: BilingualText{En: "Cover: {KeyPoints}.", Vi: "Trình bày: {KeyPoints}."}, Hint: BilingualText{En: "Use a service example.", Vi: "Dùng ví dụ service."}},
		{Level: "Mid", Title: BilingualText{En: "{Concept} pitfalls", Vi: "Pitfall {Concept}"}, Content: BilingualText{En: "What pitfalls exist with {Concept} in Spring Boot services?", Vi: "Pitfall của {Concept} trong Spring Boot service là gì?"}, CorrectAnswer: BilingualText{En: "List pitfalls and mitigations: {KeyPoints}.", Vi: "Liệt kê pitfall và giảm thiểu: {KeyPoints}."}, Hint: BilingualText{En: "Think production failures.", Vi: "Nghĩ lỗi production."}},
	}

	return BuildFromConcepts("Spring Boot", "BackEnd", concepts, formats)
}

func BuildDockerTemplates() []QuestionTemplate {
	concepts := []Concept{
		{Name: BilingualText{En: "Dockerfile best practices", Vi: "Best practice Dockerfile"}, KeyPoints: BilingualText{En: "layers, caching, multi-stage builds", Vi: "layer, cache, multi-stage"}, Hint: BilingualText{En: "Mention multi-stage builds.", Vi: "Nêu multi-stage."}},
		{Name: BilingualText{En: "image size optimization", Vi: "tối ưu kích thước image"}, KeyPoints: BilingualText{En: "minimal base, .dockerignore, remove deps", Vi: "base tối giản, .dockerignore, bỏ deps"}, Hint: BilingualText{En: "Pin versions.", Vi: "Pin version."}},
		{Name: BilingualText{En: "networking", Vi: "networking"}, KeyPoints: BilingualText{En: "bridge, host, port mapping, DNS", Vi: "bridge, host, map port, DNS"}, Hint: BilingualText{En: "Explain port mapping.", Vi: "Nêu port mapping."}},
		{Name: BilingualText{En: "volumes", Vi: "volume"}, KeyPoints: BilingualText{En: "bind mounts vs volumes, persistence", Vi: "bind mount vs volume, persistence"}, Hint: BilingualText{En: "Discuss data persistence.", Vi: "Nêu lưu dữ liệu."}},
		{Name: BilingualText{En: "container security", Vi: "bảo mật container"}, KeyPoints: BilingualText{En: "least privilege, non-root, scanning", Vi: "least privilege, non-root, scan"}, Hint: BilingualText{En: "Mention image scanning.", Vi: "Nêu scan image."}},
	}

	formats := []Format{
		{Level: "Junior", Title: BilingualText{En: "Docker: {Concept}", Vi: "Docker: {Concept}"}, Content: BilingualText{En: "Explain {Concept}.", Vi: "Giải thích {Concept}."}, CorrectAnswer: BilingualText{En: "Cover: {KeyPoints}.", Vi: "Trình bày: {KeyPoints}."}, Hint: BilingualText{En: "Use a practical example.", Vi: "Dùng ví dụ thực tế."}},
		{Level: "Mid", Title: BilingualText{En: "{Concept} in production", Vi: "{Concept} trên production"}, Content: BilingualText{En: "How do you apply {Concept} for production containers?", Vi: "Bạn áp dụng {Concept} cho container production thế nào?"}, CorrectAnswer: BilingualText{En: "Describe steps and trade-offs: {KeyPoints}.", Vi: "Nêu bước làm và trade-off: {KeyPoints}."}, Hint: BilingualText{En: "Think CI/CD integration.", Vi: "Nghĩ CI/CD."}},
		{Level: "Senior", Title: BilingualText{En: "{Concept} pitfalls", Vi: "Pitfall {Concept}"}, Content: BilingualText{En: "What pitfalls exist with {Concept} and how do you avoid them?", Vi: "Pitfall của {Concept} và cách tránh?"}, CorrectAnswer: BilingualText{En: "List pitfalls and mitigations: {KeyPoints}.", Vi: "Liệt kê pitfall và giảm thiểu: {KeyPoints}."}, Hint: BilingualText{En: "Think security and operations.", Vi: "Nghĩ bảo mật và vận hành."}},
	}

	return BuildFromConcepts("Docker", "DevOps", concepts, formats)
}

func BuildKubernetesTemplates() []QuestionTemplate {
	concepts := []Concept{
		{Name: BilingualText{En: "Deployments and rollouts", Vi: "Deployment và rollout"}, KeyPoints: BilingualText{En: "replicas, rolling updates, rollback", Vi: "replica, rolling update, rollback"}, Hint: BilingualText{En: "Explain rollout strategy.", Vi: "Nêu chiến lược rollout."}},
		{Name: BilingualText{En: "Services and networking", Vi: "Service và networking"}, KeyPoints: BilingualText{En: "ClusterIP/NodePort/LoadBalancer, DNS", Vi: "ClusterIP/NodePort/LoadBalancer, DNS"}, Hint: BilingualText{En: "Explain traffic routing.", Vi: "Nêu định tuyến traffic."}},
		{Name: BilingualText{En: "Ingress", Vi: "Ingress"}, KeyPoints: BilingualText{En: "HTTP routing, TLS termination, annotations", Vi: "routing HTTP, terminate TLS, annotation"}, Hint: BilingualText{En: "Mention controllers.", Vi: "Nêu ingress controller."}},
		{Name: BilingualText{En: "ConfigMaps and Secrets", Vi: "ConfigMap và Secret"}, KeyPoints: BilingualText{En: "config injection, rotation, RBAC", Vi: "inject config, rotate, RBAC"}, Hint: BilingualText{En: "Avoid committing secrets.", Vi: "Không commit secret."}},
		{Name: BilingualText{En: "resource requests/limits", Vi: "request/limit tài nguyên"}, KeyPoints: BilingualText{En: "CPU/memory, QoS, OOMKills", Vi: "CPU/memory, QoS, OOMKill"}, Hint: BilingualText{En: "Explain right-sizing.", Vi: "Nêu right-sizing."}},
	}

	formats := []Format{
		{Level: "Junior", Title: BilingualText{En: "Kubernetes: {Concept}", Vi: "Kubernetes: {Concept}"}, Content: BilingualText{En: "Explain {Concept} in Kubernetes.", Vi: "Giải thích {Concept} trong Kubernetes."}, CorrectAnswer: BilingualText{En: "Cover: {KeyPoints}.", Vi: "Trình bày: {KeyPoints}."}, Hint: BilingualText{En: "Use a practical example.", Vi: "Dùng ví dụ thực tế."}},
		{Level: "Mid", Title: BilingualText{En: "{Concept} troubleshooting", Vi: "Troubleshoot {Concept}"}, Content: BilingualText{En: "How do you troubleshoot issues related to {Concept}?", Vi: "Bạn troubleshoot vấn đề liên quan {Concept} thế nào?"}, CorrectAnswer: BilingualText{En: "Describe steps and tools: {KeyPoints}.", Vi: "Nêu bước và công cụ: {KeyPoints}."}, Hint: BilingualText{En: "Use kubectl describe/logs.", Vi: "Dùng kubectl describe/logs."}},
		{Level: "Senior", Title: BilingualText{En: "{Concept} trade-offs", Vi: "Trade-off {Concept}"}, Content: BilingualText{En: "What trade-offs exist with {Concept} at scale?", Vi: "Trade-off của {Concept} khi scale là gì?"}, CorrectAnswer: BilingualText{En: "Discuss trade-offs and mitigations: {KeyPoints}.", Vi: "Nêu trade-off và giảm thiểu: {KeyPoints}."}, Hint: BilingualText{En: "Think reliability and cost.", Vi: "Nghĩ reliability và cost."}},
	}

	return BuildFromConcepts("Kubernetes", "DevOps", concepts, formats)
}

func BuildAwsTemplates() []QuestionTemplate {
	concepts := []Concept{
		{Name: BilingualText{En: "IAM fundamentals", Vi: "IAM cơ bản"}, KeyPoints: BilingualText{En: "roles, policies, least privilege", Vi: "role, policy, least privilege"}, Hint: BilingualText{En: "Use least privilege.", Vi: "Dùng least privilege."}},
		{Name: BilingualText{En: "VPC networking", Vi: "Networking VPC"}, KeyPoints: BilingualText{En: "subnets, routing, security groups, NACLs", Vi: "subnet, routing, security group, NACL"}, Hint: BilingualText{En: "Explain traffic flow.", Vi: "Nêu flow traffic."}},
		{Name: BilingualText{En: "load balancing", Vi: "load balancing"}, KeyPoints: BilingualText{En: "ALB/NLB, health checks, TLS", Vi: "ALB/NLB, health check, TLS"}, Hint: BilingualText{En: "Tie to application needs.", Vi: "Gắn với nhu cầu app."}},
		{Name: BilingualText{En: "storage options", Vi: "lựa chọn storage"}, KeyPoints: BilingualText{En: "S3/EBS/EFS, durability, cost", Vi: "S3/EBS/EFS, durability, cost"}, Hint: BilingualText{En: "Mention use cases.", Vi: "Nêu use case."}},
		{Name: BilingualText{En: "observability", Vi: "observability"}, KeyPoints: BilingualText{En: "CloudWatch metrics/logs/alarms, tracing", Vi: "CloudWatch metrics/logs/alarms, tracing"}, Hint: BilingualText{En: "Tie to SLOs.", Vi: "Gắn với SLO."}},
	}

	formats := []Format{
		{Level: "Junior", Title: BilingualText{En: "AWS: {Concept}", Vi: "AWS: {Concept}"}, Content: BilingualText{En: "Explain {Concept} in AWS.", Vi: "Giải thích {Concept} trong AWS."}, CorrectAnswer: BilingualText{En: "Cover: {KeyPoints}.", Vi: "Trình bày: {KeyPoints}."}, Hint: BilingualText{En: "Use a practical example.", Vi: "Dùng ví dụ thực tế."}},
		{Level: "Mid", Title: BilingualText{En: "{Concept} design", Vi: "Thiết kế {Concept}"}, Content: BilingualText{En: "How would you design a solution using {Concept}?", Vi: "Bạn thiết kế giải pháp dùng {Concept} thế nào?"}, CorrectAnswer: BilingualText{En: "Explain architecture and trade-offs: {KeyPoints}.", Vi: "Nêu kiến trúc và trade-off: {KeyPoints}."}, Hint: BilingualText{En: "Consider cost and reliability.", Vi: "Chú ý cost và reliability."}},
		{Level: "Senior", Title: BilingualText{En: "{Concept} pitfalls", Vi: "Pitfall {Concept}"}, Content: BilingualText{En: "What pitfalls exist with {Concept} and how do you mitigate them?", Vi: "Pitfall của {Concept} và cách giảm thiểu?"}, CorrectAnswer: BilingualText{En: "List pitfalls and mitigations: {KeyPoints}.", Vi: "Liệt kê pitfall và giảm thiểu: {KeyPoints}."}, Hint: BilingualText{En: "Think security and operations.", Vi: "Nghĩ bảo mật và vận hành."}},
	}

	return BuildFromConcepts("AWS", "DevOps", concepts, formats)
}
