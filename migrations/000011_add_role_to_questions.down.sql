DROP INDEX IF EXISTS idx_questions_role;
ALTER TABLE questions DROP COLUMN IF EXISTS role;
