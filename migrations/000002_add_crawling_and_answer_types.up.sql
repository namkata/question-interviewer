-- Add crawled_questions table for Staging Layer
CREATE TABLE crawled_questions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source VARCHAR(255) NOT NULL,
    raw_title TEXT NOT NULL,
    raw_content TEXT,
    url VARCHAR(500),
    detected_topic VARCHAR(255),
    detected_level VARCHAR(50),
    status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Extend answers table for Answer Strategy
ALTER TABLE answers ADD COLUMN answer_type VARCHAR(50) DEFAULT 'community' CHECK (answer_type IN ('canonical', 'community', 'suggested'));
ALTER TABLE answers ADD COLUMN level_target VARCHAR(50) DEFAULT 'mid' CHECK (level_target IN ('junior', 'mid', 'senior'));

-- Make created_by nullable for auto-generated answers
ALTER TABLE answers ALTER COLUMN created_by DROP NOT NULL;

-- Create indexes for new columns
CREATE INDEX idx_crawled_questions_status ON crawled_questions(status);
CREATE INDEX idx_answers_answer_type ON answers(answer_type);
