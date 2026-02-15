# API Gateway

Short description
- Acts as the HTTP/gRPC entry point for the Lessons system. Routes requests to downstream services and handles gateway-specific logic.

Build & run (development)
- From repository root: `go run ./service/api-gateway`
- Or build: `go build -o bin/api-gateway ./service/api-gateway` then `./bin/api-gateway`

Docker
- Build: `docker build -t lessons-api-gateway ./service/api-gateway`
- Run: `docker run -e <ENV_VARS> -p <HOST_PORT>:<CONTAINER_PORT> lessons-api-gateway`

Notes
- Check `main.go` and the `Dockerfile` for required environment variables and exposed ports.
- Protobufs and generated stubs live in `proto/` and `pb/`; update/regenerate when APIs change.
