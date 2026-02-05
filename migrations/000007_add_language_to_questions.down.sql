DROP INDEX IF EXISTS idx_questions_language;
ALTER TABLE questions DROP COLUMN IF EXISTS language;
