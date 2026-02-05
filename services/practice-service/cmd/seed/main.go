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
		{"CV Screening", "Initial screening questions about experience and background"},
		{"Database", "Database and data layer questions"},
		{"Testing", "Testing and quality assurance questions"},
		{"DevOps", "General DevOps questions"},
		{"Docker", "Containerization with Docker"},
		{"Kubernetes", "Container Orchestration with K8s"},
		{"CI/CD", "Continuous Integration and Deployment"},
		{"Terraform", "Infrastructure as Code"},
		{"React", "React library questions"},
		{"Vue", "Vue.js framework questions"},
		{"Frontend Basic", "HTML, CSS, JavaScript questions"},
		{"Data Engineering", "Data Engineering, SQL, Big Data questions"},
		{"Java", "Java programming language questions"},
		{"JavaScript", "JavaScript language core questions"},
		{"Network", "Computer Network, HTTP, TCP/IP questions"},
		{"CSS", "CSS Styling, Layouts, Responsive Design questions"},
	}

	for _, t := range topics {
		_, err := db.Exec("INSERT INTO topics (id, name, description) VALUES (uuid_generate_v4(), $1, $2) ON CONFLICT (name) DO NOTHING", t.Name, t.Description)
		if err != nil {
			log.Printf("Failed to insert topic %s: %v", t.Name, err)
		} else {
			fmt.Printf("Ensured topic: %s\n", t.Name)
		}
	}

	// Clean up existing data to start fresh
	// Delete in order of dependencies: practice_attempts -> practice_sessions -> answers -> questions
	// We also clean legacy tables 'sessions' just in case.
	tables := []string{"practice_attempts", "practice_sessions", "answers", "sessions", "questions"}
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			// Don't fail if table doesn't exist, just log it
			log.Printf("Note: Failed to delete from %s (might not exist or other error): %v", table, err)
		} else {
			fmt.Printf("Cleaned up table: %s\n", table)
		}
	}

	questions := []struct {
		Topic         string
		Title         string
		Content       string
		Level         string
		CorrectAnswer string
		Language      string
		Role          string
		Hint          string
	}{
		{
			"Network",
			"TCP vs UDP",
			"What is the difference between TCP and UDP?",
			"Junior",
			"TCP: Connection-oriented, reliable (retries, ordering), slower (handshake). UDP: Connectionless, unreliable (fire and forget), faster. Used for streaming/gaming.",
			"en",
			"DevOps",
			"Think about reliability vs speed. TCP guarantees delivery; UDP does not.",
		},
		{
			"Network",
			"TCP vs UDP",
			"Khác biệt giữa TCP và UDP?",
			"Junior",
			"TCP: Hướng kết nối, tin cậy (gửi lại, đúng thứ tự), chậm hơn. UDP: Không kết nối, không tin cậy, nhanh hơn. Dùng cho streaming/game.",
			"vi",
			"DevOps",
			"Nghĩ về độ tin cậy so với tốc độ. TCP đảm bảo gửi đến nơi; UDP thì không.",
		},
		{
			"Network",
			"HTTP Status Codes",
			"Common HTTP status codes?",
			"Junior",
			"200 OK, 201 Created, 400 Bad Request, 401 Unauthorized, 403 Forbidden, 404 Not Found, 500 Internal Error, 502 Bad Gateway.",
			"en",
			"DevOps",
			"Group them: 2xx (Success), 3xx (Redirect), 4xx (Client Error), 5xx (Server Error).",
		},
		{
			"Network",
			"Mã trạng thái HTTP",
			"Các mã trạng thái HTTP phổ biến?",
			"Junior",
			"200 OK, 201 Created, 400 Bad Request, 401 Unauthorized, 403 Forbidden, 404 Not Found, 500 Internal Error, 502 Bad Gateway.",
			"vi",
			"DevOps",
			"Nhóm chúng lại: 2xx (Thành công), 3xx (Chuyển hướng), 4xx (Lỗi Client), 5xx (Lỗi Server).",
		},
		// --- Network (BackEnd/DevOps) ---
		{
			"Network",
			"TCP vs UDP",
			"What is the difference between TCP and UDP?",
			"Junior",
			"TCP: Connection-oriented, reliable (retries, ordering), slower (handshake). UDP: Connectionless, unreliable (fire and forget), faster. Used for streaming/gaming.",
			"en",
			"BackEnd",
			"Think about reliability vs speed. TCP guarantees delivery; UDP does not.",
		},
		{
			"Network",
			"TCP vs UDP",
			"Khác biệt giữa TCP và UDP?",
			"Junior",
			"TCP: Hướng kết nối, tin cậy (gửi lại, đúng thứ tự), chậm hơn. UDP: Không kết nối, không tin cậy, nhanh hơn. Dùng cho streaming/game.",
			"vi",
			"BackEnd",
			"Nghĩ về độ tin cậy so với tốc độ. TCP đảm bảo gửi đến nơi; UDP thì không.",
		},
		{
			"Network",
			"HTTP Status Codes",
			"Common HTTP status codes?",
			"Junior",
			"200 OK, 201 Created, 400 Bad Request, 401 Unauthorized, 403 Forbidden, 404 Not Found, 500 Internal Error, 502 Bad Gateway.",
			"en",
			"BackEnd",
			"Group them: 2xx (Success), 3xx (Redirect), 4xx (Client Error), 5xx (Server Error).",
		},
		{
			"Network",
			"Mã trạng thái HTTP",
			"Các mã trạng thái HTTP phổ biến?",
			"Junior",
			"200 OK, 201 Created, 400 Bad Request, 401 Unauthorized, 403 Forbidden, 404 Not Found, 500 Internal Error, 502 Bad Gateway.",
			"vi",
			"BackEnd",
			"Nhóm chúng lại: 2xx (Thành công), 3xx (Chuyển hướng), 4xx (Lỗi Client), 5xx (Lỗi Server).",
		},

		// --- Frontend CSS (Junior) ---
		{
			"CSS",
			"Flexbox vs Grid",
			"Difference between Flexbox and Grid?",
			"Junior",
			"Flexbox is 1D (row OR column). Grid is 2D (rows AND columns). Flexbox is content-first, Grid is layout-first.",
			"en",
			"FrontEnd",
			"One dimension vs Two dimensions. Flexbox for components, Grid for page layout.",
		},
		{
			"CSS",
			"Flexbox vs Grid",
			"Khác biệt giữa Flexbox và Grid?",
			"Junior",
			"Flexbox là 1 chiều (dòng HOẶC cột). Grid là 2 chiều (dòng VÀ cột). Flexbox ưu tiên nội dung, Grid ưu tiên bố cục.",
			"vi",
			"FrontEnd",
			"Một chiều vs Hai chiều. Flexbox dùng cho component nhỏ, Grid dùng cho bố cục trang.",
		},
		{
			"CSS",
			"Specificity",
			"How is CSS Specificity calculated?",
			"Junior",
			"Inline styles > IDs > Classes/Attributes/Pseudo-classes > Elements/Pseudo-elements. Calculated as (Inline, ID, Class, Tag).",
			"en",
			"FrontEnd",
			"Hierarchy: Inline > ID > Class > Tag. !important overrides everything.",
		},
		{
			"CSS",
			"Độ ưu tiên CSS",
			"Độ ưu tiên (Specificity) trong CSS được tính thế nào?",
			"Junior",
			"Inline styles > ID > Class > Tag. Được tính theo trọng số (Inline, ID, Class, Tag). !important ghi đè tất cả.",
			"vi",
			"FrontEnd",
			"Thứ bậc: Inline > ID > Class > Tag. !important mạnh nhất.",
		},

		// --- Frontend JavaScript (Junior) ---
		{
			"JavaScript",
			"Closures",
			"What is a Closure in JavaScript?",
			"Junior",
			"A function bundled with its lexical environment. It allows a function to access variables from its outer scope even after the outer function has finished executing.",
			"en",
			"FrontEnd",
			"It's about scope retention. A function 'remembers' variables from where it was created.",
		},
		{
			"JavaScript",
			"Closures",
			"Closure trong JavaScript là gì?",
			"Junior",
			"Một hàm đi kèm với môi trường định nghĩa của nó (lexical environment). Nó cho phép hàm truy cập biến từ phạm vi bên ngoài ngay cả khi hàm bên ngoài đã chạy xong.",
			"vi",
			"FrontEnd",
			"Nó là về việc giữ lại phạm vi. Hàm 'nhớ' các biến từ nơi nó được tạo ra.",
		},
		{
			"JavaScript",
			"Async/Await",
			"Explain Async/Await.",
			"Junior",
			"Syntactic sugar over Promises. Makes asynchronous code look and behave like synchronous code. 'await' pauses execution until the Promise resolves.",
			"en",
			"FrontEnd",
			"It's just a cleaner way to write Promises. Makes code readable.",
		},
		{
			"JavaScript",
			"Async/Await",
			"Giải thích Async/Await.",
			"Junior",
			"Cú pháp viết tắt (syntactic sugar) cho Promise. Giúp code bất đồng bộ trông giống và chạy giống code đồng bộ. 'await' tạm dừng thực thi cho đến khi Promise hoàn thành.",
			"vi",
			"FrontEnd",
			"Cách viết Promise gọn gàng hơn. Giúp code dễ đọc.",
		},

		// --- CV Screening ---
		{
			"CV Screening",
			"Recent Backend Project",
			"Describe your most recent backend project. What was the scale (Users, TPS, Data Size)?",
			"Junior",
			"Look for specific metrics (e.g., 10k DAU, 1000 TPS, 5TB data). Candidate should explain the architecture, their specific role, and how they handled the scale.",
			"en",
			"BackEnd",
			"Structure your answer: 1. Project Overview (What it does) 2. Your Role (What you built) 3. Scale Metrics (Users, Requests, Data) 4. Key Challenges & Solutions.",
		},
		{
			"CV Screening",
			"Dự án Backend gần nhất",
			"Dự án backend gần nhất anh/chị/bạn làm cái gì? Quy mô user, TPS, data size là bao nhiêu?",
			"Junior",
			"Ứng viên cần nêu rõ các con số cụ thể (ví dụ: 10k user hàng ngày, 1000 request/giây, dữ liệu 5TB). Cần mô tả kiến trúc, vai trò cụ thể và cách xử lý quy mô đó.",
			"vi",
			"BackEnd",
			"Cấu trúc câu trả lời: 1. Tổng quan dự án (Nó làm gì) 2. Vai trò của bạn (Bạn làm gì) 3. Số liệu quy mô (Người dùng, Request, Dữ liệu) 4. Thách thức & Giải pháp chính.",
		},
		{
			"CV Screening",
			"Recent Frontend Project",
			"Describe your most recent frontend project. What frameworks did you use and what were the main challenges?",
			"Junior",
			"Look for framework knowledge (React, Vue, etc.), state management complexity, performance optimization, or complex UI interactions.",
			"en",
			"FrontEnd",
			"Structure your answer: 1. Project Overview 2. Tech Stack (React, Redux, etc.) 3. Key Feature you built 4. Hardest bug/challenge you solved.",
		},
		{
			"CV Screening",
			"Dự án Frontend gần nhất",
			"Dự án frontend gần nhất bạn làm là gì? Bạn dùng framework nào và thách thức lớn nhất là gì?",
			"Junior",
			"Tìm kiếm kiến thức về framework (React, Vue...), độ phức tạp quản lý state, tối ưu hiệu năng, hoặc các tương tác UI phức tạp.",
			"vi",
			"FrontEnd",
			"Cấu trúc câu trả lời: 1. Tổng quan dự án 2. Tech Stack (React, Redux...) 3. Tính năng chính bạn làm 4. Lỗi/Thách thức khó nhất bạn đã giải quyết.",
		},

		// --- CV Screening (DevOps) ---
		{
			"CV Screening",
			"Recent DevOps Project",
			"Describe your most recent DevOps project/infrastructure. What tools did you use?",
			"Junior",
			"Look for Infrastructure as Code (Terraform/Ansible), CI/CD (Jenkins/GitLab), Containerization (Docker/K8s), and Cloud Provider (AWS/GCP).",
			"en",
			"DevOps",
			"Structure: 1. Project/Infra Overview 2. Your Role 3. Tools Used 4. Improvements you made (automation, cost, speed).",
		},
		{
			"CV Screening",
			"Dự án DevOps gần nhất",
			"Mô tả dự án DevOps hoặc hạ tầng gần nhất bạn làm. Bạn đã dùng công cụ gì?",
			"Junior",
			"Tìm kiếm kiến thức về IaC (Terraform/Ansible), CI/CD, Container (Docker/K8s), và Cloud (AWS/GCP).",
			"vi",
			"DevOps",
			"Cấu trúc: 1. Tổng quan hạ tầng 2. Vai trò 3. Công cụ 4. Cải tiến bạn đã làm (tự động hóa, chi phí, tốc độ).",
		},

		// --- CV Screening (Data Engineer) ---
		{
			"CV Screening",
			"Recent Data Project",
			"Describe your most recent data pipeline or project.",
			"Junior",
			"Look for ETL tools, Big Data tech (Spark, Kafka), Warehousing (Redshift, BigQuery), and scale of data.",
			"en",
			"Data Engineer",
			"Structure: 1. Pipeline Overview 2. Your Role 3. Tech Stack 4. Data Volume & Latency requirements.",
		},
		{
			"CV Screening",
			"Dự án Data gần nhất",
			"Mô tả pipeline dữ liệu hoặc dự án data gần nhất của bạn.",
			"Junior",
			"Tìm kiếm công cụ ETL, Big Data (Spark, Kafka), Warehouse, và quy mô dữ liệu.",
			"vi",
			"Data Engineer",
			"Cấu trúc: 1. Tổng quan Pipeline 2. Vai trò 3. Tech Stack 4. Khối lượng dữ liệu & Yêu cầu độ trễ.",
		},

		// --- Frontend Basic (Fresher/Junior) ---
		{
			"Frontend Basic",
			"Box Model",
			"Explain the CSS Box Model.",
			"Fresher",
			"Content, Padding, Border, Margin. Standard vs Border-Box sizing.",
			"en",
			"FrontEnd",
			"Think about the layers wrapping an HTML element. From inside out: Content -> ? -> ? -> ?.",
		},
		{
			"Frontend Basic",
			"Box Model",
			"Giải thích CSS Box Model.",
			"Fresher",
			"Gồm Content, Padding, Border, Margin. Phân biệt box-sizing: content-box và border-box.",
			"vi",
			"FrontEnd",
			"Hãy nghĩ về các lớp bao quanh một phần tử HTML. Từ trong ra ngoài: Nội dung -> ? -> ? -> ?.",
		},
		{
			"Frontend Basic",
			"Let vs Var vs Const",
			"What is the difference between let, var, and const in JavaScript?",
			"Fresher",
			"var: function scoped, hoisted. let: block scoped, can reassign. const: block scoped, cannot reassign (but objects are mutable).",
			"en",
			"FrontEnd",
			"Focus on two main aspects: Scope (Block vs Function) and Reassignment (Mutable vs Immutable).",
		},
		{
			"Frontend Basic",
			"Let vs Var vs Const",
			"Sự khác biệt giữa let, var và const trong JavaScript là gì?",
			"Fresher",
			"var: phạm vi hàm (function scope), hoisted. let: phạm vi khối (block scope), có thể gán lại. const: phạm vi khối, không thể gán lại (nhưng object vẫn mutable).",
			"vi",
			"FrontEnd",
			"Tập trung vào hai khía cạnh chính: Phạm vi (Block vs Function) và Khả năng gán lại (Mutable vs Immutable).",
		},

		// --- Frontend React (Junior/Mid) ---
		{
			"React",
			"React Hooks",
			"What are React Hooks? Name common ones.",
			"Fresher",
			"Hooks allow using state and other React features in functional components. Common ones: useState, useEffect, useContext, useRef.",
			"en",
			"FrontEnd",
			"They let you use state and other features without writing a class. Think about managing data (state) and side effects.",
		},
		{
			"React",
			"React Hooks",
			"React Hooks là gì? Kể tên vài hook phổ biến.",
			"Fresher",
			"Hooks cho phép dùng state và các tính năng React trong functional components. Phổ biến: useState, useEffect, useContext, useRef.",
			"vi",
			"FrontEnd",
			"Chúng cho phép bạn sử dụng state và các tính năng khác mà không cần viết class. Hãy nghĩ về quản lý dữ liệu (state) và hiệu ứng phụ (side effects).",
		},
		{
			"React",
			"Virtual DOM",
			"How does Virtual DOM work in React?",
			"Junior",
			"React creates a lightweight copy of the DOM. When state changes, it creates a new Virtual DOM tree, compares it with the previous one (diffing), and only updates the changed parts in the real DOM (reconciliation).",
			"en",
			"FrontEnd",
			"Keywords: Lightweight copy, Diffing algorithm, Reconciliation, Batch updates.",
		},
		{
			"React",
			"Virtual DOM",
			"Virtual DOM trong React hoạt động như thế nào?",
			"Junior",
			"React tạo một bản sao nhẹ của DOM. Khi state thay đổi, nó tạo cây Virtual DOM mới, so sánh với cây cũ (diffing), và chỉ cập nhật những phần thay đổi vào DOM thật (reconciliation).",
			"vi",
			"FrontEnd",
			"Từ khóa: Bản sao nhẹ (Lightweight copy), Thuật toán so sánh (Diffing), Đối chiếu (Reconciliation), Cập nhật hàng loạt.",
		},
		{
			"React",
			"useEffect Dependency Array",
			"What happens if you leave the dependency array empty in useEffect?",
			"Junior",
			"The effect runs only once after the initial render, similar to componentDidMount.",
			"en",
			"FrontEnd",
			"Compare it to the lifecycle methods of a class component (Mounting).",
		},
		{
			"React",
			"useEffect Dependency Array",
			"Điều gì xảy ra nếu bạn để trống mảng dependency trong useEffect?",
			"Junior",
			"Effect chỉ chạy một lần duy nhất sau lần render đầu tiên, tương tự như componentDidMount.",
			"vi",
			"FrontEnd",
			"So sánh nó với các phương thức vòng đời của class component (Mounting).",
		},

		// --- Backend Golang (Junior/Mid/Senior) ---
		{
			"Golang",
			"Context Usage",
			"What is context.Context used for? When should you cancel it?",
			"Mid",
			"Context is used for deadline propagation, cancellation signals, and request-scoped values across API boundaries. You should cancel it when the operation is no longer needed (e.g., client disconnected, timeout reached) to free up resources.",
			"en",
			"BackEnd",
			"Think about controlling long-running processes and passing request-scoped data. Keywords: Timeout, Cancellation, Values.",
		},
		{
			"Golang",
			"Sử dụng Context",
			"context.Context dùng để làm gì? Khi nào nên cancel?",
			"Mid",
			"Context dùng để truyền thời hạn (deadline), tín hiệu hủy (cancellation) và các giá trị trong phạm vi request. Nên cancel khi thao tác không còn cần thiết (ví dụ: client ngắt kết nối, hết thời gian chờ) để giải phóng tài nguyên.",
			"vi",
			"BackEnd",
			"Hãy nghĩ về việc kiểm soát các tiến trình chạy lâu và truyền dữ liệu trong phạm vi request. Từ khóa: Timeout, Hủy (Cancellation), Giá trị (Values).",
		},
		{
			"Golang",
			"Goroutines vs Threads",
			"How are Goroutines different from OS Threads?",
			"Junior",
			"Goroutines are lightweight, managed by Go runtime (M:N scheduling), have smaller stack size (starts at 2KB), and cheaper context switching compared to OS threads.",
			"en",
			"BackEnd",
			"Compare: Memory usage (Stack size), Scheduling (OS vs Runtime), and Creation cost.",
		},
		{
			"Golang",
			"Goroutines vs Threads",
			"Goroutines khác gì với OS Threads?",
			"Junior",
			"Goroutines nhẹ hơn, được quản lý bởi Go runtime (lịch trình M:N), kích thước stack nhỏ (bắt đầu 2KB), và chuyển đổi ngữ cảnh (context switch) rẻ hơn nhiều so với OS threads.",
			"vi",
			"BackEnd",
			"So sánh: Bộ nhớ (Kích thước Stack), Lập lịch (OS vs Runtime), và Chi phí tạo mới.",
		},
		{
			"Golang",
			"Channels",
			"What is the difference between buffered and unbuffered channels?",
			"Mid",
			"Unbuffered: Sender blocks until receiver is ready (synchronous). Buffered: Sender blocks only when buffer is full (asynchronous up to capacity).",
			"en",
			"BackEnd",
			"Think about blocking behavior. When does the sender stop waiting in each case?",
		},
		{
			"Golang",
			"Channels",
			"Sự khác biệt giữa buffered và unbuffered channels?",
			"Mid",
			"Unbuffered: Người gửi bị chặn (block) cho đến khi người nhận sẵn sàng (đồng bộ). Buffered: Người gửi chỉ bị chặn khi buffer đầy (bất đồng bộ trong giới hạn).",
			"vi",
			"BackEnd",
			"Hãy nghĩ về hành vi chặn (blocking). Khi nào người gửi ngừng chờ đợi trong mỗi trường hợp?",
		},

		// --- Backend NodeJS (Junior/Mid) ---
		{
			"NodeJS",
			"Event Loop",
			"Explain the Node.js Event Loop.",
			"Junior",
			"Single-threaded loop that handles asynchronous callbacks. Phases: Timers, Pending Callbacks, Poll, Check (setImmediate), Close Callbacks. Offloads heavy tasks to libuv worker pool.",
			"en",
			"BackEnd",
			"It's the mechanism that allows Node.js to perform non-blocking I/O operations. Mention the phases and libuv.",
		},
		{
			"NodeJS",
			"Event Loop",
			"Giải thích Node.js Event Loop.",
			"Junior",
			"Vòng lặp đơn luồng xử lý các callback bất đồng bộ. Các giai đoạn: Timers, Pending Callbacks, Poll, Check (setImmediate), Close Callbacks. Tác vụ nặng được đẩy xuống libuv worker pool.",
			"vi",
			"BackEnd",
			"Đó là cơ chế cho phép Node.js thực hiện các thao tác I/O không chặn. Hãy đề cập đến các giai đoạn và libuv.",
		},

		// --- Database (Mid/Senior) ---
		{
			"Database",
			"Redis Use Cases",
			"What are the common use cases for Redis?",
			"Mid",
			"Caching (reducing DB load), Session Store, Pub/Sub messaging, Leaderboards (Sorted Sets), Rate Limiting, Queues.",
			"en",
			"BackEnd",
			"Think about scenarios requiring high speed and temporary data storage. Not just Caching.",
		},
		{
			"Database",
			"Redis Use Cases",
			"Redis dùng cho những case nào?",
			"Mid",
			"Caching (giảm tải DB), Lưu trữ Session, Pub/Sub (bắn tin), Bảng xếp hạng (Sorted Sets), Giới hạn tốc độ (Rate Limiting), Hàng đợi (Queues).",
			"vi",
			"BackEnd",
			"Hãy nghĩ về các kịch bản yêu cầu tốc độ cao và lưu trữ dữ liệu tạm thời. Không chỉ là Caching.",
		},
		{
			"Database",
			"Indexing",
			"How does a database index work? Pros and cons?",
			"Mid",
			"Uses data structures like B-Trees to allow fast lookup (O(log n)) instead of full table scan. Pros: Faster reads. Cons: Slower writes (insert/update), increased storage usage.",
			"en",
			"BackEnd",
			"Trade-off between Read speed and Write speed/Storage.",
		},
		{
			"Database",
			"Indexing",
			"Index trong database hoạt động thế nào? Ưu nhược điểm?",
			"Mid",
			"Dùng cấu trúc dữ liệu như B-Tree để tìm kiếm nhanh (O(log n)) thay vì quét toàn bộ bảng. Ưu: Đọc nhanh. Nhược: Ghi chậm (insert/update), tốn dung lượng lưu trữ.",
			"vi",
			"BackEnd",
			"Sự đánh đổi giữa tốc độ Đọc và tốc độ Ghi/Lưu trữ.",
		},

		// --- System Design (Senior) ---
		{
			"System Design",
			"Rate Limiting",
			"How would you limit 10 concurrent uploads per user?",
			"Senior",
			"Use a distributed counter (e.g., Redis) or a Semaphore pattern. When upload starts, increment counter/acquire lock. If > 10, reject. When finished, decrement/release. Handle race conditions and timeouts.",
			"en",
			"BackEnd",
			"You need a shared state store to track active uploads across servers. Redis is a common choice.",
		},
		{
			"System Design",
			"Rate Limiting",
			"Làm sao limit 10 concurrent upload / user?",
			"Senior",
			"Dùng bộ đếm phân tán (ví dụ: Redis) hoặc mẫu Semaphore. Khi bắt đầu upload, tăng biến đếm. Nếu > 10, từ chối. Khi xong, giảm biến đếm. Cần xử lý race condition và timeout.",
			"vi",
			"BackEnd",
			"Bạn cần một kho lưu trữ trạng thái chia sẻ để theo dõi các upload đang hoạt động trên các server. Redis là lựa chọn phổ biến.",
		},
		{
			"System Design",
			"CAP Theorem",
			"Explain CAP Theorem.",
			"Senior",
			"Consistency, Availability, Partition Tolerance. In a distributed system, you can only pick two. P is mandatory in network systems, so usually choose between CP (strong consistency) or AP (high availability).",
			"en",
			"BackEnd",
			"Acronyms: C (Consistency), A (Availability), P (Partition Tolerance). Pick two.",
		},
		{
			"System Design",
			"Định lý CAP",
			"Giải thích định lý CAP.",
			"Senior",
			"Consistency (Tính nhất quán), Availability (Tính sẵn sàng), Partition Tolerance (Khả năng chịu lỗi phân vùng). Trong hệ thống phân tán, chỉ chọn được 2. P là bắt buộc, nên thường chọn giữa CP (nhất quán mạnh) hoặc AP (sẵn sàng cao).",
			"vi",
			"BackEnd",
			"Viết tắt: C (Nhất quán), A (Sẵn sàng), P (Chịu lỗi phân vùng). Chỉ chọn được hai.",
		},

		// --- Algorithms (Mid) ---
		{
			"Algorithms",
			"Time Complexity",
			"How do you determine the time complexity of a code block?",
			"Mid",
			"Analyze loops and recursion. Single loop is O(n), nested is O(n^2). Binary search is O(log n). Look for dominant operations as input size grows.",
			"en",
			"BackEnd",
			"Count the nested loops. Look for recursion. Think about Big O notation.",
		},
		{
			"Algorithms",
			"Độ phức tạp thời gian",
			"Time complexity của đoạn code là gì? Cách xác định?",
			"Mid",
			"Phân tích vòng lặp và đệ quy. Vòng lặp đơn là O(n), lồng nhau là O(n^2). Tìm kiếm nhị phân là O(log n). Quan tâm đến các thao tác chiếm ưu thế khi kích thước đầu vào tăng.",
			"vi",
			"BackEnd",
			"Đếm số vòng lặp lồng nhau. Tìm đệ quy. Nghĩ về Big O notation.",
		},

		// --- Frontend System Design (Senior) ---
		{
			"System Design",
			"Infinite Scroll",
			"Design an Infinite Scroll component. Key challenges?",
			"Senior",
			"Scroll event throttling, Virtualization (windowing) to maintain DOM size, Fetching strategy (prefetching), Error handling, Restoration of scroll position.",
			"en",
			"FrontEnd",
			"Think about Performance (DOM nodes) and User Experience (Loading states).",
		},
		{
			"System Design",
			"Cuộn vô tận (Infinite Scroll)",
			"Thiết kế component Infinite Scroll. Thách thức chính là gì?",
			"Senior",
			"Throttling sự kiện scroll, Virtualization (windowing) để giữ kích thước DOM nhỏ, Chiến lược tải (prefetching), Xử lý lỗi, Khôi phục vị trí cuộn.",
			"vi",
			"FrontEnd",
			"Hãy nghĩ về Hiệu năng (Số lượng DOM node) và Trải nghiệm người dùng (Trạng thái tải).",
		},

		// --- Frontend Algorithms (Mid) ---
		{
			"Algorithms",
			"Debounce vs Throttle",
			"Implement Debounce and Throttle functions.",
			"Mid",
			"Debounce: Delay execution until X ms have passed since last call (search bar). Throttle: Ensure execution at most once every X ms (scroll event).",
			"en",
			"FrontEnd",
			"Delaying execution vs Limiting execution frequency.",
		},
		{
			"Algorithms",
			"Debounce vs Throttle",
			"Triển khai hàm Debounce và Throttle.",
			"Mid",
			"Debounce: Trì hoãn thực thi cho đến khi X ms trôi qua từ lần gọi cuối (thanh tìm kiếm). Throttle: Đảm bảo thực thi tối đa 1 lần mỗi X ms (sự kiện cuộn).",
			"vi",
			"FrontEnd",
			"Trì hoãn thực thi vs Giới hạn tần suất thực thi.",
		},

		// --- Testing (Senior) ---
		{
			"Testing",
			"Mocking DB",
			"Is MockDB reliable? Pros and cons?",
			"Senior",
			"MockDB is good for unit tests (fast, isolated) but not fully reliable for behavior (doesn't catch constraints, triggers, complex query issues). Integration tests with real DB (e.g., Testcontainers) are preferred for data layer testing.",
			"en",
			"BackEnd",
			"Distinguish between Unit Tests (Isolation/Speed) and Integration Tests (Reliability/Real Environment).",
		},
		{
			"Testing",
			"MockDB",
			"MockDB có đáng tin không?",
			"Senior",
			"MockDB tốt cho unit test (nhanh, cô lập) nhưng không đáng tin cậy hoàn toàn về hành vi (không bắt được lỗi ràng buộc, trigger, query phức tạp). Nên ưu tiên Integration test với DB thật (ví dụ: Docker/Testcontainers) cho data layer.",
			"vi",
			"BackEnd",
			"Phân biệt giữa Unit Test (Cô lập/Nhanh) và Integration Test (Độ tin cậy/Môi trường thật).",
		},

		// --- Docker (Junior/Mid) ---
		{
			"Docker",
			"Docker Experience",
			"Have you worked with Docker? Explain a basic Dockerfile.",
			"Junior",
			"Expect familiarity with FROM, RUN, COPY, CMD/ENTRYPOINT. Understanding of layers, image vs container, and volume mounting.",
			"en",
			"DevOps",
			"Keywords: Image, Container, Dockerfile instructions (FROM, RUN, CMD).",
		},
		{
			"Docker",
			"Kinh nghiệm Docker",
			"Anh/chị đã từng làm việc với Docker chưa?",
			"Junior",
			"Mong đợi ứng viên biết về FROM, RUN, COPY, CMD. Hiểu về các layer, sự khác biệt giữa image và container, volume.",
			"vi",
			"DevOps",
			"Từ khóa: Image, Container, các lệnh Dockerfile (FROM, RUN, CMD).",
		},

		// --- CI/CD (Mid) ---
		{
			"CI/CD",
			"CI/CD Pipeline",
			"What are the stages of a standard CI/CD pipeline?",
			"Mid",
			"Code (Commit), Build (Compile/Dockerize), Test (Unit/Integration), Release (Version/Tag), Deploy (Staging/Prod), Monitor.",
			"en",
			"DevOps",
			"Think about the flow from code commit to production. Build -> Test -> Deploy.",
		},
		{
			"CI/CD",
			"Quy trình CI/CD",
			"Các giai đoạn chuẩn của một pipeline CI/CD là gì?",
			"Mid",
			"Code (Commit), Build (Compile/Dockerize), Test (Unit/Integration), Release (Version/Tag), Deploy (Staging/Prod), Monitor.",
			"vi",
			"DevOps",
			"Hãy nghĩ về quy trình từ khi commit code đến khi lên production. Build -> Test -> Deploy.",
		},

		// --- Kubernetes (Mid/Senior) ---
		{
			"Kubernetes",
			"K8s Components",
			"What are the main components of Kubernetes?",
			"Mid",
			"Control Plane (API Server, etcd, Scheduler, Controller Manager). Nodes (Kubelet, Kube-proxy, Container Runtime).",
			"en",
			"DevOps",
			"Control Plane vs Worker Node components.",
		},
		{
			"Kubernetes",
			"Thành phần Kubernetes",
			"Các thành phần chính của Kubernetes là gì?",
			"Mid",
			"Control Plane (API Server, etcd, Scheduler, Controller Manager). Nodes (Kubelet, Kube-proxy, Container Runtime).",
			"vi",
			"DevOps",
			"Phân biệt thành phần Control Plane và Worker Node.",
		},

		// --- Terraform (Mid/Senior) ---
		{
			"Terraform",
			"Terraform State",
			"What is Terraform State and why is it important?",
			"Mid",
			"State file tracks the mapping between configuration and real-world resources. It stores metadata, improves performance, and enables locking (to prevent concurrent updates).",
			"en",
			"DevOps",
			"It's the 'source of truth' for Terraform about what exists in the cloud.",
		},
		{
			"Terraform",
			"Terraform State",
			"Terraform State là gì và tại sao nó quan trọng?",
			"Mid",
			"File State theo dõi ánh xạ giữa cấu hình và tài nguyên thực tế. Nó lưu metadata, cải thiện hiệu năng và cho phép khóa (locking) để tránh cập nhật đồng thời.",
			"vi",
			"DevOps",
			"Nó là 'nguồn sự thật' của Terraform về những gì đang tồn tại trên cloud.",
		},

		// --- Data Engineering (Fresher/Junior/Mid) ---
		{
			"Data Engineering",
			"SQL Group By",
			"What does GROUP BY do in SQL?",
			"Fresher",
			"Groups rows that have the same values into summary rows, often used with aggregate functions (COUNT, MAX, MIN, SUM, AVG).",
			"en",
			"Data Engineer",
			"It combines rows with identical values. Usually used with COUNT, SUM, AVG.",
		},
		{
			"Data Engineering",
			"SQL Group By",
			"Lệnh GROUP BY trong SQL dùng để làm gì?",
			"Fresher",
			"Nhóm các hàng có cùng giá trị thành các hàng tóm tắt, thường dùng với các hàm tổng hợp (COUNT, MAX, MIN, SUM, AVG).",
			"vi",
			"Data Engineer",
			"Nó gộp các hàng có giá trị giống nhau. Thường dùng với COUNT, SUM, AVG.",
		},
		{
			"Data Engineering",
			"ETL Process",
			"What is ETL?",
			"Junior",
			"Extract, Transform, Load. The process of copying data from various sources into a destination system which represents the data differently from the source.",
			"en",
			"Data Engineer",
			"Acronym: E (Extract), T (Transform), L (Load).",
		},
		{
			"Data Engineering",
			"Quy trình ETL",
			"ETL là gì?",
			"Junior",
			"Extract (Trích xuất), Transform (Chuyển đổi), Load (Tải). Quy trình sao chép dữ liệu từ nhiều nguồn vào hệ thống đích với định dạng hoặc cấu trúc khác.",
			"vi",
			"Data Engineer",
			"Viết tắt: E (Extract - Trích xuất), T (Transform - Chuyển đổi), L (Load - Tải).",
		},
		{
			"Data Engineering",
			"Data Warehouse vs Data Lake",
			"Difference between Data Warehouse and Data Lake?",
			"Mid",
			"Warehouse: Structured data, schema-on-write, optimized for analysis (OLAP). Lake: Raw data (structured/unstructured), schema-on-read, low cost storage.",
			"en",
			"Data Engineer",
			"Structured vs Unstructured. Schema-on-write vs Schema-on-read.",
		},
		{
			"Data Engineering",
			"Data Warehouse vs Data Lake",
			"Khác biệt giữa Data Warehouse và Data Lake?",
			"Mid",
			"Warehouse: Dữ liệu có cấu trúc, schema-on-write, tối ưu cho phân tích (OLAP). Lake: Dữ liệu thô (cấu trúc/phi cấu trúc), schema-on-read, chi phí lưu trữ thấp.",
			"vi",
			"Data Engineer",
			"Có cấu trúc vs Phi cấu trúc. Schema-on-write vs Schema-on-read.",
		},

		// --- Java (Fresher/Junior) ---
		{
			"Java",
			"OOP Principles",
			"Name the 4 main principles of OOP.",
			"Fresher",
			"Encapsulation, Abstraction, Inheritance, Polymorphism.",
			"en",
			"BackEnd",
			"EAIP: Encapsulation, Abstraction, Inheritance, Polymorphism.",
		},
		{
			"Java",
			"Nguyên lý OOP",
			"Kể tên 4 nguyên lý chính của OOP.",
			"Fresher",
			"Đóng gói (Encapsulation), Trừu tượng (Abstraction), Kế thừa (Inheritance), Đa hình (Polymorphism).",
			"vi",
			"BackEnd",
			"EAIP: Đóng gói, Trừu tượng, Kế thừa, Đa hình.",
		},
		{
			"Java",
			"Interface vs Abstract Class",
			"Difference between Interface and Abstract Class in Java?",
			"Junior",
			"Interface: Multiple inheritance, only method signatures (until Java 8 default methods). Abstract Class: Single inheritance, can have state and implemented methods.",
			"en",
			"BackEnd",
			"Multiple Inheritance vs Single Inheritance. State vs No State.",
		},
		{
			"Java",
			"Interface vs Abstract Class",
			"Khác biệt giữa Interface và Abstract Class trong Java?",
			"Junior",
			"Interface: Đa kế thừa, chỉ có chữ ký hàm (trừ default methods từ Java 8). Abstract Class: Đơn kế thừa, có thể có biến (state) và hàm đã implement.",
			"vi",
			"BackEnd",
			"Đa kế thừa vs Đơn kế thừa. Có trạng thái vs Không trạng thái.",
		},

		{
			"Behavioral",
			"Incident Handling",
			"Describe the last incident you caused. How did you handle it?",
			"Mid",
			"Look for ownership (admitting mistake), immediate mitigation (rollback/fix), communication (alerting stakeholders), and prevention (post-mortem, fixing root cause).",
			"en",
			"DevOps",
			"Use the STAR method (Situation, Task, Action, Result). Focus on Ownership and Prevention.",
		},
		{
			"Behavioral",
			"Xử lý sự cố",
			"Lần gần nhất anh/chị gây ra incident là gì?",
			"Mid",
			"Tìm kiếm tinh thần làm chủ (thừa nhận lỗi), khắc phục ngay lập tức (rollback/fix), giao tiếp (thông báo cho bên liên quan), và phòng ngừa (post-mortem, sửa nguyên nhân gốc rễ).",
			"vi",
			"DevOps",
			"Sử dụng phương pháp STAR (Tình huống, Nhiệm vụ, Hành động, Kết quả). Tập trung vào Tinh thần làm chủ và Phòng ngừa.",
		},
		// --- Behavioral (Mid) ---
		{
			"Behavioral",
			"Incident Handling",
			"Describe the last incident you caused. How did you handle it?",
			"Mid",
			"Look for ownership (admitting mistake), immediate mitigation (rollback/fix), communication (alerting stakeholders), and prevention (post-mortem, fixing root cause).",
			"en",
			"BackEnd",
			"Use the STAR method (Situation, Task, Action, Result). Focus on Ownership and Prevention.",
		},
		{
			"Behavioral",
			"Xử lý sự cố",
			"Lần gần nhất anh/chị gây ra incident là gì?",
			"Mid",
			"Tìm kiếm tinh thần làm chủ (thừa nhận lỗi), khắc phục ngay lập tức (rollback/fix), giao tiếp (thông báo cho bên liên quan), và phòng ngừa (post-mortem, sửa nguyên nhân gốc rễ).",
			"vi",
			"BackEnd",
			"Sử dụng phương pháp STAR (Tình huống, Nhiệm vụ, Hành động, Kết quả). Tập trung vào Tinh thần làm chủ và Phòng ngừa.",
		},
		{
			"Behavioral",
			"Conflict Resolution",
			"Tell me about a time you disagreed with a team member. How did you resolve it?",
			"Mid",
			"Look for professional communication, focus on the problem not the person, seeking compromise or better solution, and maintaining relationship.",
			"en",
			"FrontEnd",
			"STAR method again. Focus on 'We' vs 'I'. Show empathy and logical reasoning.",
		},
		{
			"Behavioral",
			"Giải quyết mâu thuẫn",
			"Kể về một lần bạn bất đồng ý kiến với đồng nghiệp. Bạn giải quyết nó thế nào?",
			"Mid",
			"Tìm kiếm giao tiếp chuyên nghiệp, tập trung vào vấn đề không phải con người, tìm kiếm thỏa hiệp hoặc giải pháp tốt hơn, và duy trì mối quan hệ.",
			"vi",
			"FrontEnd",
			"Phương pháp STAR. Tập trung vào 'Chúng tôi' thay vì 'Tôi'. Thể hiện sự thấu cảm và lý luận logic.",
		},

		// --- Frontend Vue (Junior/Mid) ---
		{
			"Vue",
			"Vue Lifecycle",
			"Describe the Vue.js component lifecycle.",
			"Junior",
			"Creation (beforeCreate, created), Mounting (beforeMount, mounted), Updating (beforeUpdate, updated), Destruction (beforeDestroy, destroyed).",
			"en",
			"FrontEnd",
			"Think about the 4 phases: Create -> Mount -> Update -> Destroy.",
		},
		{
			"Vue",
			"Vòng đời Vue",
			"Mô tả vòng đời (lifecycle) của component Vue.js.",
			"Junior",
			"Khởi tạo (beforeCreate, created), Gắn kết (beforeMount, mounted), Cập nhật (beforeUpdate, updated), Hủy (beforeDestroy, destroyed).",
			"vi",
			"FrontEnd",
			"Hãy nghĩ về 4 giai đoạn: Khởi tạo -> Gắn kết -> Cập nhật -> Hủy.",
		},
		{
			"Vue",
			"Computed vs Watch",
			"Difference between Computed properties and Watchers?",
			"Junior",
			"Computed: Cached based on dependencies, synchronous, for derived state. Watch: Runs side effects on change, async supported, no caching.",
			"en",
			"FrontEnd",
			"Caching vs Side Effects. Synchronous vs Asynchronous.",
		},
		{
			"Vue",
			"Computed vs Watch",
			"Khác biệt giữa Computed properties và Watchers?",
			"Junior",
			"Computed: Được cache dựa trên dependencies, đồng bộ, dùng cho dữ liệu dẫn xuất. Watch: Chạy side effects khi thay đổi, hỗ trợ bất đồng bộ, không cache.",
			"vi",
			"FrontEnd",
			"Caching vs Side Effects (Hiệu ứng phụ). Đồng bộ vs Bất đồng bộ.",
		},

		// --- Frontend React (Advanced) ---
		{
			"React",
			"Context API vs Redux",
			"When to use Context API vs Redux?",
			"Mid",
			"Context: Low frequency updates, global themes/user data, simple state. Redux: High frequency updates, complex state logic, middleware needs, devtools debugging.",
			"en",
			"FrontEnd",
			"Complexity and Frequency. Simple global data vs Complex state management.",
		},
		{
			"React",
			"Context API vs Redux",
			"Khi nào dùng Context API so với Redux?",
			"Mid",
			"Context: Cập nhật ít thường xuyên, theme/user data, state đơn giản. Redux: Cập nhật liên tục, logic phức tạp, cần middleware, debug tool mạnh.",
			"vi",
			"FrontEnd",
			"Độ phức tạp và Tần suất. Dữ liệu toàn cục đơn giản vs Quản lý state phức tạp.",
		},

		// --- Backend Python (Junior/Mid) ---
		{
			"Python",
			"Decorators",
			"What is a Python decorator?",
			"Junior",
			"A design pattern that allows a user to add new functionality to an existing object without modifying its structure. Syntax @decorator_name.",
			"en",
			"BackEnd",
			"It wraps a function to extend its behavior. Syntax: @name.",
		},
		{
			"Python",
			"Decorators",
			"Decorator trong Python là gì?",
			"Junior",
			"Một mẫu thiết kế cho phép thêm chức năng mới vào đối tượng hiện có mà không thay đổi cấu trúc của nó. Cú pháp @decorator_name.",
			"vi",
			"BackEnd",
			"Nó bao bọc một hàm để mở rộng hành vi của nó. Cú pháp: @name.",
		},
		{
			"Python",
			"GIL",
			"What is the Global Interpreter Lock (GIL)?",
			"Senior",
			"A mutex that allows only one thread to hold the control of the Python interpreter, limiting multi-threaded performance in CPU-bound tasks.",
			"en",
			"BackEnd",
			"It prevents multiple native threads from executing Python bytecodes at once. Limits CPU-bound concurrency.",
		},
		{
			"Python",
			"GIL",
			"Global Interpreter Lock (GIL) là gì?",
			"Senior",
			"Một mutex chỉ cho phép một luồng nắm giữ quyền kiểm soát trình thông dịch Python, làm hạn chế hiệu năng đa luồng trong các tác vụ CPU-bound.",
			"vi",
			"BackEnd",
			"Nó ngăn nhiều luồng thực thi mã byte Python cùng lúc. Hạn chế đa luồng cho tác vụ CPU-bound.",
		},

		{
			"System Design",
			"Load Balancer",
			"How does a Load Balancer work? Algorithms?",
			"Any",
			"Distributes traffic across servers. Algorithms: Round Robin, Least Connections, IP Hash.",
			"en",
			"Any",
			"Traffic Distribution. Algorithms: Round Robin, Least Connections.",
		},
		{
			"System Design",
			"Load Balancer",
			"Load Balancer hoạt động thế nào? Các thuật toán?",
			"Any",
			"Phân phối lưu lượng truy cập qua các server. Thuật toán: Round Robin, Least Connections, IP Hash.",
			"vi",
			"Any",
			"Phân phối lưu lượng. Thuật toán: Round Robin, Least Connections.",
		},
		// --- Backend General (System Design/Microservices) ---
		{
			"System Design",
			"Microservices vs Monolith",
			"Pros and Cons of Microservices?",
			"Any",
			"Pros: Independent scaling, technology agnostic, fault isolation. Cons: Complexity, network latency, data consistency (distributed transactions).",
			"en",
			"Any",
			"Scale vs Complexity. Independence vs Consistency.",
		},
		{
			"System Design",
			"Microservices vs Monolith",
			"Ưu nhược điểm của Microservices?",
			"Any",
			"Ưu: Scale độc lập, đa công nghệ, cô lập lỗi. Nhược: Phức tạp, độ trễ mạng, tính nhất quán dữ liệu (giao dịch phân tán).",
			"vi",
			"Any",
			"Scale vs Phức tạp. Độc lập vs Nhất quán.",
		},

		// --- Data Engineering (Advanced) ---
		{
			"Data Engineering",
			"Window Functions",
			"What are SQL Window Functions?",
			"Mid",
			"Perform calculations across a set of table rows that are somehow related to the current row (e.g., RANK, LEAD, LAG, ROW_NUMBER) without collapsing rows like GROUP BY.",
			"en",
			"Data Engineer",
			"Calculations across related rows without grouping. Keywords: OVER(), PARTITION BY.",
		},
		{
			"Data Engineering",
			"Window Functions",
			"SQL Window Functions là gì?",
			"Mid",
			"Thực hiện tính toán trên một tập hợp các hàng liên quan đến hàng hiện tại (ví dụ: RANK, LEAD, LAG, ROW_NUMBER) mà không gộp hàng như GROUP BY.",
			"vi",
			"Data Engineer",
			"Tính toán trên các hàng liên quan mà không gộp nhóm. Từ khóa: OVER(), PARTITION BY.",
		},
		{
			"Data Engineering",
			"Spark RDD vs DataFrame",
			"Difference between RDD and DataFrame in Spark?",
			"Mid",
			"RDD: Low-level, type-safe, slower (no optimization). DataFrame: High-level, schema-aware, optimized (Catalyst optimizer).",
			"en",
			"Data Engineer",
			"Low-level vs High-level. Optimization (Catalyst).",
		},
		{
			"Data Engineering",
			"Spark RDD vs DataFrame",
			"Khác biệt giữa RDD và DataFrame trong Spark?",
			"Mid",
			"RDD: Mức thấp, an toàn kiểu, chậm hơn (không tối ưu). DataFrame: Mức cao, có schema, được tối ưu hóa (Catalyst optimizer).",
			"vi",
			"Data Engineer",
			"Mức thấp vs Mức cao. Tối ưu hóa (Catalyst).",
		},

		// --- DevOps (Kubernetes) ---
		{
			"Kubernetes",
			"Kubernetes Pod",
			"What is a Pod in Kubernetes?",
			"Any",
			"The smallest deployable unit in K8s. Represents a single instance of a running process, can contain one or more containers sharing storage/network.",
			"en",
			"DevOps",
			"Smallest unit in K8s. Can hold one or more containers.",
		},
		{
			"Kubernetes",
			"Kubernetes Pod",
			"Pod trong Kubernetes là gì?",
			"Any",
			"Đơn vị triển khai nhỏ nhất trong K8s. Đại diện cho một instance của process đang chạy, có thể chứa một hoặc nhiều container chia sẻ storage/network.",
			"vi",
			"DevOps",
			"Đơn vị nhỏ nhất trong K8s. Có thể chứa một hoặc nhiều container.",
		},

		// --- FullStack (Added for completeness) ---
		{
			"React",
			"React Hooks",
			"What are React Hooks? Name common ones.",
			"Fresher",
			"Hooks allow using state and other React features in functional components. Common ones: useState, useEffect, useContext, useRef.",
			"en",
			"FullStack",
			"They let you use state and other features without writing a class. Think about managing data (state) and side effects.",
		},
		{
			"React",
			"React Hooks",
			"React Hooks là gì? Kể tên vài hook phổ biến.",
			"Fresher",
			"Hooks cho phép dùng state và các tính năng React trong functional components. Phổ biến: useState, useEffect, useContext, useRef.",
			"vi",
			"FullStack",
			"Chúng cho phép bạn sử dụng state và các tính năng khác mà không cần viết class. Hãy nghĩ về quản lý dữ liệu (state) và hiệu ứng phụ (side effects).",
		},
		{
			"NodeJS",
			"Event Loop",
			"Explain the Node.js Event Loop.",
			"Junior",
			"Single-threaded loop that handles asynchronous callbacks. Phases: Timers, Pending Callbacks, Poll, Check (setImmediate), Close Callbacks. Offloads heavy tasks to libuv worker pool.",
			"en",
			"FullStack",
			"It's the mechanism that allows Node.js to perform non-blocking I/O operations. Mention the phases and libuv.",
		},
		{
			"NodeJS",
			"Event Loop",
			"Giải thích Node.js Event Loop.",
			"Junior",
			"Vòng lặp đơn luồng xử lý các callback bất đồng bộ. Các giai đoạn: Timers, Pending Callbacks, Poll, Check (setImmediate), Close Callbacks. Tác vụ nặng được đẩy xuống libuv worker pool.",
			"vi",
			"FullStack",
			"Đó là cơ chế cho phép Node.js thực hiện các thao tác I/O không chặn. Hãy đề cập đến các giai đoạn và libuv.",
		},
		{
			"System Design",
			"Rate Limiting",
			"How would you limit 10 concurrent uploads per user?",
			"Senior",
			"Use a distributed counter (e.g., Redis) or a Semaphore pattern. When upload starts, increment counter/acquire lock. If > 10, reject. When finished, decrement/release. Handle race conditions and timeouts.",
			"en",
			"FullStack",
			"You need a shared state store to track active uploads across servers. Redis is a common choice.",
		},
		{
			"System Design",
			"Rate Limiting",
			"Làm sao limit 10 concurrent upload / user?",
			"Senior",
			"Dùng bộ đếm phân tán (ví dụ: Redis) hoặc mẫu Semaphore. Khi bắt đầu upload, tăng biến đếm. Nếu > 10, từ chối. Khi xong, giảm biến đếm. Cần xử lý race condition và timeout.",
			"vi",
			"FullStack",
			"Bạn cần một kho lưu trữ trạng thái chia sẻ để theo dõi các upload đang hoạt động trên các server. Redis là lựa chọn phổ biến.",
		},

		// --- Additional Questions (Crawler Sync) ---
		{
			"Data Engineering",
			"Apache Kafka Basics",
			"Apache Kafka là gì? Các thành phần chính?",
			"Mid",
			"Kafka là nền tảng phân tán xử lý luồng sự kiện. Thành phần: Producer, Consumer, Broker, Topic, Partition, Offset. Hiệu năng cao nhờ Sequential I/O và Zero Copy.",
			"vi",
			"Data Engineer",
			"Distributed Event Streaming. Sequential I/O. Zero Copy.",
		},
		{
			"Data Engineering",
			"Airflow & Workflow Orchestration",
			"Apache Airflow là gì? Khái niệm DAG?",
			"Senior",
			"Airflow là platform lập lịch workflow. DAG (Directed Acyclic Graph) biểu diễn luồng công việc. Backfill dùng để chạy lại task quá khứ.",
			"vi",
			"Data Engineer",
			"Python-based. Scheduler. Operators.",
		},
		{
			"Database",
			"Redis Use Cases",
			"Khi nào nên dùng Redis?",
			"Mid",
			"Dùng làm Cache (tăng tốc đọc), Session Storage, Message Broker (Pub/Sub), Leaderboard (Sorted Set), Rate Limiting.",
			"vi",
			"BackEnd",
			"In-memory. Key-Value. Pub/Sub.",
		},
		{
			"Frontend Basic",
			"Frontend Performance",
			"Các kỹ thuật tối ưu performance cho web app?",
			"Senior",
			"Lazy Loading, Code Splitting, Caching (HTTP cache, Service Worker), Tối ưu ảnh (WebP, lazy load), Tối ưu Critical Rendering Path (inline CSS critical, defer JS).",
			"vi",
			"FrontEnd",
			"Giảm size bundle. Giảm request. Tối ưu render blocking.",
		},
		{
			"Terraform",
			"Infrastructure as Code (Terraform)",
			"Tại sao nên dùng Terraform (IaC)?",
			"Mid",
			"Nhất quán (tránh lỗi cấu hình tay), Version Control (lưu code trong Git), Tái sử dụng (Module), Tài liệu hóa hạ tầng.",
			"vi",
			"DevOps",
			"Consistency, Version Control. State file.",
		},

		// --- DevOps Additional (Senior/Behavioral) ---
		{
			"CI/CD",
			"Blue Green vs Canary",
			"Explain Blue/Green Deployment vs Canary Release.",
			"Senior",
			"Blue/Green: Two identical environments, switch traffic instantly (fast rollback, expensive). Canary: Roll out to small % of users first, then expand (safer, slower).",
			"en",
			"DevOps",
			"Traffic switching strategies. Cost vs Safety.",
		},
		{
			"CI/CD",
			"Blue Green vs Canary",
			"Giải thích Blue/Green Deployment và Canary Release.",
			"Senior",
			"Blue/Green: Hai môi trường giống hệt nhau, chuyển traffic ngay lập tức (rollback nhanh, tốn kém). Canary: Triển khai cho một phần nhỏ user trước, rồi mở rộng dần (an toàn hơn, chậm hơn).",
			"vi",
			"DevOps",
			"Chiến lược chuyển đổi traffic. Chi phí so với An toàn.",
		},
		{
			"System Design",
			"High Availability DevOps",
			"Design a HA architecture for a Kubernetes cluster.",
			"Senior",
			"Multi-master control plane (3+ nodes), etcd across zones, worker nodes across Availability Zones (AZs), Load Balancer distribution, Auto-scaling groups.",
			"en",
			"DevOps",
			"Avoid Single Point of Failure (SPOF). Multi-AZ.",
		},
		{
			"System Design",
			"High Availability DevOps",
			"Thiết kế kiến trúc HA cho Kubernetes cluster.",
			"Senior",
			"Control plane đa master (3+ node), etcd phân tán qua các zone, worker nodes rải rác qua các Availability Zones (AZ), Load Balancer phân phối, Auto-scaling groups.",
			"vi",
			"DevOps",
			"Tránh điểm chết duy nhất (SPOF). Đa vùng (Multi-AZ).",
		},
		{
			"Behavioral",
			"Learning New Tech",
			"How do you keep up with new DevOps tools?",
			"Junior",
			"Follow blogs (Hacker News, Medium), official docs, POC projects, conferences (KubeCon), team sharing.",
			"en",
			"DevOps",
			"Continuous learning strategies. Hands-on practice.",
		},
		{
			"Behavioral",
			"Learning New Tech",
			"Bạn cập nhật công nghệ DevOps mới thế nào?",
			"Junior",
			"Theo dõi blog, tài liệu chính thức, làm dự án POC, tham gia hội thảo (KubeCon), chia sẻ trong team.",
			"vi",
			"DevOps",
			"Chiến lược học tập liên tục. Thực hành thực tế.",
		},

		// --- DevOps (Docker) ---
		{
			"Docker",
			"Docker vs VM",
			"Difference between Docker Containers and Virtual Machines?",
			"Any",
			"Docker: Shares OS kernel, lightweight, faster startup. VM: Has own OS, heavier, stronger isolation.",
			"en",
			"DevOps",
			"Shared Kernel vs Full OS. Lightweight vs Isolation.",
		},
		{
			"Docker",
			"Docker vs VM",
			"Khác biệt giữa Docker Container và Virtual Machine?",
			"Any",
			"Docker: Chia sẻ kernel OS, nhẹ, khởi động nhanh. VM: Có OS riêng, nặng hơn, cô lập tốt hơn.",
			"vi",
			"DevOps",
			"Chia sẻ Kernel vs Full OS. Nhẹ vs Cô lập.",
		},

		// --- DevOps (Terraform) ---
		{
			"Terraform",
			"Infrastructure as Code (Terraform)",
			"Why use Terraform (IaC)?",
			"Any",
			"Consistency (avoid manual errors), Version Control (Git), Reusability (Modules), Documentation of infrastructure.",
			"en",
			"DevOps",
			"Consistency, Version Control. State file.",
		},
		{
			"Terraform",
			"Infrastructure as Code (Terraform)",
			"Tại sao dùng Terraform (IaC)?",
			"Any",
			"Tính nhất quán (tránh lỗi thủ công), Quản lý phiên bản (Git), Tái sử dụng (Modules), Tài liệu hóa hạ tầng.",
			"vi",
			"DevOps",
			"Nhất quán. Version Control. State file.",
		},
		{
			"Terraform",
			"Terraform State",
			"What is Terraform State?",
			"Any",
			"It maps real world resources to your configuration, keeps track of metadata, and improves performance for large infrastructures. Stored locally or remotely (S3).",
			"en",
			"DevOps",
			"Mapping configuration to reality. Metadata. Locking.",
		},
		{
			"Terraform",
			"Terraform State",
			"Terraform State là gì?",
			"Any",
			"Nó ánh xạ tài nguyên thực tế với cấu hình, theo dõi metadata và cải thiện hiệu năng. Lưu trữ local hoặc remote (S3).",
			"vi",
			"DevOps",
			"Ánh xạ cấu hình với thực tế. Metadata. Locking.",
		},

		// --- DevOps (CI/CD) ---
		{
			"CI/CD",
			"CI/CD Pipeline Stages",
			"Common stages in a CI/CD pipeline?",
			"Any",
			"Build (Compile), Test (Unit/Integration), Security Scan, Artifact (Docker Image), Deploy (Staging/Prod).",
			"en",
			"DevOps",
			"Build -> Test -> Scan -> Artifact -> Deploy.",
		},
		{
			"CI/CD",
			"CI/CD Pipeline Stages",
			"Các giai đoạn phổ biến trong CI/CD pipeline?",
			"Any",
			"Build (Biên dịch), Test (Unit/Integration), Security Scan, Artifact (Docker Image), Deploy (Staging/Prod).",
			"vi",
			"DevOps",
			"Build -> Test -> Scan -> Artifact -> Deploy.",
		},

		// --- DevOps (CV Screening) ---
		{
			"CV Screening",
			"DevOps Challenges",
			"What was the biggest infrastructure challenge you faced?",
			"Any",
			"Look for scaling issues, downtime recovery, security breaches, or migration challenges (e.g., On-prem to Cloud).",
			"en",
			"DevOps",
			"Situation -> Task -> Action -> Result (STAR). Focus on problem solving.",
		},
		{
			"CV Screening",
			"DevOps Challenges",
			"Thách thức hạ tầng lớn nhất bạn từng gặp?",
			"Any",
			"Tìm kiếm các vấn đề về mở rộng (scaling), khôi phục sau sự cố (downtime), bảo mật, hoặc di chuyển hệ thống (migration).",
			"vi",
			"DevOps",
			"Situation -> Task -> Action -> Result (STAR). Tập trung vào giải quyết vấn đề.",
		},

		// --- DevOps (Behavioral) ---
		{
			"Behavioral",
			"Production Incident",
			"Describe a production incident you handled.",
			"Mid",
			"Focus on: Detection (Monitoring), Response (Incident Management), Resolution (Fix), and Post-Mortem (Prevention).",
			"en",
			"DevOps",
			"Detection -> Response -> Resolution -> Post-Mortem.",
		},
		{
			"Behavioral",
			"Production Incident",
			"Mô tả một sự cố production bạn đã xử lý.",
			"Mid",
			"Tập trung vào: Phát hiện (Monitoring), Phản ứng (Quản lý sự cố), Giải quyết (Fix), và Hậu kỳ (Ngăn chặn tái diễn).",
			"vi",
			"DevOps",
			"Phát hiện -> Phản ứng -> Giải quyết -> Hậu kỳ (Post-Mortem).",
		},
		// --- JavaScript (General) ---
		{
			"JavaScript",
			"Let vs Var vs Const",
			"What is the difference between let, var, and const in JavaScript?",
			"Any",
			"var: function-scoped, hoisted. let: block-scoped, not hoisted. const: block-scoped, immutable reference.",
			"en",
			"Any",
			"Scope (Function vs Block), Hoisting, Mutability.",
		},
		{
			"JavaScript",
			"Let vs Var vs Const",
			"Sự khác biệt giữa let, var và const trong JavaScript?",
			"Any",
			"var: phạm vi hàm (function-scoped), hoisting. let: phạm vi khối (block-scoped), không hoisting. const: phạm vi khối, tham chiếu không đổi.",
			"vi",
			"Any",
			"Phạm vi (Scope) và Hoisting. Tính thay đổi (Mutability).",
		},
	}

	for _, q := range questions {
		var topicID string
		err := db.QueryRow("SELECT id FROM topics WHERE name = $1", q.Topic).Scan(&topicID)
		if err != nil {
			log.Printf("Topic not found for %s: %v", q.Topic, err)
			continue
		}

		// Check if question exists by content
		var existingID string
		err = db.QueryRow("SELECT id FROM questions WHERE content = $1", q.Content).Scan(&existingID)

		if err == nil {
			// Exists: Update it (especially Level and Role)
			_, err = db.Exec(`
				UPDATE questions 
				SET title = $1, level = $2, topic_id = $3, correct_answer = $4, language = $5, role = $6, hint = $7
				WHERE id = $8
			`, q.Title, q.Level, topicID, q.CorrectAnswer, q.Language, q.Role, q.Hint, existingID)

			if err != nil {
				log.Printf("Failed to update question %s: %v", q.Title, err)
			} else {
				fmt.Printf("Updated question: %s (%s) [%s]\n", q.Title, q.Language, q.Role)
			}
		} else {
			// Not exists: Insert
			_, err = db.Exec(`
				INSERT INTO questions (id, title, content, level, topic_id, created_by, status, correct_answer, language, role, hint)
				VALUES (uuid_generate_v4(), $1, $2, $3, $4, '123e4567-e89b-12d3-a456-426614174000', 'published', $5, $6, $7, $8)
			`, q.Title, q.Content, q.Level, topicID, q.CorrectAnswer, q.Language, q.Role, q.Hint)

			if err != nil {
				log.Printf("Failed to insert question %s: %v", q.Title, err)
			} else {
				fmt.Printf("Inserted question: %s (%s) [%s]\n", q.Title, q.Language, q.Role)
			}
		}
	}
}
