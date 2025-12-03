# MCPM Daemon - Runs all MCP profiles as Streamable HTTP services
# Using Python 3.11 for better compatibility with fastmcp/pydantic
# Transport: Streamable HTTP (MCP 2025-03-26 spec, replaces deprecated SSE)
FROM python:3.11-slim

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Install system dependencies, uv, and MCPM in one layer
# hadolint ignore=DL3008,DL3013
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    git \
    nodejs \
    npm \
    supervisor \
    && rm -rf /var/lib/apt/lists/* \
    && curl -LsSf https://astral.sh/uv/install.sh | sh \
    && pip install --no-cache-dir "pydantic==2.11.10" "fastmcp==2.10.2" mcpm \
    && mkdir -p /var/log/mcpm /etc/supervisor/conf.d

ENV PATH="/root/.local/bin:$PATH"

# Copy entrypoint
COPY --chmod=755 entrypoint.sh /entrypoint.sh

# Expose ports for profiles
# p-pokeedge: 6276, memory: 6277, morph: 6278, qdrant: 6279, p-new: 6280
EXPOSE 6276 6277 6278 6279 6280

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:6276/health || exit 1

ENTRYPOINT ["/entrypoint.sh"]
