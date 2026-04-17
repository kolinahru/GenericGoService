# GenericGoService
A production-style backend service in Go demonstrating:

- REST API with `net/http`
- PostgreSQL integration
- Redis caching (cache-aside)
- Clean architecture (Handler → Service → Repository)
- Concurrency with goroutines + worker pools
- Background job processing

---

# Architecture Overview
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

# Features

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

# Project Structure
GenericGoService/  
&nbsp;&nbsp;main.go  
&nbsp;&nbsp;db.go  
  
&nbsp;&nbsp;models/  
&nbsp;&nbsp;&nbsp;&nbsp;item.go  
  
&nbsp;&nbsp;repository/  
&nbsp;&nbsp;&nbsp;&nbsp;item_repository.go  
  
&nbsp;&nbsp;service/  
&nbsp;&nbsp;&nbsp;&nbsp;item_service.go  
  
&nbsp;&nbsp;handlers/  
&nbsp;&nbsp;&nbsp;&nbsp;item_handler.go  
&nbsp;&nbsp;&nbsp;&nbsp;response.go  
  
&nbsp;&nbsp;jobs/  
&nbsp;&nbsp;&nbsp;&nbsp;job.go  
&nbsp;&nbsp;&nbsp;&nbsp;queue.go  
&nbsp;&nbsp;&nbsp;&nbsp;worker_pool.go  
  
seed.sql  

---

# Database Setup

## Create table
`CREATE TABLE IF NOT EXISTS items (` 
`id SERIAL PRIMARY KEY,`  
`name TEXT NOT NULL);`

---

# Seed Data

## `seed.sql`  
`INSERT INTO items (name) VALUES`  
`('Keyboard'),`  
`('Mouse'),`  
`('Monitor'),`  
`('Laptop'),`  
`('Desk'),`  
`('Chair'),`  
`('Webcam'),`  
`('Microphone');`   

---

## Run seed
`psql -h localhost -U postgres -d goapp -f seed.sql` 
Password: postgres  

---

# Running the App: 
`go run .`

Server: http://localhost:8080

---

# API Endpoints

## GET all: 
`GET /items` 
## GET by ID: 
`GET /items/{id}` 
## POST: 
`POST /items` 
`{` 
 ` "name": "New Item"` 
`}` 
## PUT: 
`PUT /items/{id}` 
`{` 
  `"name": "Updated Item"` 
`}` 
## DELETE:
`DELETE /items/{id}` 

# Testing (PowerShell)
## Create
`Invoke-RestMethod` 
  `-Uri http://localhost:8080/items` 
  `-Method POST` 
  `-Body '{"name":"Test Item"}'` 
  `-ContentType "application/json"` 
## Get
`Invoke-RestMethod http://localhost:8080/items` 
## Update
`Invoke-RestMethod` 
  `-Uri http://localhost:8080/items/1` 
  `-Method PUT` 
  `-Body '{"name":"Updated"}'` 
  `-ContentType "application/json"` 
##D elete
`Invoke-RestMethod` 
  `-Uri http://localhost:8080/items/1` 
  `-Method DELETE` 
