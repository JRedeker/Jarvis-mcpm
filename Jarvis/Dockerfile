# Jarvis MCP Server - Dockerfile
# Version: 3.0 (Consolidated Tools)
#
# This builds a containerized Jarvis that can be used as a standalone MCP server.
# For full functionality, mount the MCPM config directory.

# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
# hadolint ignore=DL3018
RUN apk add --no-cache git

# Download dependencies first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o jarvis .

# Runtime stage - using Node.js 22 LTS for MCPM compatibility
FROM node:22-bookworm-slim

LABEL org.opencontainers.image.title="Jarvis MCP Server"
LABEL org.opencontainers.image.description="Intelligent MCP management gateway with 8 consolidated tools"
LABEL org.opencontainers.image.version="3.0"
LABEL org.opencontainers.image.source="https://github.com/pathintegral-institute/mcpm.sh"

WORKDIR /app

# Install runtime dependencies
# hadolint ignore=DL3008
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    git \
    && rm -rf /var/lib/apt/lists/*

# Copy built binary from builder
COPY --from=builder /app/jarvis .

# Create config directory
RUN mkdir -p /root/.config/mcpm

# Health check - Jarvis responds to MCP initialize
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD echo '{"jsonrpc":"2.0","method":"ping","id":1}' | timeout 5 ./jarvis 2>/dev/null | grep -q "pong" || exit 1

# Run as stdio MCP server
ENTRYPOINT ["./jarvis"]
