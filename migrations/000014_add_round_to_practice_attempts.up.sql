ALTER TABLE practice_attempts
    ADD COLUMN round_index INT,
    ADD COLUMN round_name TEXT;

CREATE INDEX IF NOT EXISTS idx_practice_attempts_round_name ON practice_attempts(round_name);

