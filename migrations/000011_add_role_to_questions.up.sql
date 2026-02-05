ALTER TABLE questions ADD COLUMN role VARCHAR(100);
CREATE INDEX idx_questions_role ON questions(role);
