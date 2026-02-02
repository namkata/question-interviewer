ALTER TABLE questions DROP COLUMN correct_answer;
ALTER TABLE questions ALTER COLUMN title SET NOT NULL;
