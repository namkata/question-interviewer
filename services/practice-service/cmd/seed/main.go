package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dsn := "host=localhost port=5432 user=user password=password dbname=question_db sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open db: %v", err)
	}
	defer db.Close()

	topics := []struct {
		Name        string
		Description string
	}{
		{"Python", "Python programming language questions"},
		{"Golang", "Golang programming language questions"},
		{"NodeJS", "Node.js runtime questions"},
		{"NestJS", "NestJS framework questions"},
		{"Behavioral", "Behavioral interview questions"},
		{"Leadership", "Leadership and management questions"},
		{"Algorithms", "Data structures and algorithms questions"},
		{"System Design", "System design questions"},
		{"Design Patterns", "Design Patterns 101 questions"},
	}

	for _, t := range topics {
		// Use uuid_generate_v4() as seen in existing migrations
		_, err := db.Exec("INSERT INTO topics (id, name, description) VALUES (uuid_generate_v4(), $1, $2) ON CONFLICT (name) DO NOTHING", t.Name, t.Description)
		if err != nil {
			log.Printf("Failed to insert topic %s: %v", t.Name, err)
		} else {
			fmt.Printf("Ensured topic: %s\n", t.Name)
		}
	}

	questions := []struct {
		Topic         string
		Title         string
		Content       string
		Level         string
		CorrectAnswer string
	}{
		{
			"Python",
			"Python Lists vs Tuples",
			"What is the difference between lists and tuples in Python?",
			"Junior",
			"Lists are mutable, meaning they can be changed after creation. Tuples are immutable and cannot be changed. Lists use square brackets [], while tuples use parentheses (). Tuples are generally faster and safer for fixed data.",
		},
		{
			"Python",
			"Python GIL",
			"What is the Global Interpreter Lock (GIL)?",
			"Senior",
			"The GIL is a mutex that allows only one thread to hold the control of the Python interpreter. This means that only one thread can be in a state of execution at any point in time. It effectively limits the parallelism of Python threads on multi-core systems.",
		},
		{
			"Golang",
			"Goroutines vs Threads",
			"Explain the difference between Goroutines and OS threads.",
			"Mid",
			"Goroutines are lightweight threads managed by the Go runtime, not the OS. They have a smaller stack size (starts at 2KB) which can grow/shrink. OS threads are heavier (typically 1MB stack) and managed by the kernel. Context switching is cheaper for Goroutines.",
		},
		{
			"Golang",
			"Interfaces",
			"How do interfaces work in Go? Explain duck typing.",
			"Junior",
			"Go interfaces are implemented implicitly. If a type provides the methods declared in an interface, it implements that interface. This is known as 'duck typing' - if it walks like a duck and quacks like a duck, it's a duck. No 'implements' keyword is needed.",
		},
		{
			"NodeJS",
			"Event Loop",
			"Explain the Node.js Event Loop.",
			"Mid",
			"The Event Loop is what allows Node.js to perform non-blocking I/O operations despite being single-threaded. It offloads operations to the system kernel whenever possible. It has phases like Timer, Pending Callbacks, Poll, Check, and Close Callbacks.",
		},
		{
			"NodeJS",
			"Buffer",
			"What is the purpose of the Buffer class in Node.js?",
			"Junior",
			"The Buffer class is used to handle binary data. Since JavaScript historically didn't have a mechanism for reading or manipulating streams of binary data, Buffer was introduced. It's now similar to Uint8Array but optimized for Node.js use cases.",
		},
		{
			"NestJS",
			"Dependency Injection",
			"How does Dependency Injection work in NestJS?",
			"Mid",
			"NestJS uses a built-in DI container. You define providers (services, repositories) and inject them into controllers or other services using constructor injection. The @Injectable() decorator marks a class as a provider.",
		},
		{
			"NestJS",
			"Guards vs Interceptors",
			"What is the difference between Guards and Interceptors in NestJS?",
			"Senior",
			"Guards determine whether a request should be handled by the route handler (authorization). Interceptors can intercept the request/response before/after the handler execution (logging, transformation, timeout). Guards run before Interceptors.",
		},
		{
			"Behavioral",
			"Conflict Resolution",
			"Tell me about a time you had a conflict with a coworker.",
			"Junior",
			"Use the STAR method: Situation (disagreement on API design), Task (needed to finalize spec), Action (scheduled a meeting, listed pros/cons, compromised), Result (delivered on time, improved relationship).",
		},
		{
			"Behavioral",
			"Weakness",
			"What is your greatest weakness?",
			"Junior",
			"Choose a real weakness but one you are working on. Example: 'I sometimes focus too much on details. I'm improving by setting strict time limits for tasks and focusing on the big picture first.'",
		},
		{
			"Leadership",
			"Motivating Team",
			"How do you motivate a team under tight deadlines?",
			"Senior",
			"I focus on transparency and purpose. I explain 'why' the deadline is important. I remove blockers, prioritize tasks to reduce scope creep, and ensure the team feels supported rather than pressured. Celebrating small wins is also key.",
		},
		{
			"Leadership",
			"Mentoring",
			"Describe a time you mentored a junior engineer.",
			"Senior",
			"I mentored a junior dev who struggled with testing. I set up pair programming sessions, reviewed their PRs with detailed explanations (not just corrections), and encouraged them to write the test plan before coding. They eventually became a testing advocate.",
		},
		{
			"Algorithms",
			"Binary Search",
			"Implement binary search on a sorted array.",
			"Mid",
			"Algorithm: Compare target with middle element. If equal, return index. If target < middle, search left half. If target > middle, search right half. Time Complexity: O(log n). Space: O(1) iterative.",
		},
		{
			"Algorithms",
			"Two Sum",
			"Given an array of integers and a target, return indices of the two numbers such that they add up to target.",
			"Junior",
			"Use a hash map to store the difference (target - current_value) and its index. Iterate through the array; if current value exists in map, return [map[current], current_index]. Time: O(n), Space: O(n).",
		},
		{
			"System Design",
			"Design Twitter",
			"Design a simplified version of Twitter.",
			"Senior",
			"Key components: User Service, Tweet Service, Timeline Service (Fan-out on write for active users, Fan-out on read for passive), Redis for caching timelines, Load Balancers. Database: SQL for users, NoSQL (Cassandra/DynamoDB) for tweets.",
		},
		{
			"System Design",
			"Design URL Shortener",
			"Design a URL shortening service like bit.ly.",
			"Senior",
			"Core: generate unique short ID (base62 encoding of database auto-inc ID or distributed ID generator like Snowflake). Store mapping (short_id -> long_url). 301 Redirect. High read/write ratio, heavy caching.",
		},
	}

	for _, q := range questions {
		// Get Topic ID
		var topicID string
		err := db.QueryRow("SELECT id FROM topics WHERE name = $1", q.Topic).Scan(&topicID)
		if err != nil {
			log.Printf("Skipping question for %s: topic not found", q.Topic)
			continue
		}

		// Insert Question
		_, err = db.Exec(`
			INSERT INTO questions (id, title, content, level, topic_id, created_by, status, correct_answer)
			VALUES (uuid_generate_v4(), $1, $2, $3, $4, '123e4567-e89b-12d3-a456-426614174000', 'published', $5)
			ON CONFLICT DO NOTHING
		`, q.Title, q.Content, q.Level, topicID, q.CorrectAnswer)

		if err != nil {
			log.Printf("Failed to insert question '%s': %v", q.Title, err)
		} else {
			fmt.Printf("Ensured question: %s\n", q.Title)
		}
	}
}
