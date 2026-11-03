# Multi-stage Dockerfile for CollabFlow API Engine
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum* ./
COPY . .

# Ensure dependencies are resolved cleanly
RUN go mod tidy

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/server ./cmd/server

# Stage 2: Minimal Runtime
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /home/appuser

COPY --from=builder /bin/server /usr/local/bin/server

USER appuser

EXPOSE 8080 10000

ENTRYPOINT ["/usr/local/bin/server"]
