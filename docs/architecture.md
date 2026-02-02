# Architecture Overview

## Purpose

High-level architecture for the Interview Q&A Platform using Go microservices, FastAPI adapters, and Next.js BFF.

## Components

### Frontend / BFF (Next.js)

* Server-side rendering for SEO
* Authentication & session handling
* Aggregates APIs from backend services

### Core Backend (Go Microservices)

* Question Service
* Answer Service
* Practice Service

Responsibilities:

* Business logic
* Validation
* Transactional writes

### Adapter Services (FastAPI)

* Search Service
* AI / Content Service
* Background Worker

Responsibilities:

* Full-text search
* Async processing
* External integrations

## Communication

* REST/JSON over HTTP
* Event-based (optional) via message broker

## Data Source

* PostgreSQL as source of truth
* Redis / Search engine as read optimization
