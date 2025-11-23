# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o jarvis .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install mcpm dependencies if needed (or we mount the binary)
# For now, we assume we will mount the host's mcpm binary or config.
# But to be safe, let's ensure basic tools are present.
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/jarvis .

ENTRYPOINT ["./jarvis"]