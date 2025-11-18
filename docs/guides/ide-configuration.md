# IDE Configuration Guide for MCPJungle

**Version:** 1.0
**Date:** 2025-11-18
**Status:** Configuration Reference

---

## Overview

This guide provides step-by-step instructions for configuring popular IDEs to connect to MCPJungle (jarvis) as the single MCP aggregation point.

## Supported IDEs

- **Cline** (VS Code extension)
- **Kilo Code** (Cursor-based IDE)
- **Claude Desktop** (standalone application)
- **Other MCP-compatible editors**

---

## General Configuration Principles

### **Connection Details**
- **URL**: `http://localhost:8080/mcp`
- **Transport**: Streamable-HTTP (primary)
- **Authentication**: None (development) / Bearer token (production)
- **Protocol**: MCP 2024-11-05 specification

### **Configuration Format**
All IDEs use JSON configuration files with similar structure:
```json
{
  "mcpServers": {
    "mcpjungle": {
      "url": "http://localhost:8080/mcp",
      "transport": "streamable-http"
    }
  }
}
```

---

## Cline Configuration

### **VS Code Settings**
1. Open VS Code Settings (Ctrl/Cmd + ,)
2. Search for "cline mcp"
3. Edit settings.json directly:

```json
{
  "cline.mcp": {
    "mcpServers": {
      "mcpjungle": {
        "url": "http://localhost:8080/mcp",
        "transport": "streamable-http",
        "timeout": 30000
      }
    }
  }
}
```

### **Manual Configuration File**
Create or edit `~/.vscode/settings.json`:
```json
{
  "cline.mcp": {
    "mcpServers": {
      "mcpjungle": {
        "url": "http://localhost:8080/mcp",
        "transport": "streamable-http",
        "timeout": 30000,
        "retryAttempts": 3,
        "retryDelay": 1000
      }
    }
  }
}
```

### **Workspace-Specific Configuration**
Create `.vscode/settings.json` in your project:
```json
{
  "cline.mcp": {
    "mcpServers": {
      "mcpjungle": {
        "url": "http://localhost:8080/mcp",
        "transport": "streamable-http"
      }
    }
  }
}
```

---

## Kilo Code Configuration

### **Application Settings**
1. Open Kilo Code Settings
2. Navigate to MCP Configuration
3. Add new MCP server:

```json
{
  "mcpServers": {
    "mcpjungle": {
      "url": "http://localhost:8080/mcp",
      "transport": "streamable-http"
    }
  }
}
```

### **Configuration File Location**
Edit `~/.cursor/mcp.json` or similar configuration file:
```json
{
  "mcpServers": {
    "mcpjungle": {
      "url": "http://localhost:8080/mcp",
      "transport": "streamable-http",
      "headers": {
        "User-Agent": "KiloCode/1.0"
      }
    }
  }
}
```

---

## Claude Desktop Configuration

### **Configuration File Location**
Edit `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS):
```json
{
  "mcpServers": {
    "mcpjungle": {
      "url": "http://localhost:8080/mcp",
      "transport": "streamable-http"
    }
  }
}
```

### **Windows Configuration**
Edit `%APPDATA%\Claude\claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "mcpjungle": {
      "url": "http://localhost:8080/mcp",
      "transport": "streamable-http"
    }
  }
}
```

---

## Production Configuration (with Authentication)

### **Bearer Token Authentication**
```json
{
  "mcpServers": {
    "mcpjungle": {
      "url": "http://localhost:8080/mcp",
      "transport": "streamable-http",
      "headers": {
        "Authorization": "Bearer your-access-token"
      }
    }
  }
}
```

### **Custom Headers**
```json
{
  "mcpServers": {
    "mcpjungle": {
      "url": "http://localhost:8080/mcp",
      "transport": "streamable-http",
      "headers": {
        "X-Client-ID": "cline-v1.0",
        "X-Environment": "development"
      }
    }
  }
}
```

---

## Multi-Environment Configuration

### **Development Environment**
```json
{
  "mcpServers": {
    "mcpjungle-dev": {
      "url": "http://localhost:8080/mcp",
      "transport": "streamable-http"
    }
  }
}
```

### **Staging Environment**
```json
{
  "mcpServers": {
    "mcpjungle-staging": {
      "url": "https://staging.jarvis.example.com/mcp",
      "transport": "streamable-http",
      "headers": {
        "Authorization": "Bearer staging-token"
      }
    }
  }
}
```

### **Production Environment**
```json
{
  "mcpServers": {
    "mcpjungle-prod": {
      "url": "https://jarvis.example.com/mcp",
      "transport": "streamable-http",
      "headers": {
        "Authorization": "Bearer production-token"
      }
    }
  }
}
```

---

## Tool Groups Configuration

### **Development Tool Group**
Access only essential development tools:
```json
{
  "mcpServers": {
    "mcpjungle-dev": {
      "url": "http://localhost:8080/v0/groups/development/mcp",
      "transport": "streamable-http"
    }
  }
}
```

### **Research Tool Group**
Access research and search tools:
```json
{
  "mcpServers": {
    "mcpjungle-research": {
      "url": "http://localhost:8080/v0/groups/research/mcp",
      "transport": "streamable-http"
    }
  }
}
```

---

## Troubleshooting

### **Connection Issues**

#### **Service Not Available**
```bash
# Check if jarvis is running
curl http://localhost:8080/health

# Check port binding
netstat -tulpn | grep 8080

# Check Docker status
docker ps | grep mcpjungle
```

#### **Connection Timeout**
```json
{
  "mcpServers": {
    "mcpjungle": {
      "url": "http://localhost:8080/mcp",
      "transport": "streamable-http",
      "timeout": 60000,  // 60 seconds
      "retryAttempts": 3,
      "retryDelay": 2000
    }
  }
}
```

#### **Authentication Errors**
```bash
# Check token validity
curl -H "Authorization: Bearer your-token" \
  http://localhost:8080/mcp

# Verify token in jarvis logs
docker compose logs mcpjungle | grep -i auth
```

### **Tool Discovery Issues**

#### **Tools Not Appearing**
```bash
# List tools via CLI
mcpjungle list tools

# Check specific server
mcpjungle list tools --server context7

# Verify server registration
mcpjungle list servers
```

#### **Tool Names**
Tools follow the pattern: `<server-name>__<tool-name>`
- `context7__get-library-docs`
- `brave-search__brave_web_search`
- `filesystem__read_file`

---

## Testing Configuration

### **Basic Connectivity Test**
```bash
# Test health endpoint
curl http://localhost:8080/health

# Test tool listing
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/list","id":1}'
```

### **Tool Invocation Test**
```bash
# Test specific tool
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "context7__get-library-docs",
      "arguments": {"library": "lodash/lodash"}
    },
    "id": 1
  }'
```

---

## Migration from Direct MCP Connections

### **Before (Direct Connections)**
```json
{
  "mcpServers": {
    "context7": {
      "url": "https://mcp.context7.com/mcp",
      "transport": "streamable-http"
    },
    "brave-search": {
      "command": "npx",
      "args": ["-y", "@brave/brave-search-mcp-server"],
      "env": {"BRAVE_API_KEY": "key"}
    }
  }
}
```

### **After (Via jarvis)**
```json
{
  "mcpServers": {
    "mcpjungle": {
      "url": "http://localhost:8080/mcp",
      "transport": "streamable-http"
    }
  }
}
```

---

## Performance Optimization

### **Connection Pooling**
```json
{
  "mcpServers": {
    "mcpjungle": {
      "url": "http://localhost:8080/mcp",
      "transport": "streamable-http",
      "connectionPool": {
        "maxSize": 10,
        "idleTimeout": 30000
      }
    }
  }
}
```

### **Caching**
```json
{
  "mcpServers": {
    "mcpjungle": {
      "url": "http://localhost:8080/mcp",
      "transport": "streamable-http",
      "cache": {
        "enabled": true,
        "ttl": 300000  // 5 minutes
      }
    }
  }
}
```

---

## Security Considerations

### **Local Development**
- Use localhost connections
- No authentication required
- All tools available

### **Production Deployment**
- Use HTTPS connections
- Implement authentication
- Restrict tool access
- Monitor usage patterns

### **Network Security**
```bash
# Firewall rules
sudo ufw allow from 127.0.0.1 to any port 8080
sudo ufw deny 8080  # Block external access

# Docker network isolation
docker network create mcp-network
docker network connect mcp-network mcpjungle
```

---

## Monitoring and Logging

### **IDE-Side Logging**
Enable debug logging in IDE settings to troubleshoot MCP communication.

### **Server-Side Monitoring**
```bash
# Monitor jarvis logs
docker compose logs -f mcpjungle

# Check metrics (if enabled)
curl http://localhost:8080/metrics

# Monitor connection patterns
docker compose logs mcpjungle | grep -E "connected|disconnected|error"
```

---

## Backup and Recovery

### **Configuration Backup**
```bash
# Backup IDE settings
cp ~/.vscode/settings.json settings-backup.json
cp ~/.cursor/mcp.json mcp-backup.json

# Backup jarvis configurations
mcpjungle list servers --json > servers-backup.json
mcpjungle list groups --json > groups-backup.json
```

### **Recovery Procedures**
1. Restore jarvis from backup
2. Re-register MCP servers
3. Update IDE configurations
4. Test connectivity

---

## Advanced Features

### **Dynamic Tool Groups**
```json
{
  "mcpServers": {
    "mcpjungle-dynamic": {
      "url": "http://localhost:8080/v0/groups/dynamic/mcp",
      "transport": "streamable-http",
      "headers": {
        "X-Context": "development",
        "X-User-Role": "developer"
      }
    }
  }
}
```

### **Multi-Server Setup**
```json
{
  "mcpServers": {
    "mcpjungle-primary": {
      "url": "http://localhost:8080/mcp",
      "transport": "streamable-http",
      "priority": 1
    },
    "mcpjungle-backup": {
      "url": "http://backup.localhost:8081/mcp",
      "transport": "streamable-http",
      "priority": 2
    }
  }
}
```

---

## Version History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-11-18 | Kilo Code | Initial IDE configuration guide |
| 1.1 | TBD | TBD | Add troubleshooting section |
| 1.2 | TBD | TBD | Add advanced features |

**Status**: ✅ **Complete** - IDE configuration procedures documented

**Next Steps**:
1. Test configurations with actual IDEs during Phase 4
2. Update with any IDE-specific issues discovered
3. Add screenshots and visual guides
