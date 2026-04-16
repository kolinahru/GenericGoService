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
