FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy full repo (build context is repository root)
COPY . .

RUN go build -o api-gateway ./service/api-gateway

FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /app/api-gateway .
COPY --from=builder /app/static ./static

CMD ["./api-gateway"]
