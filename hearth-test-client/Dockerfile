# Stage 1: Builder
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o test-client main.go

# Stage 2: Runner
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/test-client .
# Copy templates for the web UI
COPY --from=builder /app/templates ./templates
# Copy .env (optional, but good for defaults in container)
COPY --from=builder /app/.env .

EXPOSE 3000

CMD ["./test-client"]