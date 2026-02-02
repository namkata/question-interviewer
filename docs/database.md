# Database Design (PostgreSQL)

## Core Tables

### users

* id (UUID, PK)
* email (unique)
* username (unique)
* role
* created_at

### topics

* id (UUID, PK)
* name (unique)
* description

### questions

* id (UUID, PK)
* title
* content
* level
* topic_id (FK)
* created_by (FK)
* status
* created_at
* updated_at

Indexes:

* topic_id
* level

### answers

* id (UUID, PK)
* question_id (FK)
* content
* created_by (FK)
* vote_count
* is_accepted
* created_at

### votes

* user_id (PK)
* answer_id (PK)
* value (+1 / -1)

Constraint:

* One vote per user per answer

### bookmarks

* user_id (PK)
* question_id (PK)
* created_at

### practice_sessions

* id (UUID, PK)
* user_id (FK)
* score
* started_at
* ended_at
