# Content Service

Short description
- Provides persistent storage and business logic for lesson content. Implements data access and content-related operations.

Build & run (development)
- From repository root: `go run ./service/content`
- Or build: `go build -o bin/content ./service/content` then `./bin/content`

Docker
- Build: `docker build -t lessons-content ./service/content`
- Run: `docker run -e <ENV_VARS> -p <HOST_PORT>:<CONTAINER_PORT> lessons-content`

Notes
- Review `internal/service.go` and `internal/store.go` for DB and configuration requirements.
- Ensure DB migrations in `db/migrations` have been applied before starting.
