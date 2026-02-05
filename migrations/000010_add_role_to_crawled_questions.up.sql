ALTER TABLE crawled_questions ADD COLUMN detected_role VARCHAR(100);
CREATE INDEX idx_crawled_questions_role ON crawled_questions(detected_role);
