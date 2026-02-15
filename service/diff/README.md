# Diff Service

Short description
- Handles diffing/caching responsibilities used by the Lessons system. Contains in-memory/cache helpers and diff logic.

Build & run (development)
- From repository root: `go run ./service/diff`
- Or build: `go build -o bin/diff ./service/diff` then `./bin/diff`

Docker
- Build: `docker build -t lessons-diff ./service/diff`
- Run: `docker run -e <ENV_VARS> -p <HOST_PORT>:<CONTAINER_PORT> lessons-diff`

Notes
- See `internal/cache.go` and `internal/service.go` for cache configuration (Redis, in-memory) and expected env vars.
