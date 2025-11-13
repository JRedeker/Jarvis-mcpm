# MCP Tool Conflict Resolution & Fallback Strategies

Based on analysis of 134 error entries in cipher system logs, this document provides robust conflict resolution mechanisms for the enhanced routing system.

## Critical Issues Identified

### 1. Process Conflicts (Major)
- **firecrawl-mcp**: Multiple processes running simultaneously
- **memory-bank**: Process conflicts and resource exhaustion
- **filesystem**: Duplicate process instances
- **Impact**: Resource exhaustion, failed tool calls, system instability

### 2. Connection Failures (High)
- **SSE Connection Timeouts**: 1-5 retry attempts before failure
- **Database Connectivity**: MySQL connection refused (127.0.0.1:3306)
- **API Rate Limiting**: Brave Search and other external APIs

### 3. Tool Execution Timeouts (High)
- **firecrawl_search**: Regular timeout failures
- **Code Index Operations**: Slow response times
- **Memory Bank Operations**: Intermittent failures

## Conflict Resolution Framework

### Tier 1: Prevention & Monitoring

#### Process Health Monitoring
```javascript
// Monitor for conflicting processes before tool selection
function checkProcessHealth() {
  const conflictProcesses = [
    'firecrawl-mcp',
    'memory-bank-mcp',
    'filesystem-mcp'
  ];

  return conflictProcesses.map(proc => {
    const count = processCount(proc);
    return {
      name: proc,
      instances: count,
      healthy: count <= 1,
      action: count > 1 ? 'cleanup_required' : 'ready'
    };
  });
}
```

#### Connection Health Checks
```javascript
// Check SSE server health before tool calls
async function verifyConnectionHealth() {
  try {
    const response = await fetch('http://localhost:3020/health');
    return response.ok;
  } catch (error) {
    logConnectionError(error);
    return false;
  }
}
```

### Tier 2: Immediate Conflict Resolution

#### MCP Server Process Cleanup
```javascript
// Clean up conflicting processes before tool execution
async function cleanupConflictingProcesses() {
  const cleanupScript = `./mcp-manager.sh kill-all`;
  exec(cleanupScript, (error, stdout, stderr) => {
    if (error) {
      logSystemEvent('Cleanup failed', error);
    } else {
      logSystemEvent('Process cleanup completed', stdout);
    }
  });
}
```

#### Database Connection Recovery
```javascript
// Handle MySQL connection failures
async function handleDatabaseFailure() {
  // Check if MySQL service is running
  const mysqlCheck = execSync('systemctl is-active mysql', {encoding: 'utf8'});

  if (mysqlCheck.trim() !== 'active') {
    // Start MySQL service
    execSync('sudo systemctl start mysql');
    logSystemEvent('MySQL service restarted');
  }

  // Test connection
  try {
    await testDatabaseConnection();
    return true;
  } catch (error) {
    logSystemEvent('Database recovery failed', error);
    return false;
  }
}
```

### Tier 3: Tool-Specific Fallback Chains

#### firecrawl-mcp Fallback Chain
```
Primary: firecrawl-mcp tools
├─ Failure: Check process health → Restart if needed
├─ Still Failing: Use brave-search for web content discovery
└─ Last Resort: Use fetch with manual HTML parsing
```

#### memory-bank-mcp Fallback Chain
```
Primary: memory-bank-mcp tools
├─ Failure: Verify database connectivity
├─ Still Failing: Use cipher built-in memory tools
└─ Last Resort: Store locally in filesystem
```

#### github-mcp Fallback Chain
```
Primary: github-mcp tools
├─ Failure: Check GitHub API rate limits
├─ Still Failing: Use fetch with GitHub REST API
└─ Last Resort: Manual git commands via cipher_bash
```

#### filesystem-mcp Fallback Chain
```
Primary: filesystem-mcp tools
├─ Failure: Check process conflicts → Cleanup
├─ Still Failing: Use cipher_bash for basic operations
└─ Last Resort: Direct file I/O operations
```

### Tier 4: Retry Logic & Backoff Strategies

#### Exponential Backoff Implementation
```javascript
async function executeWithRetry(operation, maxRetries = 3, baseDelay = 1000) {
  for (let attempt = 1; attempt <= maxRetries; attempt++) {
    try {
      return await operation();
    } catch (error) {
      if (attempt === maxRetries) {
        throw new Error(`Operation failed after ${maxRetries} attempts: ${error.message}`);
      }

      const delay = baseDelay * Math.pow(2, attempt - 1);
      logSystemEvent(`Attempt ${attempt} failed, retrying in ${delay}ms`, error);
      await sleep(delay);
    }
  }
}
```

#### Tool-Specific Retry Policies
```javascript
const retryPolicies = {
  'firecrawl_search': { maxRetries: 2, baseDelay: 2000, exponential: true },
  'memory_bank_store': { maxRetries: 3, baseDelay: 1000, exponential: true },
  'github_create_repo': { maxRetries: 1, baseDelay: 500, exponential: false },
  'filesystem_read': { maxRetries: 2, baseDelay: 500, exponential: true }
};
```

### Tier 5: Circuit Breaker Pattern

#### Circuit Breaker Implementation
```javascript
class ToolCircuitBreaker {
  constructor(threshold = 5, timeout = 60000) {
    this.threshold = threshold;
    this.timeout = timeout;
    this.failureCount = 0;
    this.lastFailureTime = null;
    this.state = 'CLOSED'; // CLOSED, OPEN, HALF_OPEN
  }

  async execute(toolName, operation) {
    if (this.state === 'OPEN') {
      if (Date.now() - this.lastFailureTime > this.timeout) {
        this.state = 'HALF_OPEN';
      } else {
        throw new Error(`Circuit breaker is OPEN for ${toolName}`);
      }
    }

    try {
      const result = await operation();
      this.onSuccess();
      return result;
    } catch (error) {
      this.onFailure();
      throw error;
    }
  }

  onSuccess() {
    this.failureCount = 0;
    this.state = 'CLOSED';
  }

  onFailure() {
    this.failureCount++;
    this.lastFailureTime = Date.now();

    if (this.failureCount >= this.threshold) {
      this.state = 'OPEN';
      logSystemEvent(`Circuit breaker opened for ${this.failureCount} failures`);
    }
  }
}
```

## Dynamic Fallback Selection

### Task-Based Fallback Logic
```javascript
function selectFallbackTool(primaryTool, error, context) {
  const taskContext = extractTaskContext(context);

  switch (primaryTool) {
    case 'firecrawl_search':
      if (error.code === 'ETIMEDOUT') {
        return { tool: 'brave_web_search', reason: 'timeout_fallback' };
      }
      if (error.code === 'ECONNREFUSED') {
        return { tool: 'fetch', reason: 'connection_fallback' };
      }
      break;

    case 'memory_bank_store':
      if (error.code === 'ECONNREFUSED') {
        return { tool: 'filesystem_write', reason: 'db_fallback' };
      }
      break;

    case 'github_create_repo':
      if (error.code === 'RATE_LIMITED') {
        return { tool: 'cipher_bash', reason: 'rate_limit_fallback' };
      }
      break;
  }

  return { tool: 'generic_error_handler', reason: 'unknown_error' };
}
```

### Performance-Based Fallback
```javascript
// Monitor tool performance and auto-fallback to alternatives
class PerformanceMonitor {
  constructor() {
    this.performanceMetrics = new Map();
  }

  recordToolPerformance(toolName, executionTime, success) {
    const metrics = this.performanceMetrics.get(toolName) || {
      totalCalls: 0,
      successCalls: 0,
      totalTime: 0,
      averageTime: 0
    };

    metrics.totalCalls++;
    metrics.totalTime += executionTime;

    if (success) {
      metrics.successCalls++;
    }

    metrics.averageTime = metrics.totalTime / metrics.totalCalls;
    metrics.successRate = metrics.successCalls / metrics.totalCalls;

    this.performanceMetrics.set(toolName, metrics);
  }

  shouldUseFallback(toolName) {
    const metrics = this.performanceMetrics.get(toolName);
    if (!metrics) return false;

    return metrics.successRate < 0.7 || metrics.averageTime > 10000;
  }
}
```

## Implementation Guidelines

### 1. Pre-Execution Checks
- Verify process health for critical MCP servers
- Check database connectivity before memory operations
- Validate API rate limits before external calls

### 2. Graceful Degradation
- Implement progressive fallbacks (Specialized → Generic → Manual)
- Use circuit breakers to prevent cascading failures
- Provide meaningful error messages for human intervention

### 3. Recovery Mechanisms
- Automatic process cleanup for conflicts
- Database service restart for connection failures
- API key validation and refresh for authentication errors

### 4. Monitoring & Alerting
- Track success/failure rates for all tools
- Monitor average response times and trend analysis
- Alert on circuit breaker openings and repeated failures

### 5. Configuration Management
- Tool-specific retry policies based on criticality
- Performance thresholds for automatic fallback activation
- Resource limits for process management

## Success Metrics

- **90%+ tool selection accuracy** for optimal tools
- **50% reduction in average task completion time**
- **70% reduction in tool-related failures**
- **Zero process conflicts** through proactive monitoring
- **<5% circuit breaker activation rate**

This framework ensures robust operation of the enhanced routing system despite the conflicts and failures identified in the system logs.