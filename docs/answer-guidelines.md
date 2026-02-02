# Answer Guideline & Suggested Answers Strategy

## 1. Mục tiêu

Nhiều người đọc **không biết cách trả lời phỏng vấn như thế nào cho đúng kỳ vọng interviewer**. Vì vậy hệ thống **không chỉ hiển thị câu trả lời**, mà còn **gợi ý cách trả lời chuẩn**, có cấu trúc, dễ học và dễ áp dụng khi đi phỏng vấn thật.

---

## 2. Nguyên tắc thiết kế Answer

### Không chỉ là "đáp án đúng"

* Tránh answer kiểu sách giáo khoa
* Tập trung vào **cách diễn đạt khi phỏng vấn**
* Có context + ví dụ + trade-off

### Phân cấp câu trả lời

* Junior: hiểu khái niệm
* Mid: biết áp dụng
* Senior: phân tích trade-off, edge case

---

## 3. Cấu trúc Answer Chuẩn (Answer Template)

Mỗi câu trả lời nên theo format sau:

### 1️⃣ Short Answer (TL;DR)

> Trả lời ngắn gọn trong 2–3 câu, đúng kiểu nói miệng khi phỏng vấn.

### 2️⃣ Detailed Explanation

* Giải thích khái niệm
* Nguyên lý hoạt động
* Khi nào nên dùng / không nên dùng

### 3️⃣ Example (Rất quan trọng)

* Ví dụ code (nếu có)
* Ví dụ hệ thống thực tế

### 4️⃣ Trade-offs / Pitfalls

* Ưu điểm
* Nhược điểm
* Sai lầm thường gặp

### 5️⃣ Follow-up Questions

* Intervier có thể hỏi gì tiếp theo?
* Gợi ý cách trả lời tiếp

---

## 4. Answer Types trong hệ thống

### Canonical Answer

* Do maintainer / admin viết
* Được review kỹ
* Dùng làm chuẩn học

### Community Answer

* Do user đóng góp
* Có vote
* Có thể bổ sung góc nhìn khác

### Suggested Answer (Gợi ý)

* Sinh ra tự động (AI hoặc rule-based)
* Dùng cho người mới
* Luôn có disclaimer

---

## 5. Suggested Answer – Chiến lược triển khai

### 5.1 Rule-based (Phase 1)

* Template cố định
* Fill nội dung theo topic
* An toàn, dễ kiểm soát

### 5.2 AI-assisted (Phase 2)

* Generate answer theo level
* Rewrite cho interview-style
* Sinh follow-up questions

⚠️ Lưu ý:

* Không auto-publish
* Chỉ dùng làm gợi ý

---

## 6. Database Design (Answer Extension)

### answers (extend)

* id
* question_id
* content
* answer_type (canonical | community | suggested)
* level_target (junior | mid | senior)
* created_by (nullable)
* vote_count

---

## 7. UI / UX Gợi ý

### Khi người dùng mở câu hỏi

Tabs:

* Suggested Answer (default)
* Canonical Answer
* Community Answers

Badge:

* "Interview-style"
* "Beginner-friendly"

---

## 8. Ví dụ minh hoạ

### Question

> What is a goroutine?

### Suggested Answer (Junior)

**Short Answer:**
A goroutine is a lightweight thread managed by the Go runtime that allows concurrent execution of functions.

**Explanation:**
Goroutines are cheaper than OS threads and are s
