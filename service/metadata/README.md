# Metadata Service

Short description
- Manages metadata for lessons (indexes, lookup data, and related storage). Provides metadata APIs used by other services.

Build & run (development)
- From repository root: `go run ./service/metadata`
- Or build: `go build -o bin/metadata ./service/metadata` then `./bin/metadata`

Docker
- Build: `docker build -t lessons-metadata ./service/metadata`
- Run: `docker run -e <ENV_VARS> -p <HOST_PORT>:<CONTAINER_PORT> lessons-metadata`

Notes
- Review `internal/service.go` and `internal/store.go` for DB requirements and `internal/errors.go` for error types.
