# Read / Write Strategy

## Write Path (Strong Consistency)

* Client → Next.js → Go Service → PostgreSQL Primary
* Used for:

  * Create / update / delete
  * Voting
  * Practice progress

Characteristics:

* ACID transactions
* Row-level locking

## Read Path (Optimized)

### Option 1: Read Replica

* Question listing
* Question detail

### Option 2: Cache / Search

* Popular questions
* Full-text search

## CQRS-lite

* Write model: normalized schema
* Read model: denormalized view

Example:

* Materialized view for question list
