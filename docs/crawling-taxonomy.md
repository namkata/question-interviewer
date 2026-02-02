# Data Crawling & Interview Question Taxonomy

## 1. Mục tiêu

Thu thập (crawl/scrape) dữ liệu câu hỏi phỏng vấn Software Engineer từ nhiều nguồn, chuẩn hóa, phân loại và đưa vào hệ thống Interview Q&A Platform.

---

## 2. Nguồn dữ liệu (Recommended Sources)

### Open / Public

* GitHub repositories (interview questions)
* Blog kỹ thuật (Medium, Dev.to)
* Engineering blogs (Google, Meta, Uber…)
* Public interview handbooks

### Community-driven

* StackOverflow (chỉ crawl câu hỏi, không copy nguyên câu trả lời)
* Reddit (r/cscareerquestions – summary)

⚠️ Lưu ý pháp lý:

* Chỉ crawl **question / topic / idea**, không copy nguyên answer có bản quyền
* Ưu tiên license: MIT, Apache-2.0, CC

---

## 3. Crawling Architecture

### Flow

1. FastAPI Crawler fetch raw data
2. Normalize content
3. Classify topic / level
4. Store into staging tables
5. Review / approve
6. Publish to main tables

### Components

* **Crawler (FastAPI + httpx / playwright)**
* **Parser / Normalizer**
* **Classifier (rule-based + AI)**
* **Admin Review UI**

---

## 4. Database – Staging Layer

### crawled_questions

* id (UUID)
* source
* raw_title
* raw_content
* url
* detected_topic
* detected_level
* status (pending / approved / rejected)
* created_at

Purpose:

* Không ghi thẳng vào questions
* Tránh dirty data

---

## 5. Question Taxonomy (Phân loại chuẩn)

### 5.1 By Level

* Junior
* Mid-level
* Senior
* Staff / Principal

---

### 5.2 By Category (Top-level)

#### 1. Computer Science Fundamentals

* Data Structures
* Algorithms
* Complexity Analysis
* Memory Management

#### 2. Backend Engineering

* API Design
* Concurrency
* Database Design
* Caching
* Messaging / Queue

#### 3. Frontend Engineering

* JavaScript Core
* React / Vue
* Browser Internals
* Performance

#### 4. System Design

* Scalability
* Load Balancing
* Consistency Models
* Distributed Systems

#### 5. DevOps / Infrastructure

* Docker
* Kubernetes
* CI/CD
* Cloud Architecture

#### 6. Language-specific

* Go
* Java
* Python
* JavaScript

---

## 6. Sample Question Set (Curated)

### Golang

* What is a goroutine?
* How does channel work internally?
* Mutex vs RWMutex
* Context cancellation patterns

### System Design

* Design a URL shortener
* Design a rate limiter
* Design a notification system

### Database

* Index vs Composite Index
* ACID vs BASE
* When to use read replica?

### Algorithms

* Difference between BFS and DFS
* Time complexity of quicksort
* How hash table works?

---

## 7. Classification Strategy

### Rule-based

* Keyword matching
* Source-based hints

### AI-assisted (Optional)

* LLM classify topic / level
* Deduplicate similar questions
* Generate canonical question title

---

## 8. Deduplication Strategy

* Hash normalized title
* Semantic similarity (embedding)
* Manual merge via admin UI

---

## 9. Review & Publish Workflow

1. Crawled → pending
2. Admin review
3. Edit & normalize
4. Publish
5. Sync search index

---

## 10. Roadmap

Phase 1:

* Manual curated questions

Phase 2:

* Automated crawling
* Rule-based classification

Phase 3:

* AI-assisted enrichment

Phase 4:

* Community contribution
