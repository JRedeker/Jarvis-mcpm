# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o jarvis .

# Runtime stage
FROM debian:bookworm-slim

WORKDIR /app

# Install dependencies: python3, pipx
RUN apt-get update && apt-get install -y \
    ca-certificates \
    python3 \
    pipx \
    && rm -rf /var/lib/apt/lists/*

ENV PATH="/root/.local/bin:$PATH"

# Install mcpm via pipx
RUN pipx install mcpm

COPY --from=builder /app/jarvis .

ENTRYPOINT ["./jarvis"]