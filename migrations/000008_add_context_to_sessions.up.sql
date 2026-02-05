ALTER TABLE practice_sessions ADD COLUMN topic_id UUID REFERENCES topics(id);
ALTER TABLE practice_sessions ADD COLUMN level VARCHAR(50);
ALTER TABLE practice_sessions ADD COLUMN language VARCHAR(10) DEFAULT 'en';
