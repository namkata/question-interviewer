ALTER TABLE practice_sessions ADD COLUMN status VARCHAR(50) DEFAULT 'in_progress';
UPDATE practice_sessions SET status = 'in_progress' WHERE status IS NULL;
