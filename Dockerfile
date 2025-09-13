FROM golang:1.21-apline AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o migration-service cmd/server/main.go

FROM apline:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/migration-service .

RUN adduser -D -g '' appuser && chown -R appuser /app
USER appuser

EXPOSE 8080

ENTRYPOINT ["./migration-service"]
