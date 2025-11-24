# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o jarvis .

# Runtime stage
FROM node:18-bookworm-slim

WORKDIR /app

# Install runtime dependencies if needed (e.g. git for some MCP servers)
RUN apt-get update && apt-get install -y \
    ca-certificates \
    git \
    && rm -rf /var/lib/apt/lists/*

# Copy built binary from builder
COPY --from=builder /app/jarvis .

# We need the MCPM source code to be available and linked
# Assuming the build context is the root of the repo, but standard docker build context for this file is usually Jarvis/
# If we want to bundle MCPM with Jarvis image, we need to change how this image is built or copy it in.
# For now, let's assume we copy MCPM from a 'mcpm_context' or we just don't support the full bundled image yet 
# because the "bootstrap" flow relies on local folders.

# HOWEVER, to make this image standalone-functional as requested by "working state":
# We should copy the MCPM directory from the context if possible. 
# But strictly speaking, the user's new flow is "run jarvis binary locally".
# So this Dockerfile is for a "containerized jarvis" deployment.

# Let's at least make it a valid Node environment so it COULD run MCPM if mounted/installed.
ENTRYPOINT ["./jarvis"]