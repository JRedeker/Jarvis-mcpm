# Routing Metadata MCP Server Usage Guide

## Overview
The routing-metadata MCP server provides intelligent tool routing analytics and session management capabilities for the cipher-aggregator system. It helps track tool usage patterns, validate routing decisions, and provide analytics for routing optimization.

## Architecture
- **Transport**: JSON-RPC over stdio
- **Framework**: MCP Python SDK (async-native)
- **Logging**: OpenTelemetry with graceful fallback
- **Environment**: Uses eval_llm_venv for consistent dependencies

## Installation & Setup

### Dependencies
The server requires Python packages that are installed in the eval_llm_venv:
```bash
# Dependencies are pre-installed in eval_llm_venv
/home/jrede/dev/MCP/eval_llm_venv/bin/python -c "import mcp, opentelemetry, yaml; print('All dependencies available')"
```

### Configuration
The server is configured in `cipher.yml`:
```yaml
routing-metadata:
  type: stdio
  command: "/home/jrede/dev/MCP/eval_llm_venv/bin/python"
  args:
    - /home/jrede/dev/MCP/servers/routing-metadata-mcp.py
  env:
    OTEL_LOGS_EXPORTER: "otlp"
    OTEL_EXPORTER_OTLP_PROTOCOL: "grpc"
    OTEL_EXPORTER_OTLP_ENDPOINT: "http://localhost:4317"
    # ... other OTel settings
  enabled: true
  timeout: 10000
  connectionMode: lenient
```

## Available Tools

### 1. validate_tool_selection
Validates routing decisions and provides analytics for tool selection optimization.

**Parameters:**
- `task_description` (string): Description of the task to be performed
- `selected_tools` (array): List of tools selected for the task
- `routing_context` (object, optional): Additional routing context

**Example:**
```json
{
  "task_description": "Analyze Python code for security vulnerabilities",
  "selected_tools": ["code-index", "filesystem"],
  "routing_context": {
    "domain": "security",
    "complexity": "high"
  }
}
```

### 2. track_tool_execution
Tracks tool execution metrics for performance analysis.

**Parameters:**
- `tool_name` (string): Name of the executed tool
- `execution_time_ms` (number): Execution time in milliseconds
- `success` (boolean): Whether the tool execution was successful
- `error_message` (string, optional): Error message if execution failed

**Example:**
```json
{
  "tool_name": "code-index",
  "execution_time_ms": 1250,
  "success": true
}
```

### 3. initialize_session
Initializes a new routing analytics session.

**Parameters:**
- `session_id` (string): Unique session identifier
- `user_context` (object, optional): User context information

**Example:**
```json
{
  "session_id": "session-12345",
  "user_context": {
    "user_type": "developer",
    "project": "security-audit"
  }
}
```

### 4. get_routing_analytics
Retrieves routing analytics and performance metrics.

**Parameters:**
- `time_range` (string): Time range for analytics (e.g., "1h", "24h", "7d")
- `metrics` (array, optional): Specific metrics to retrieve
- `filters` (object, optional): Filters to apply to analytics

**Example:**
```json
{
  "time_range": "24h",
  "metrics": ["success_rate", "avg_execution_time"],
  "filters": {
    "tool_domain": "security"
  }
}
```

## Integration with Cipher-Aggregator

### Connection Status
The routing-metadata server connects automatically when cipher-aggregator starts:
- **Successful Connection**: Log shows "Successfully connected to routing-metadata"
- **Health Monitoring**: Connection status available via aggregator health endpoint
- **Error Handling**: Graceful fallback with "lenient" connection mode

### Request Flow
1. Agent makes tool request → cipher-aggregator
2. Aggregator consults routing-metadata for validation → stdio RPC
3. Routing metadata provides analytics/validation → stdio response
4. Aggregator makes routing decision → tool execution
5. Execution results tracked back to routing-metadata → analytics update

### Logging Integration
- **OpenTelemetry**: Structured logging with OTel exporter
- **Fallback**: Automatic fallback to stdlib logging if OTel collector unavailable
- **Performance**: Minimal overhead with async logging

## Testing & Validation

### Local Testing
Test the server directly:
```bash
# Start server directly
/home/jrede/dev/MCP/eval_llm_venv/bin/python /home/jrede/dev/MCP/servers/routing-metadata-mcp.py

# Test tool listing
curl -X POST http://localhost:3020/sse \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc": "2.0", "id": 1, "method": "tools/list"}'
```

### Aggregator Integration Testing
1. **Check Connection**: Monitor aggregator logs for successful connection
2. **Tool Availability**: Verify tools are listed in aggregator tool catalog
3. **Request Routing**: Test routing decisions through aggregator
4. **Analytics Collection**: Confirm execution tracking works

### Troubleshooting

#### Connection Issues
- **Verify venv**: Ensure eval_llm_venv has all required dependencies
- **Check logs**: Look for OTel and MCP SDK initialization messages
- **Test standalone**: Run server directly to isolate issues

#### Performance Issues
- **Monitor timeouts**: 10-second timeout may be too short for complex analytics
- **OTel collector**: Consider running local OTel collector for better observability
- **Resource usage**: Monitor memory usage during heavy analytics operations

## Development

### Code Structure
```
servers/routing-metadata-mcp.py
├── Server initialization (MCP SDK)
├── Tool handlers (@app.list_tools, @app.call_tool)
├── OTel logging integration
├── Analytics processing logic
└── Error handling & fallback
```

### Adding New Tools
1. Define tool schema in `list_tools()` function
2. Implement handler in `call_tool()` function
3. Add appropriate logging and error handling
4. Test integration with aggregator

### Updating Dependencies
```bash
# Update eval_llm_venv
source /home/jrede/dev/MCP/eval_llm_venv/bin/activate
pip install --upgrade mcp opentelemetry-api opentelemetry-sdk
```

## Monitoring & Observability

### Key Metrics
- **Connection Success Rate**: Monitor aggregator logs for connection status
- **Tool Response Times**: Track validate_tool_selection performance
- **Error Rates**: Watch for SDK and OTel-related errors
- **Session Analytics**: Monitor routing analytics collection

### Log Locations
- **Aggregator Logs**: `/home/jrede/dev/MCP/logs/cipher-aggregator-*.log`
- **Server Logs**: Redirected to aggregator logs via stdio
- **OTel Logs**: If collector running, check collector output

## Security Considerations
- **Input Validation**: All tool parameters validated via MCP SDK
- **Session Management**: Unique session IDs prevent conflicts
- **Resource Limits**: Timeout and connection mode provide safety bounds
- **Logging**: Sensitive data filtering in OTel export

---
**Last Updated**: 2025-11-12 17:14:15
**Version**: 1.0.0
**Status**: Production Ready
