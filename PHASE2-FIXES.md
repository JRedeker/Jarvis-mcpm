# PHASE 2: MCP Server Fixes - Action Plan

**Date:** 2025-11-15
**Status:** In Progress

---

## 🎉 BREAKTHROUGH DISCOVERY

**Cipher IS loading mcpServers correctly!**

Log evidence:
```
15:29:41 INFO: MCP Manager: Initializing 12 servers (2 disabled)
```

The aggregator mode is working. The issue is **individual server connection failures**.

---

## 📊 Current Server Status

### ✅ Working Servers (4/12)
1. **brave-search** - Web search working
2. **filesystem** - File operations working
3. **file-batch** - Batch file operations working
4. **schemathesis** - API schema testing working

### ❌ Failed Servers (8/12)

#### 1. routing-metadata
**Error:** `ModuleNotFoundError: No module named 'servers.log_setup'`
**Fix:** Create missing log_setup module or remove dependency

#### 2. httpie
**Error:** Connection timeout (30s)
**Fix:** Check if server script has errors, increase timeout

#### 3. pytest
**Error:** Connection timeout (30s)
**Fix:** Check if server script has errors, increase timeout

#### 4. firecrawl
**Error:** `Either FIRECRAWL_API_KEY or FIRECRAWL_API_URL must be provided`
**Fix:** Add real API key or disable server

#### 5. memory-bank
**Error:** `npm error 404 Not Found - @modelcontextprotocol/memory-bank-mcp-server`
**Fix:** Use correct package name: `@modelcontextprotocol/server-memory`

#### 6. gpt-researcher
**Error:** `npm error 404 Not Found - gptr-mcp`
**Fix:** Use correct package name: `@gptr/mcp-server` or similar

#### 7. context7
**Error:** Connection failed to `https://mcp.context7.com/mcp`
**Fix:** Verify URL, check if service is available

#### 8. morph-fast-apply
**Error:** Connection failed to `https://mcp.morph.ai/fast-apply`
**Fix:** Verify URL, check if service requires authentication

---

## 🔧 Immediate Actions

### Priority 1: Fix Python Servers (routing-metadata, httpie, pytest)

**Issue:** Missing `servers/log_setup.py` module

**Solution A:** Create the missing module
**Solution B:** Remove the import from servers
**Solution C:** Disable these servers temporarily

### Priority 2: Fix NPM Package Names

**memory-bank:**
```yaml
# WRONG
command: npx
args: ["-y", "@modelcontextprotocol/memory-bank-mcp-server"]

# CORRECT (need to verify actual package name)
command: npx
args: ["-y", "@modelcontextprotocol/server-memory"]
```

**gpt-researcher:**
```yaml
# WRONG
command: npx
args: ["-y", "-p", "gptr-mcp@latest", "gptr-mcp"]

# CORRECT (need to verify actual package name)
command: npx
args: ["-y", "@gptr/mcp-server"]
```

### Priority 3: Add Missing API Keys

Add to `.env`:
```bash
FIRECRAWL_API_KEY=your-firecrawl-api-key
# OR disable firecrawl in cipher.yml
```

### Priority 4: Verify Remote URLs

Test remote endpoints:
```bash
curl -v https://mcp.context7.com/mcp
curl -v https://mcp.morph.ai/fast-apply
```

---

## 📋 Implementation Steps

### Step 1: Create log_setup.py
```python
# servers/log_setup.py
import logging
import sys

def init_logging(name="mcp-server", level=logging.INFO):
    """Initialize logging for MCP servers"""
    logger = logging.getLogger(name)
    logger.setLevel(level)

    handler = logging.StreamHandler(sys.stderr)
    handler.setLevel(level)

    formatter = logging.Formatter(
        '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
    )
    handler.setFormatter(formatter)

    logger.addHandler(handler)
    return logger
```

### Step 2: Research Correct Package Names

Need to verify:
- [ ] memory-bank correct npm package
- [ ] gpt-researcher correct npm package
- [ ] context7 endpoint status
- [ ] morph-fast-apply endpoint status

### Step 3: Update cipher.yml

Create corrected version with:
- Fixed package names
- Increased timeouts for slow servers
- Proper error handling
- Optional servers marked as `enabled: false`

### Step 4: Test Each Server Individually

```bash
# Test memory-bank
npx -y @modelcontextprotocol/server-memory

# Test gpt-researcher
npx -y @gptr/mcp-server

# Test context7
curl https://mcp.context7.com/mcp
```

---

## 🎯 Success Criteria

- [ ] All Python servers (routing-metadata, httpie, pytest) connect successfully
- [ ] memory-bank connects with correct package name
- [ ] gpt-researcher connects with correct package name OR disabled if unavailable
- [ ] firecrawl connects with API key OR disabled
- [ ] context7 connects OR disabled if service unavailable
- [ ] morph-fast-apply connects OR disabled if service unavailable
- [ ] Cipher logs show: "Successfully connected to server: X" for all enabled servers
- [ ] Tool count increases from 13 to 30+ tools

---

## 📈 Expected Outcome

**Before:** 13 tools (Brave + Cipher memory)
**After:** 30+ tools (All working MCP servers)

**Working servers:**
- brave-search ✅
- filesystem ✅
- file-batch ✅
- schemathesis ✅
- routing-metadata (after fix)
- httpie (after fix)
- pytest (after fix)
- memory-bank (after package name fix)
- gpt-researcher (after package name fix)
- firecrawl (with API key)
- context7 (if available)
- morph-fast-apply (if available)

---

## 🚀 Next Steps

1. Create `servers/log_setup.py`
2. Research correct npm package names
3. Update cipher.yml with fixes
4. Restart Cipher aggregator
5. Verify all tools available
6. Document working configuration
