# Port Allocation Matrix

**Version:** 2.0
**Date:** 2025-11-22
**Status:** Active

---

## Overview

This document defines the port allocation for the MCPM-managed environment. Since `mcpm` primarily manages local processes via `stdio`, the need for extensive port mapping is significantly reduced compared to the legacy gateway architecture.

## Port Allocation Table

| Service | Port | Protocol | Transport | Status | Notes | Conflict Risk |
|---------|------|----------|-----------|--------|--------|---------------|
| **PostgreSQL** | 5432 | TCP | SQL | ⚪ Optional | Backend for specific tools | Low |
| **Qdrant** | 6333 | HTTP | HTTP | ⚪ Available | Vector database for memory | Low |
| **Remote MCP** | 8080 | HTTP/SSE | SSE | ⚪ Optional | If running a remote MCP server | Medium |

## Port Status Legend

- ✅ **Active**: Currently in use by default configuration
- ⚪ **Optional**: Only used if specific services are enabled
- ❌ **Avoid**: Known conflicts or reserved ports

## Service Details

### **PostgreSQL - Port 5432**
- **Default Port**: 5432
- **Use Cases**:
  - Backend storage for `sqlite` or custom database tools.
  - Metadata storage for enterprise deployments.
- **Configuration**:
  Managed via Docker Compose or local installation.

### **Qdrant - Port 6333**
- **Default Port**: 6333
- **Use Cases**:
  - Vector storage for the `memory` MCP server.
  - Semantic search capabilities.
- **Configuration**:
  Managed via Docker Compose.

### **Remote MCP Servers - Port 8080+**
- **Default Port**: 8080 (Common convention)
- **Use Cases**:
  - Hosting an MCP server over HTTP/SSE for remote access.
  - `mcpm` can register these remote endpoints, but does not manage the port binding itself.

## Port Conflict Analysis

### **Common Port Conflicts**

| Port | Common Service | Risk Level | Mitigation |
|------|----------------|------------|------------|
| 5432 | Local PostgreSQL | Low | Use Docker port mapping (e.g., 5433:5432) if local DB exists. |
| 6333 | Local Qdrant | Low | Use Docker port mapping. |
| 8080 | Web Servers | Medium | Configure remote servers to use alternative ports (e.g., 3000, 8081). |

## Development Environment

### **Checking Port Availability**

Before starting Docker services, ensure the required ports are free:

```bash
# Check for listeners on common ports
netstat -tulpn | grep -E ':(5432|6333|8080)'
```

### **Docker Configuration**

If you encounter conflicts, modify the `docker-compose.yml` to map to different host ports:

```yaml
services:
  postgres:
    ports:
      - "5433:5432"  # Maps host 5433 to container 5432
```
