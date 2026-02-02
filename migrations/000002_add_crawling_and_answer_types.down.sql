DROP INDEX IF EXISTS idx_answers_answer_type;
DROP INDEX IF EXISTS idx_crawled_questions_status;

ALTER TABLE answers ALTER COLUMN created_by SET NOT NULL;
ALTER TABLE answers DROP COLUMN level_target;
ALTER TABLE answers DROP COLUMN answer_type;

DROP TABLE IF EXISTS crawled_questions;
