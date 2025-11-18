# MCPJungle Installation Guide

**Version:** 1.0
**Date:** 2025-11-18
**Status:** Installation Reference

---

## Overview

This guide provides step-by-step instructions for installing and setting up MCPJungle (jarvis) on various platforms and deployment methods.

## Prerequisites

### **System Requirements**
- **OS**: macOS, Linux, or Windows (with WSL)
- **Memory**: 2GB RAM minimum, 4GB recommended
- **Storage**: 1GB for installation + additional for data
- **Network**: Internet access for downloading packages

### **Software Dependencies**
- **Docker** (recommended) or **Node.js** 18+
- **Homebrew** (macOS) or package manager
- **Git** for cloning examples

---

## Installation Methods

### **Method 1: Homebrew (Recommended for macOS)**

```bash
# Add MCPJungle tap
brew tap mcpjungle/mcpjungle

# Install MCPJungle
brew install mcpjungle/mcpjungle/mcpjungle

# Verify installation
mcpjungle version

# Expected output: MCPJungle version 0.2.16 or later
```

### **Method 2: Direct Binary Download**

```bash
# Download latest release
curl -L -o mcpjungle \
  https://github.com/mcpjungle/MCPJungle/releases/latest/download/mcpjungle-darwin-amd64

# Make executable
chmod +x mcpjungle

# Move to PATH
sudo mv mcpjungle /usr/local/bin/

# Verify installation
mcpjungle version
```

### **Method 3: Docker (Recommended for Production)**

```bash
# Pull the image
docker pull mcpjungle/mcpjungle:latest-stdio

# Verify image
docker images | grep mcpjungle
```

---

## Quick Start with Docker Compose

### **Development Setup**

1. **Create docker-compose.yml**:
```yaml
version: '3.8'

services:
  mcpjungle:
    image: mcpjungle/mcpjungle:latest-stdio
    ports:
      - "8080:8080"
    volumes:
      - ./data:/data
      - .:/host  # For filesystem access
    environment:
      - DATABASE_URL=sqlite:///data/mcpjungle.db
      - LOG_LEVEL=info
    restart: unless-stopped

  postgres:
    image: postgres:15
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      - POSTGRES_DB=mcpjungle_db
      - POSTGRES_USER=mcpjungle
      - POSTGRES_PASSWORD=secure_password
    restart: unless-stopped

volumes:
  postgres_data:
```

2. **Start the services**:
```bash
# Start services
docker compose up -d

# Check status
docker compose ps

# View logs
docker compose logs -f mcpjungle
```

3. **Test the installation**:
```bash
# Test health endpoint
curl http://localhost:8080/health

# Should return: {"status":"healthy"}
```

### **Production Setup**

1. **Create docker-compose.prod.yml**:
```yaml
version: '3.8'

services:
  mcpjungle:
    image: mcpjungle/mcpjungle:latest-stdio
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://mcpjungle:secure_password@postgres:5432/mcpjungle_db
      - SERVER_MODE=enterprise
      - OTEL_ENABLED=true
      - LOG_LEVEL=warn
    depends_on:
      - postgres
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  postgres:
    image: postgres:15
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      - POSTGRES_DB=mcpjungle_db
      - POSTGRES_USER=mcpjungle
      - POSTGRES_PASSWORD=secure_password
    restart: unless-stopped

volumes:
  postgres_data:
```

2. **Start production services**:
```bash
# Start production setup
docker compose -f docker-compose.prod.yml up -d

# Verify deployment
docker compose -f docker-compose.prod.yml ps
```

---

## Manual Installation (Without Docker)

### **Step 1: Download Binary**

```bash
# For macOS (Intel)
curl -L -o mcpjungle \
  https://github.com/mcpjungle/MCPJungle/releases/latest/download/mcpjungle-darwin-amd64

# For macOS (Apple Silicon)
curl -L -o mcpjungle \
  https://github.com/mcpjungle/MCPJungle/releases/latest/download/mcpjungle-darwin-arm64

# For Linux (AMD64)
curl -L -o mcpjungle \
  https://github.com/mcpjungle/MCPJungle/releases/latest/download/mcpjungle-linux-amd64

# For Linux (ARM64)
curl -L -o mcpjungle \
  https://github.com/mcpjungle/MCPJungle/releases/latest/download/mcpjungle-linux-arm64
```

### **Step 2: Install Binary**

```bash
# Make executable
chmod +x mcpjungle

# Move to system PATH
sudo mv mcpjungle /usr/local/bin/

# Verify installation
mcpjungle version
```

### **Step 3: Start Server**

```bash
# Development mode (SQLite)
mcpjungle start

# Production mode (PostgreSQL)
export DATABASE_URL=postgres://user:pass@localhost:5432/mcpjungle_db
mcpjungle start --enterprise

# With custom port
export PORT=8081
mcpjungle start
```

---

## Verification Steps

### **Basic Health Check**
```bash
# Test health endpoint
curl http://localhost:8080/health

# Expected output: {"status":"healthy"}

# Test with custom port
curl http://localhost:8081/health
```

### **CLI Verification**
```bash
# Check version
mcpjungle version

# List available commands
mcpjungle --help

# Test CLI functionality
mcpjungle list servers
```

### **Service Status**
```bash
# Check if service is running
ps aux | grep mcpjungle

# Check port binding
netstat -tulpn | grep 8080

# Check Docker status
docker ps | grep mcpjungle
```

---

## Configuration

### **Environment Variables**
```bash
# Core configuration
export PORT=8080                    # Server port
export DATABASE_URL=sqlite:///data/mcpjungle.db  # Database connection
export LOG_LEVEL=info               # Logging level
export SERVER_MODE=development      # Server mode

# Enterprise features
export OTEL_ENABLED=true            # OpenTelemetry
export OTEL_RESOURCE_ATTRIBUTES=deployment.environment.name=dev
```

### **Configuration File**
Create `jarvis.config.json`:
```json
{
  "port": 8080,
  "database": {
    "url": "sqlite:///data/mcpjungle.db"
  },
  "logging": {
    "level": "info"
  },
  "enterprise": {
    "enabled": false
  }
}
```

---

## Troubleshooting

### **Common Issues**

#### **Port Already in Use**
```bash
# Find process using port 8080
lsof -i :8080

# Kill process (careful!)
kill -9 <PID>

# Or use different port
export PORT=8081
mcpjungle start
```

#### **Permission Denied**
```bash
# Check file permissions
ls -la /usr/local/bin/mcpjungle

# Fix permissions
sudo chmod +x /usr/local/bin/mcpjungle

# Check if running as root is required
sudo mcpjungle start
```

#### **Database Connection Issues**
```bash
# Test PostgreSQL connection
psql $DATABASE_URL -c "SELECT 1;"

# Check SQLite file permissions
ls -la data/mcpjungle.db

# Create data directory
mkdir -p data
chmod 755 data
```

#### **Binary Not Found**
```bash
# Check PATH
echo $PATH

# Add to PATH if needed
export PATH=$PATH:/usr/local/bin

# Or use full path
/usr/local/bin/mcpjungle version
```

### **Docker Issues**

#### **Container Won't Start**
```bash
# Check logs
docker compose logs mcpjungle

# Check port conflicts
docker compose ps

# Restart services
docker compose restart
```

#### **Permission Issues**
```bash
# Fix file permissions
sudo chown -R $USER:$USER ./data

# Check Docker permissions
docker info
```

---

## Security Setup

### **Development Environment**
- No authentication required
- Localhost access only
- All tools available

### **Production Environment**
```bash
# Enable enterprise mode
export SERVER_MODE=enterprise

# Set up authentication (future feature)
export AUTH_ENABLED=true
export JWT_SECRET=your-secret-key

# Configure access control
export ADMIN_TOKEN=admin-secret-token
```

---

## Performance Tuning

### **Development**
```bash
# Reduce logging
export LOG_LEVEL=warn

# Disable metrics
export OTEL_ENABLED=false
```

### **Production**
```bash
# Enable all features
export SERVER_MODE=enterprise
export OTEL_ENABLED=true
export LOG_LEVEL=error

# Connection pooling
export MAX_CONNECTIONS=100
export CONNECTION_TIMEOUT=30000
```

---

## Backup & Recovery

### **Configuration Backup**
```bash
# Backup server registrations
mcpjungle list servers --json > servers-backup.json

# Backup tool groups
mcpjungle list groups --json > groups-backup.json
```

### **Database Backup**
```bash
# PostgreSQL backup
pg_dump mcpjungle_db > mcpjungle-backup.sql

# SQLite backup
cp data/mcpjungle.db data/mcpjungle-backup.db
```

---

## Uninstallation

### **Homebrew Uninstall**
```bash
brew uninstall mcpjungle
brew untap mcpjungle/mcpjungle
```

### **Manual Uninstall**
```bash
# Remove binary
sudo rm /usr/local/bin/mcpjungle

# Remove data (careful!)
rm -rf data/
rm -rf ~/.mcpjungle/
```

### **Docker Cleanup**
```bash
# Stop and remove containers
docker compose down

# Remove images
docker rmi mcpjungle/mcpjungle:latest-stdio

# Remove volumes (careful!)
docker volume prune
```

---

## Next Steps

After successful installation:

1. **Verify installation** using health checks
2. **Register MCP servers** using CLI commands
3. **Configure tool groups** for organization
4. **Test IDE connections** to jarvis
5. **Set up monitoring** for production use

---

## Version History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-11-18 | Kilo Code | Initial installation guide |
| 1.1 | TBD | TBD | Add troubleshooting section |
| 1.2 | TBD | TBD | Add performance tuning |

**Status**: ✅ **Complete** - Installation procedures documented

**Next Steps**:
1. Test installation procedures during Phase 1
2. Update guide with any issues discovered
3. Add platform-specific troubleshooting
