# Services Design

## Go Microservices

### Question Service

* Create / update / delete questions
* Assign topic, level, tags
* Publish / unpublish

### Answer Service

* Submit answers
* Vote / unvote
* Mark accepted answer

### Practice Service

* Start practice session
* Randomize questions
* Calculate score

## FastAPI Adapter Services

### Search Service

* Sync questions & answers
* Full-text search
* Ranking

### AI / Content Service

* Suggest answers
* Summarize questions
* Difficulty classification

### Background Worker

* Re-index search
* Analytics events
* Notifications

## Architecture Pattern

* Hexagonal / Clean Architecture
* Domain isolated from adapters
