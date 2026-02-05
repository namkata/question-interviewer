ALTER TABLE questions ADD COLUMN language VARCHAR(10) DEFAULT 'en';
CREATE INDEX idx_questions_language ON questions(language);
