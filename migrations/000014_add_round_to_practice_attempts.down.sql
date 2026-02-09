DROP INDEX IF EXISTS idx_practice_attempts_round_name;

ALTER TABLE practice_attempts
    DROP COLUMN IF EXISTS round_index,
    DROP COLUMN IF EXISTS round_name;

