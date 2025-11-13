# Alert Creation Workflow

## Overview

Alerts are system-detected issues that should be automatically identified and reported when patterns of failure or problems are detected. This document defines when and how to create alerts.

## Alert Detection Triggers

### When to Create an Alert

Create an alert when you detect:

1. **Repeated Tool Failures** (3+ failures of same operation)
   - MCP tool timeouts or connection errors
   - Repeated JSON-RPC errors
   - Tool not found errors (especially after recent changes)

2. **Configuration Issues**
   - Missing required environment variables
   - Invalid configuration values detected
   - Configuration drift (expected vs actual state)

3. **Dependency Problems**
   - Missing Python packages
   - Node modules not installed
   - Binary tools not found in PATH

4. **Performance Degradation**
   - Response times exceeding thresholds (>5s consistently)
   - Memory usage patterns indicating leaks
   - Disk space warnings

5. **Integration Failures**
   - External API repeatedly unavailable
   - Database connection issues
   - File system permission errors

6. **Setup/Tooling Issues**
   - Server startup failures
   - Port conflicts
   - Version mismatches

## Alert Notification Pattern

### Step 1: Detect and Log

When you detect an issue that warrants an alert:

```markdown
⚠️ **ALERT DETECTED**: [Brief issue description]

**Pattern**: [What triggered this - e.g., "MCP tool 'firecrawl_scrape' failed 3 times in a row"]
**Impact**: [How this affects the system - e.g., "Web scraping operations currently unavailable"]
**Severity**: [critical/warning/info]
**Recommendation**: [What should be done - e.g., "Check MCP server status and configuration"]

Would you like me to create an alert in `./open/alerts/` for tracking and resolution?
```

### Step 2: Get User Approval

Present the user with:
- Clear description of the issue
- Severity level
- Impact assessment
- Recommended next steps
- Offer to create formal alert

### Step 3: Create Alert Structure

If user approves, create:

```
./open/alerts/{severity}-{issue-name}/
├── {severity}-{issue-name}.md     # Main alert document
├── detection-log.md                # When/how/where detected
└── context/                        # Evidence
    ├── error-logs.txt
    ├── screenshots/ (if applicable)
    └── stack-traces.txt
```

## Alert Document Template

### Main Alert File: `{severity}-{issue-name}.md`

```markdown
# {Severity}: {Issue Name}

**Status**: Open
**Created**: {Date/Time}
**Severity**: {critical/warning/info}
**Detected By**: System
**Impact**: {High/Medium/Low}

## Summary

[1-2 sentence description of the issue]

## Observed Behavior

[What is happening - be specific]

## Expected Behavior

[What should be happening instead]

## Impact Assessment

**User Impact**: [How does this affect end users?]
**System Impact**: [How does this affect system operations?]
**Workaround**: [Is there a temporary workaround?]

## Detection Details

**First Detected**: {timestamp}
**Occurrence Count**: {number}
**Pattern**: {what pattern triggered the alert}

See: `./detection-log.md` for full timeline

## Evidence

See: `./context/` directory for:
- Error logs
- Stack traces
- Configuration snapshots
- Screenshots (if applicable)

## Recommended Actions

1. {Action item 1}
2. {Action item 2}
3. {Action item 3}

## Root Cause Analysis

[To be filled in during investigation]

## Resolution

[To be filled in when resolved]

## Prevention

[To be filled in - how to prevent this in future]

## Related

- **Tickets**: [Links to related tickets if any]
- **Specs**: [Links to specs that might resolve this]
```

### Detection Log File: `detection-log.md`

```markdown
# Detection Log: {Issue Name}

## Timeline

### {Timestamp 1}
- **Event**: {What happened}
- **Context**: {Where/when/how}
- **Error**: {Error message if applicable}

### {Timestamp 2}
- **Event**: {What happened}
- **Context**: {Where/when/how}
- **Error**: {Error message if applicable}

### {Timestamp 3} - Alert Triggered
- **Event**: {What happened}
- **Threshold**: {What threshold was crossed}
- **Decision**: Alert created

## Pattern Analysis

**Frequency**: {How often is this occurring}
**Consistency**: {Is it consistent or sporadic}
**Correlation**: {Any correlation with other events}
```

## Severity Levels

### Critical
- Prefix: `critical-`
- Criteria: System down, blocking all work
- Examples:
  - MCP aggregator won't start
  - Database unreachable
  - Core dependency missing

### Warning
- Prefix: `warning-`
- Criteria: Degraded performance, needs attention soon
- Examples:
  - Repeated tool timeouts
  - Performance degradation
  - Partial functionality unavailable

### Info
- Prefix: `info-`
- Criteria: Configuration drift, non-urgent
- Examples:
  - Version mismatch detected
  - Deprecated API usage
  - Optimization opportunity identified

## Example Alert Notifications

### Example 1: MCP Tool Failures

```markdown
⚠️ **ALERT DETECTED**: MCP Tool Repeated Failures

**Pattern**: Tool 'firecrawl_scrape' failed 3 consecutive times with timeout errors
**Impact**: Web scraping operations currently unavailable
**Severity**: warning
**Last Error**: "Request timed out after 30000ms"

**Recommendation**:
1. Check if Firecrawl MCP server is running
2. Verify network connectivity
3. Check Firecrawl API rate limits

Would you like me to create an alert in `./open/alerts/warning-firecrawl-timeout/`
for proper tracking and resolution?
```

### Example 2: Configuration Issue

```markdown
⚠️ **ALERT DETECTED**: Missing Environment Variable

**Pattern**: OPENROUTER_API_KEY not set in environment
**Impact**: LLM inference requests will fail
**Severity**: critical
**Detection**: Attempted to call llm_inference_auto but API key missing

**Recommendation**:
1. Set OPENROUTER_API_KEY in .env file
2. Restart MCP servers to pick up new environment

Would you like me to create an alert in `./open/alerts/critical-missing-api-key/`
for tracking and resolution?
```

### Example 3: Performance Degradation

```markdown
⚠️ **ALERT DETECTED**: Response Time Degradation

**Pattern**: Average response time for code search increased from 200ms to 5000ms
**Impact**: Development workflow significantly slowed
**Severity**: warning
**Timeframe**: Started approximately 2 hours ago

**Recommendation**:
1. Check system resource usage (CPU, memory, disk)
2. Review recent index rebuilds
3. Check for file system issues

Would you like me to create an alert in `./open/alerts/warning-code-search-slow/`
for investigation?
```

## Alert Lifecycle

1. **Detection**: System detects pattern requiring alert
2. **Notification**: Present to user with clear details
3. **Creation**: User approves, alert created in `./open/alerts/`
4. **Investigation**: Root cause analysis added to alert
5. **Resolution**:
   - If requires spec: Move to `./open/specs/{spec-name}/alerts/`
   - If fixed directly: Document resolution in alert
   - When complete: Delete from `./open/` (git preserves history)

## Integration with Specs

When an alert requires architectural changes:

1. Create or identify relevant spec
2. Move alert: `./open/alerts/{name}/` → `./open/specs/{spec-name}/alerts/{name}/`
3. Link alert to spec in both directions
4. Track resolution as part of spec implementation

## Best Practices

### Detection Thresholds
- Don't alert on first failure (could be transient)
- Use 3 consecutive failures or 5 in 10 minutes as threshold
- Consider severity when setting thresholds

### Alert Content
- Be specific about what's broken
- Include error messages verbatim
- Provide actionable recommendations
- Document evidence thoroughly

### User Communication
- Keep notification concise but informative
- Present severity upfront
- Offer to create alert (don't create automatically)
- Provide immediate workarounds if available

### Avoiding Alert Fatigue
- Don't create duplicate alerts for same issue
- Group related failures into single alert
- Update existing alerts rather than creating new ones
- Mark alerts as resolved when fixed

## Created: 2025-11-12
## Related: open-directory-structure.md, ticket-workflow-preferences.md
