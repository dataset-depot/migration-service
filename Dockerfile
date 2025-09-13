# Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git and ca-certificates for go modules
RUN apk add --no-cache git ca-certificates

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o migration-service \
    cmd/server/main.go

# Final stage
FROM alpine:3.19

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/migration-service .

# Create non-root user
RUN adduser -D -g '' appuser && chown -R appuser /app
USER appuser

EXPOSE 8080

ENTRYPOINT ["./migration-service"]
