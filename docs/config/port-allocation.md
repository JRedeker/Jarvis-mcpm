# Port Allocation Matrix

**Version:** 1.0
**Date:** 2025-11-18
**Status:** Configuration Reference

---

## Overview

This document defines the complete port allocation for the simplified MCPJungle architecture, ensuring no conflicts and clear service identification.

## Port Allocation Table

| Service | Port | Protocol | Transport | Status | Notes | Conflict Risk |
|---------|------|----------|-----------|--------|--------|---------------|
| **MCPJungle (jarvis)** | 8080 | HTTP/WS | streamable-http | ✅ Planned | Primary MCP endpoint | Low |
| **PostgreSQL** | 5432 | TCP | SQL | ⚪ Optional | Production metadata | Low |
| **Qdrant** | 6333 | HTTP | HTTP | ⚪ Available | Vector database (Phase 3) | Low |
| **Cipher** | 3021 | HTTP | stdio | ⚪ Available | Memory-only mode (Phase 3) | Low |

## Port Status Legend

- ✅ **Planned**: Will be used in current implementation
- ⚪ **Available**: Can be used if needed in future phases
- ❌ **Avoid**: Known conflicts or reserved ports
- 🔧 **Configurable**: Can be changed via configuration

## Service Details

### **MCPJungle (jarvis) - Port 8080**
- **Default Port**: 8080 (configurable)
- **Protocol**: HTTP for REST API, WebSocket for IDE connections
- **Use Cases**:
  - Tool discovery: `GET /mcp`
  - Tool invocation: `POST /mcp`
  - Health checks: `GET /health`
  - Metrics: `GET /metrics` (enterprise mode)
- **Configuration**:
  ```bash
  # Environment variable
  export PORT=8080

  # Docker
  ports:
    - "8080:8080"
  ```

### **PostgreSQL - Port 5432**
- **Default Port**: 5432 (standard PostgreSQL)
- **Protocol**: TCP
- **Use Cases**:
  - Server metadata storage
  - Tool group definitions
  - Access control data (enterprise)
  - Usage analytics
- **Configuration**:
  ```bash
  # Environment variable
  export DATABASE_URL=postgres://user:pass@localhost:5432/mcpjungle_db

  # Docker
  ports:
    - "5432:5432"
  ```

### **Qdrant - Port 6333**
- **Default Port**: 6333 (standard Qdrant)
- **Protocol**: HTTP REST API
- **Use Cases**:
  - Vector storage for advanced memory (Phase 3)
  - Semantic search capabilities
  - Embedding storage
- **Configuration**:
  ```bash
  # Environment variable
  export VECTOR_STORE_URL=http://localhost:6333

  # Docker
  ports:
    - "6333:6333"
  ```

### **Cipher - Port 3021**
- **Default Port**: 3021 (Cipher default)
- **Protocol**: HTTP for MCP stdio transport
- **Use Cases**:
  - Advanced memory tools (Phase 3)
  - Vector search and reasoning
  - Workspace memory
- **Configuration**:
  ```bash
  # Command line
  cipher --mode mcp --mcp-port 3021

  # Environment
  export MCP_PORT=3021
  ```

## Port Conflict Analysis

### **Common Port Conflicts**

| Port | Common Service | Risk Level | Mitigation |
|------|----------------|------------|------------|
| 3000 | Node.js dev servers | Medium | Use 8080 for jarvis |
| 3001 | React dev servers | Medium | Use 8080 for jarvis |
| 5000 | Flask/Python servers | Medium | Use 8080 for jarvis |
| 5432 | PostgreSQL | Low | Standard port, usually available |
| 6333 | Qdrant | Low | Specialized service |
| 8080 | Generic HTTP | Medium | Most common conflict |

### **Conflict Resolution**

#### **If Port 8080 is Unavailable**
```bash
# Option 1: Use different port
export PORT=8081
mcpjungle start

# Option 2: Check what's using it
lsof -i :8080
netstat -tulpn | grep 8080

# Option 3: Kill conflicting process
kill -9 <PID>
```

#### **If PostgreSQL Port 5432 is Unavailable**
```bash
# Option 1: Use different PostgreSQL port
docker run -p 5433:5432 postgres:15

# Option 2: Update connection string
export DATABASE_URL=postgres://user:pass@localhost:5433/mcpjungle_db
```

#### **If Qdrant Port 6333 is Unavailable**
```bash
# Option 1: Use different Qdrant port
docker run -p 6334:6333 qdrant/qdrant

# Option 2: Update environment
export VECTOR_STORE_URL=http://localhost:6334
```

## Network Configuration

### **Firewall Rules**
```bash
# Allow jarvis (if external access needed)
sudo ufw allow 8080/tcp

# Allow PostgreSQL (production only)
sudo ufw allow 5432/tcp

# Allow Qdrant (if external access needed)
sudo ufw allow 6333/tcp
```

### **Docker Networking**
```yaml
# docker-compose.yml
networks:
  mcp-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16

services:
  mcpjungle:
    networks:
      - mcp-network
    ports:
      - "8080:8080"

  postgres:
    networks:
      - mcp-network
    ports:
      - "5432:5432"

  qdrant:
    networks:
      - mcp-network
    ports:
      - "6333:6333"
```

## Development Environment

### **Local Development Ports**
```bash
# Check all ports in use
netstat -tulpn | grep -E ':(8080|5432|6333|3021)'

# Check specific port
lsof -i :8080

# Check port availability
nc -z localhost 8080 && echo "Port 8080 in use" || echo "Port 8080 available"
```

### **Port Monitoring Script**
```bash
#!/bin/bash
# save as check-ports.sh

PORTS=(8080 5432 6333 3021)
echo "Checking MCPJungle port availability..."

for port in "${PORTS[@]}"; do
    if lsof -i :$port > /dev/null 2>&1; then
        echo "❌ Port $port - IN USE"
        lsof -i :$port | head -2
    else
        echo "✅ Port $port - AVAILABLE"
    fi
done
```

## Production Considerations

### **Load Balancing**
- **jarvis**: Can be load balanced behind nginx/traefik
- **PostgreSQL**: Use connection pooling (PgBouncer)
- **Qdrant**: Built-in clustering support

### **Security**
- **8080**: Should be behind reverse proxy in production
- **5432**: Should not be exposed externally
- **6333**: Should be internal network only

### **Monitoring**
```bash
# Monitor port usage
watch -n 5 'netstat -tulpn | grep -E ":(8080|5432|6333|3021)"'

# Log port access
sudo tcpdump -i any port 8080 -w jarvis-traffic.pcap
```

## Configuration Management

### **Environment Variables**
```bash
# .env file
MCPJUNGLE_PORT=8080
POSTGRES_PORT=5432
QDRANT_PORT=6333
CIPHER_PORT=3021
```

### **Docker Environment**
```yaml
# docker-compose.override.yml
version: '3.8'
services:
  mcpjungle:
    ports:
      - "${MCPJUNGLE_PORT:-8080}:8080"

  postgres:
    ports:
      - "${POSTGRES_PORT:-5432}:5432"

  qdrant:
    ports:
      - "${QDRANT_PORT:-6333}:6333"
```

## Troubleshooting

### **Port Already in Use**
```bash
# Find process using port
sudo lsof -i :8080

# Kill process (careful!)
kill -9 <PID>

# Or use different port
export PORT=8081
mcpjungle start
```

### **Permission Denied**
```bash
# Check if port requires root
sudo netstat -tulpn | grep 8080

# Use higher port (above 1024)
export PORT=18080
mcpjungle start
```

### **Connection Refused**
```bash
# Check if service is running
curl http://localhost:8080/health

# Check firewall
sudo ufw status

# Check Docker port mapping
docker ps | grep mcpjungle
```

## Version History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-11-18 | Kilo Code | Initial port allocation matrix |
| 1.1 | TBD | TBD | Add production considerations |
| 1.2 | TBD | TBD | Add troubleshooting section |

**Status**: ✅ **Complete** - Port allocation defined and documented

**Next Steps**:
1. Verify ports during Phase 1 installation
2. Update matrix with actual usage findings
3. Add any new ports discovered during implementation
