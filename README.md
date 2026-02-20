# About 
Demo utilizing rest api, grpcs for backend <-> backend communication. Next on the to-do list: work on the infrastructure, such as adding tests (i.e unit tests), remove hard coded env values, and more.  

Time spent on this: 90-minutes 

# Lesson Versioning System

A scalable, multi-service Go system for lesson versioning with:

- Fiber API Gateway (HTTP)
- gRPC backend services
- Postgres for metadata
- MinIO (S3-compatible) for content, the blob storage. 
- Redis for diff caching

## Architecture

- `api-gateway` (Fiber, HTTP)
- `content-service` (gRPC, MinIO)
- `metadata-service` (gRPC, Postgres)
- `diff-service` (gRPC, Redis + ContentService)

## Running locally

1. Generate protobufs:

```bash
protoc --go_out=. --go-grpc_out=. proto/lesson.proto
```

## Run with Docker Compose

Prerequisites: Docker and Docker Compose installed.

Start the full stack (build images):

```bash
docker compose up --build
```

The gateway will be available at: `http://127.0.0.1:8080`.

Notes:
- The repository includes a `docker-compose.yml` that builds service images from the repository root. The `api-gateway` service mounts the local `static/` directory so you can iterate on frontend files without rebuilding the image.
- If you change Go code, rebuild the specific service image: `docker compose up -d --build api-gateway` (replace service name as needed).

## Database migrations

Apply migrations before creating content (if you started with a fresh DB):

From your host (requires `psql`):

```bash
psql "postgres://lessons:lessons@localhost:5432/lessons" -f db/migrations/001_init.sql
```

Or run inside the running Postgres container:

```bash
docker compose exec -T postgres psql -U lessons -d lessons -f /var/lib/postgresql/data/docker-entrypoint-initdb.d/001_init.sql
```

## QA — manual and automated

Manual (web UI):
- Open `http://127.0.0.1:8080`.
- Use the forms to: Create Lesson, Create Version, Publish Version, Get Version, and Diff. Each form shows the HTTP status and response in the result box.

Automated QA scripts:
- Bash (Linux/macOS/WSL): `bash qa/run_qa.sh`
- PowerShell (Windows): `.\qa\run_qa.ps1`

The scripts perform the full flow (create lesson → create version → publish → get → diff) and print responses. If any step fails they exit non-zero and print the response body.

## Troubleshooting
- View logs: `docker compose logs -f api-gateway` or `docker compose logs -f metadata`.
- If the API returns DB errors, ensure migrations were applied and Postgres is healthy.
- If static files don't reflect changes, confirm `./static` is present and the `api-gateway` service has the bind mount in `docker-compose.yml`.
