INSERT INTO users (id, email, username, password_hash, role)
VALUES ('123e4567-e89b-12d3-a456-426614174000', 'demo@example.com', 'demo_user', 'hashed_password', 'user')
ON CONFLICT (id) DO NOTHING;

INSERT INTO topics (id, name, description)
VALUES 
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Golang', 'Go programming language questions'),
    ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'System Design', 'System design and architecture questions')
ON CONFLICT (name) DO NOTHING;

INSERT INTO questions (id, title, content, level, topic_id, created_by, status)
VALUES
    (uuid_generate_v4(), 'Goroutines vs Threads', 'Explain the difference between Goroutines and OS threads.', 'Junior', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '123e4567-e89b-12d3-a456-426614174000', 'published'),
    (uuid_generate_v4(), 'Channels', 'What are channels in Go and how are they used?', 'Mid', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '123e4567-e89b-12d3-a456-426614174000', 'published'),
    (uuid_generate_v4(), 'CAP Theorem', 'Explain CAP theorem in distributed systems.', 'Senior', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', '123e4567-e89b-12d3-a456-426614174000', 'published')
ON CONFLICT DO NOTHING;
