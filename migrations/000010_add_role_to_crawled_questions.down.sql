DROP INDEX IF EXISTS idx_crawled_questions_role;
ALTER TABLE crawled_questions DROP COLUMN IF EXISTS detected_role;
