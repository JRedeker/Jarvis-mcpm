# Slash Command System - Specification

**Status**: Proposed
**Priority**: Medium
**Owner**: Platform/Agents

## Overview

This specification defines a centralized slash command system for AI agents (Cline, Kilocode) that enables structured workflows, tier signaling, and workflow reuse through a registry-based approach.

## Problem Statement

Currently, users must rely on natural language to trigger specific workflows or tier preferences. This leads to:
- Inconsistent tier selection (can't force m4 for deep reasoning)
- No standardized way to trigger common workflows (ticket creation, research, updates)
- Repeated workflow descriptions in prompts
- Difficulty enforcing Plan mode behavior programmatically

## Solution

Implement a hybrid slash command system that:
1. **Recognizes structured commands** like `/ticket`, `/research`, `/m4`
2. **Signals tier preferences** to the LLM routing system
3. **Triggers predefined workflows** stored in memory-bank
4. **Supports both direct matching and fuzzy/LLM inference**

## Directory Structure

```
open/specs/slash-command-system/
├── slash-command-system.md  # This file - spec index & overview
└── spec.md                   # Detailed specification with comparison table
```

## Main Specification

**File**: `./spec.md`

Contains:
- Detailed comparison table (3 implementation options)
- Architecture and design
- Command registry structure
- Integration points
- Implementation plan

## Use Cases

### 1. Tier Signaling
```
User: "/m4 analyze this complex algorithm"
→ Forces m4 tier (claude-sonnet-4.5)
→ Ensures maximum reasoning capability
```

### 2. Workflow Triggers
```
User: "/ticket bug in search functionality"
→ Loads ticket creation workflow
→ Creates ./open/tickets/{name}/{name}.md
→ Uses template from memory-bank
```

### 3. Research Mode
```
User: "/research latest trends in MCP servers"
→ Forces m3 tier for reasoning
→ Triggers research workflow (brave_search + firecrawl + synthesis)
→ Stores findings in memory-bank
```

### 4. Update Detection
```
User: "/up ./open/specs/vscode-lm-provider-cipher/spec.md"
→ Scans spec for outdated information
→ Searches for latest details
→ Proposes updates with replace_in_file
```

### 5. Plan Mode Forcing
```
User: "/plan how should we implement X?"
→ Forces Plan mode behavior
→ Explores with read_file, list_files
→ Presents architectural plan before acting
```

## Key Benefits

1. **Predictable Behavior**: Same command → same workflow
2. **Tier Control**: Explicit tier selection when needed
3. **Workflow Reuse**: Centralized workflow definitions
4. **Fast Execution**: Direct matching for known commands
5. **Flexible Parsing**: Fuzzy matching + LLM fallback for variations
6. **Self-Documenting**: Commands describe their purpose and parameters

## Implementation Status

- **Status**: Proposed - awaiting approval
- **Next Steps**: Review comparison table, select approach, implement

## Quick Links

- [Detailed Specification](./spec.md) - Full spec with comparison table
- [Memory Bank](../../../data/memory-bank/development-workflow/) - Workflow storage

## Related Work

- `data/memory-bank/development-workflow/` - Existing workflow documentation
- `.clinerules/` - Cline integration point
- `.kilocoderules` - Kilocode integration point
- `servers/llm-inference-mcp.py` - Tier routing system

## Decision Points

1. **Implementation Approach**: Direct mapping, LLM inference, or hybrid?
2. **Storage Format**: YAML, JSON, or code-based registry?
3. **Fuzzy Matching**: Levenshtein distance threshold?
4. **LLM Fallback**: Always enabled or opt-in?
5. **Command Namespacing**: Flat or hierarchical? (e.g., `/git/commit` vs `/gitcommit`)

## Contact & Ownership

**Owners**: Platform/Agents team
**Status**: Proposed - review comparison table in spec.md
**Next Steps**: Select implementation approach and create PoC
