ALTER TABLE questions ADD COLUMN correct_answer TEXT;
ALTER TABLE questions ALTER COLUMN title DROP NOT NULL;
