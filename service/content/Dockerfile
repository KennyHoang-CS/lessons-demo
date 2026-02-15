FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy full repo (build context is repository root)
COPY . .

RUN go build -o content-service ./service/content

FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /app/content-service .

CMD ["./content-service"]
