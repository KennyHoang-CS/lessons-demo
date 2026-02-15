FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy full repo (build context is repository root)
COPY . .

RUN go build -o diff-service ./service/diff

FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /app/diff-service .

CMD ["./diff-service"]
