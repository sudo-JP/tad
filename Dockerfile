# Build stage
FROM golang:1.26-alpine AS builder

# mattn/go-sqlite3 needs cgo.
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o tad ./cmd/tad

# goose's own sqlite driver is pure Go, so it builds without cgo.
RUN CGO_ENABLED=0 go install github.com/pressly/goose/v3/cmd/goose@latest

# Runtime stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

RUN adduser -D -u 10001 appuser

WORKDIR /app

COPY --from=builder /app/tad .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY migrations ./migrations

RUN mkdir -p /app/data && chown -R appuser:appuser /app/data

USER appuser

CMD ["sh", "-c", "goose -dir ./migrations sqlite3 \"$DATABASE_URL\" up && exec ./tad"]
