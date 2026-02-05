ALTER TABLE questions
    DROP COLUMN IF EXISTS sample_answer,
    DROP COLUMN IF EXISTS sample_feedback,
    DROP COLUMN IF EXISTS sample_suggestions,
    DROP COLUMN IF EXISTS sample_source;

