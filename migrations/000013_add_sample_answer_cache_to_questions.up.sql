ALTER TABLE questions
    ADD COLUMN sample_answer TEXT,
    ADD COLUMN sample_feedback TEXT,
    ADD COLUMN sample_suggestions JSONB,
    ADD COLUMN sample_source VARCHAR(20) DEFAULT 'seed';

