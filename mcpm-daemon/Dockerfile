# MCPM Daemon - Runs all MCP profiles as SSE services
# Using Python 3.11 for better compatibility with fastmcp/pydantic
FROM python:3.11-slim

# Install system dependencies
RUN apt-get update && apt-get install -y \
    curl \
    git \
    nodejs \
    npm \
    supervisor \
    && rm -rf /var/lib/apt/lists/*

# Install uv for Python package management
RUN curl -LsSf https://astral.sh/uv/install.sh | sh
ENV PATH="/root/.local/bin:$PATH"

# Install MCPM with exact versions that are known to work (from host)
RUN pip install "pydantic==2.11.10" "fastmcp==2.10.2" mcpm

# Create directories
RUN mkdir -p /var/log/mcpm /etc/supervisor/conf.d

# Copy supervisor config (will be generated at runtime)
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Expose ports for profiles
# p-pokeedge: 6276, memory: 6277, morph: 6278, qdrant: 6279, p-new: 6280
EXPOSE 6276 6277 6278 6279 6280

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:6276/health || exit 1

ENTRYPOINT ["/entrypoint.sh"]
