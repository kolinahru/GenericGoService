# GenericGoService
A production-style backend service in Go demonstrating:

- REST API with `net/http`
- PostgreSQL integration
- Redis caching (cache-aside)
- Clean architecture (Handler → Service → Repository)
- Concurrency with goroutines + worker pools
- Background job processing

---

# 🧭 Architecture Overview
Client
↓
HTTP Handler
↓
Service Layer (business logic, caching, jobs)
↓
Repository (PostgreSQL)

Service also interacts with:
→ Redis (cache)
→ Worker Pool (async jobs)

---

# 🚀 Features

- Full CRUD API
- PostgreSQL as source of truth
- Redis cache:
  - `items:all`
  - `items:{id}`
- Cache invalidation on writes
- Background jobs using:
  - channels
  - worker pool
- Clean, testable architecture

---

# 📁 Project Structure
go-day3/
main.go
db.go

models/
item.go

repository/
item_repository.go

service/
item_service.go

handlers/
item_handler.go
response.go

jobs/
job.go
queue.go
worker_pool.go

seed.sql
docker-compose.yml


---

# ⚙️ Configuration (Environment Variables)

Create a `.env` file (optional):
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=goapp

REDIS_ADDR=localhost:6379
PORT=8080
---

# 🐳 Docker Compose (Recommended)

## `docker-compose.yml`
version: '3.9'

services:
postgres:
image: postgres:15
container_name: go-postgres
environment:
POSTGRES_USER: postgres
POSTGRES_PASSWORD: postgres
POSTGRES_DB: goapp
ports:
- "5432:5432"
volumes:
- pgdata:/var/lib/postgresql/data

redis:
image: redis:7
container_name: go-redis
ports:
- "6379:6379"

volumes:
pgdata:
---

## Start services
docker-compose up -d

---

# 🗄️ Database Setup

## Create table
CREATE TABLE IF NOT EXISTS items (
id SERIAL PRIMARY KEY,
name TEXT NOT NULL
);

---

# 🌱 Seed Data

## `seed.sql`
INSERT INTO items (name) VALUES
('Keyboard'),
('Mouse'),
('Monitor'),
('Laptop'),
('Desk'),
('Chair'),
('Webcam'),
('Microphone');

---

## Run seed
psql -h localhost -U postgres -d goapp -f seed.sql

Password: postgres


---

# ▶️ Running the App: go run .

Server:


http://localhost:8080


---

# 📡 API Endpoints

## GET all


GET /items


---

## GET by ID


GET /items/{id}


---

## POST


POST /items


```json
{
  "name": "New Item"
}
PUT
PUT /items/{id}
{
  "name": "Updated Item"
}
DELETE
DELETE /items/{id}
🧪 Testing (PowerShell)
Create
Invoke-RestMethod `
  -Uri http://localhost:8080/items `
  -Method POST `
  -Body '{"name":"Test Item"}' `
  -ContentType "application/json"
Get
Invoke-RestMethod http://localhost:8080/items
Update
Invoke-RestMethod `
  -Uri http://localhost:8080/items/1 `
  -Method PUT `
  -Body '{"name":"Updated"}' `
  -ContentType "application/json"
Delete
Invoke-RestMethod `
  -Uri http://localhost:8080/items/1 `
  -Method DELETE
⚡ Caching Strategy
Pattern: Cache-Aside
Read Flow
Check Redis
Cache miss → query PostgreSQL
Store result in Redis
Write Flow
Invalidate items:all
Update or delete items:{id}
🔁 Background Jobs
In-memory queue using Go channels
Fixed worker pool (3 workers)
Triggered on:
create
update
delete

Example jobs:

reindex-item
delete-from-index
🧵 Concurrency Model
Goroutines for workers
Channels for job queue
WaitGroup for coordination
⚠️ Limitations (Intentional)
In-memory job queue (not durable)
No retries / dead-letter queue
No authentication
No distributed scaling
🚀 Production Improvements
Infrastructure
Move jobs to Google Cloud Pub/Sub
Add Kubernetes deployment (GKE)
Observability
Structured logging (Zap / Logrus)
Metrics (Prometheus)
Tracing (OpenTelemetry)
Reliability
Retry logic + dead-letter queue
Circuit breakers
Graceful shutdown
API
Replace manual routing with router (Chi / Gorilla)
Add validation middleware

