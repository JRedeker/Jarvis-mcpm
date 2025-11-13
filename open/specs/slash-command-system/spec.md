# Slash Command System - Detailed Specification

## 1. Executive Summary

This specification proposes a centralized slash command system for AI agents that enables:
- Structured command recognition (`/ticket`, `/m4`, `/research`)
- Explicit tier signaling to LLM routing
- Reusable workflow definitions
- Flexible parsing (direct matching + fuzzy + LLM inference)

## 2. Implementation Options Comparison

### Comparison Table

| Dimension | Option 1: Direct Mapping | Option 2: LLM Inference | Option 3: Hybrid (Recommended) |
|-----------|-------------------------|------------------------|-------------------------------|
| **Performance** | ⚡⚡⚡ <1ms | 🐌 200-500ms | ⚡⚡⚡ <1ms (direct)<br>⚡ 10ms (fuzzy)<br>🐌 200-500ms (LLM fallback) |
| **Latency Detail** | O(1) hash lookup | Network + inference time | Cascading: direct → fuzzy → LLM |
| **Cost** | 💰 $0.00 (free) | 💸 ~$0.0001 per command | 💰 ~$0.00001 average<br>(LLM rarely used) |
| **Token Usage** | 0 tokens | ~100-200 tokens/command | ~5-20 tokens average |
| **Flexibility** | ❌ Exact match only<br>No typo tolerance | ✅ Full natural language<br>✅ Variations understood | ✅✅ Best of both:<br>- Exact matching<br>- Typo correction<br>- Natural language fallback |
| **Example Flexibility** | `/ticket` ✅<br>`/tickt` ❌<br>`make ticket` ❌ | `/ticket` ✅<br>`/tickt` ✅<br>`make ticket` ✅ | `/ticket` ✅ (instant)<br>`/tickt` ✅ (fuzzy)<br>`make ticket` ✅ (LLM) |
| **Determinism** | ✅✅✅ 100% deterministic | ❌ Non-deterministic<br>May vary per request | ✅✅ 95% deterministic<br>(LLM only for unknowns) |
| **Accuracy** | 100% (for exact matches)<br>0% (for variations) | ~90-95% (depends on LLM) | ~98% (combined strength) |
| **Maintainability** | ✅✅ Easy:<br>Edit YAML file | ⚠️ Complex:<br>Tune prompts, manage memory | ✅ Moderate:<br>YAML + fuzzy threshold |
| **Extensibility** | ✅ Add to YAML | ✅ Learns from memory<br>(but less explicit) | ✅✅ Add to YAML<br>+ automatic learning |
| **Implementation** | Simple:<br>- YAML parser<br>- Dict lookup | Complex:<br>- LLM integration<br>- Memory management<br>- Prompt engineering | Moderate:<br>- YAML parser<br>- Fuzzy matcher<br>- LLM fallback |
| **Dependencies** | None (stdlib only) | - evalLLM/OpenRouter<br>- memory-bank | - Fuzzy lib (fuzzywuzzy)<br>- evalLLM (optional) |
| **Error Handling** | Easy: "Command not found" | Complex: Ambiguous intents | Clear: "Did you mean X?" |
| **User Experience** | ⚠️ Frustrating for typos | ✅ Natural and forgiving | ✅✅ Best UX:<br>Fast + forgiving |
| **Debugging** | ✅✅ Easy to trace | ❌ Opaque (LLM black box) | ✅ Clear cascade path |
| **Testability** | ✅✅ Unit testable | ⚠️ Requires LLM mocking | ✅ Each layer testable |
| **Offline Support** | ✅ Works offline | ❌ Requires network | ⚠️ Degrades gracefully |
| **Latency p50** | <1ms | 250ms | 2ms (95% direct/fuzzy) |
| **Latency p95** | <1ms | 450ms | 15ms |
| **Latency p99** | <1ms | 600ms | 300ms (LLM fallback) |
| **Resource Usage** | Negligible RAM | ~100MB (model cache) | ~5MB (fuzzy index) |
| **Scalability** | ✅✅ Scales to 1000s | ⚠️ Rate limited by API | ✅✅ Scales well |
| **Learning Curve** | Low (users learn commands) | Very low (natural language) | Low (works both ways) |
| **Documentation** | Self-documenting YAML | Needs examples in memory | Self-documenting + learned |
| **Version Control** | ✅ Git-friendly (YAML) | ⚠️ Memory in DB/JSON | ✅ YAML in git |
| **Multi-language** | ✅ Easy (add translations) | ✅✅ Built-in understanding | ✅✅ Best of both |
| **Command Discovery** | ❌ Must know commands | ✅ Can guess intent | ✅ Suggestions on mismatch |
| **Consistency** | ✅✅ Always consistent | ❌ May change over time | ✅ Mostly consistent |

### Visual Performance Comparison

```
Latency Distribution (p50):

Option 1: ▓ <1ms
Option 2: ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ 250ms
Option 3: ▓ 2ms

Cost per 1000 commands:

Option 1: $0.00
Option 2: $0.10
Option 3: $0.01
```

## 3. Detailed Architecture (Hybrid Approach)

### 3.1 Data Structure

**File**: `data/slash-commands.yml`

```yaml
version: "1.0"
metadata:
  last_updated: "2025-11-12"
  command_count: 15

# Core command definitions
commands:
  ticket:
    description: "Create a new ticket in ./open/tickets/"
    tier: null  # Use auto-routing
    mode: null  # Don't force mode
    parameters:
      - name: "title"
        required: false
        description: "Ticket title (will prompt if not provided)"
    workflow_ref: "data/memory-bank/workflows/ticket-workflow.md"
    template_ref: "data/memory-bank/templates/ticket-template.md"
    examples:
      - "/ticket fix search bug"
      - "/ticket"

  research:
    description: "Deep research and evaluation of topic"
    tier: "m3"  # Force reasoning tier
    mode: null
    parameters:
      - name: "topic"
        required: true
        description: "Research topic"
    workflow_ref: "data/memory-bank/workflows/research-workflow.md"
    tools:
      - "brave_web_search"
      - "firecrawl_scrape"
      - "memory_bank_write"
    examples:
      - "/research latest MCP server patterns"
      - "/research OpenRouter model pricing"

  up:
    description: "Update outdated information in a document"
    tier: "m2"  # Balanced tier
    mode: null
    parameters:
      - name: "context"
        required: true
        description: "File path or context to update"
    workflow_ref: "data/memory-bank/workflows/update-workflow.md"
    tools:
      - "read_file"
      - "brave_web_search"
      - "replace_in_file"
    examples:
      - "/up ./open/specs/vscode-lm-provider-cipher/spec.md"
      - "/up OpenRouter model list"

  plan:
    description: "Force Plan mode behavior"
    tier: "m1"  # Light tier for planning
    mode: "plan"  # Force Plan mode
    workflow_ref: "data/memory-bank/workflows/plan-mode-workflow.md"
    behavior:
      - "Use read_file and list_files to explore"
      - "Ask clarifying questions"
      - "Present architectural plan"
      - "Wait for approval before acting"
    examples:
      - "/plan how should we implement X?"
      - "/plan"

  m4:
    description: "Force maximum reasoning tier (m4)"
    tier: "m4"
    mode: null
    workflow: "Continue with user's request using m4 tier"
    model_info:
      tier: "m4"
      model: "claude-sonnet-4.5-20241022"
      context: "1M tokens"
      output: "48K tokens"
    examples:
      - "/m4 analyze this algorithm"
      - "/m4 deep dive into architecture"

  deepthink:
    description: "Deep reasoning mode (alias for /m4)"
    tier: "m4"
    alias_of: "m4"
    workflow_ref: "data/memory-bank/workflows/deep-reasoning-workflow.md"
    examples:
      - "/deepthink complex problem"

  alert:
    description: "Create a new alert in ./open/alerts/"
    tier: "m1"
    mode: null
    parameters:
      - name: "severity"
        required: true
        description: "critical, warning, or info"
      - name: "issue"
        required: true
        description: "Issue name"
    workflow_ref: "data/memory-bank/workflows/alert-workflow.md"
    examples:
      - "/alert warning mcp-timeout"
      - "/alert critical missing-api-key"

  spec:
    description: "Create a new spec in ./open/specs/"
    tier: "m2"
    mode: "plan"  # Force planning for specs
    workflow_ref: "data/memory-bank/workflows/spec-workflow.md"
    examples:
      - "/spec new-feature-name"

# Command aliases for variations
aliases:
  # Tier aliases
  dt: "deepthink"
  deep: "deepthink"
  max: "m4"
  fast: "l0"
  light: "m1"
  balanced: "m2"
  reason: "m3"

  # Workflow aliases
  newticket: "ticket"
  createticket: "ticket"
  bug: "ticket"
  feature: "ticket"

  update: "up"
  refresh: "up"

  planning: "plan"
  explore: "plan"

  investigate: "research"
  study: "research"

# Fuzzy matching configuration
fuzzy_matching:
  enabled: true
  algorithm: "levenshtein"
  threshold: 0.8  # 80% similarity required
  max_distance: 2  # Maximum edit distance
  case_sensitive: false
  suggest_on_mismatch: true

# LLM fallback configuration
llm_fallback:
  enabled: true
  use_evalLLM: true
  confidence_threshold: 0.7
  max_attempts: 1
  prompt_template: |
    User typed: "{input}"
    Available commands: {command_list}
    Known aliases: {alias_list}

    What command did the user intend? Respond with just the command name or "UNKNOWN".

# Tier signal mappings
tier_signals:
  l0: ["fast", "quick", "l0"]
  m1: ["plan", "planning", "light", "m1"]
  m2: ["up", "update", "balanced", "m2"]
  m3: ["research", "investigate", "reason", "m3"]
  m4: ["m4", "deepthink", "deep", "max", "dt"]
```

### 3.2 Processing Pipeline

```
User Input: "/tickt bug in search"
       ↓
┌─────────────────────────┐
│  1. Direct Match        │
│  Check: commands[input] │
│  Result: NOT FOUND      │
└─────────────────────────┘
       ↓
┌─────────────────────────┐
│  2. Alias Match         │
│  Check: aliases[input]  │
│  Result: NOT FOUND      │
└─────────────────────────┘
       ↓
┌─────────────────────────┐
│  3. Fuzzy Match         │
│  Calculate similarity   │
│  Result: "ticket" (90%) │ ← MATCH!
└─────────────────────────┘
       ↓
┌─────────────────────────┐
│  4. Load Workflow       │
│  Read: ticket-workflow  │
│  Execute steps          │
└─────────────────────────┘
```

### 3.3 Integration Points

#### A. .clinerules Integration

```markdown
# file: .clinerules/slash-commands.md

## Slash Command Processing

When a message starts with `/`, process as a slash command:

### Step 1: Parse Command
Extract command and parameters:
- `/ticket fix bug` → command="ticket", params=["fix", "bug"]
- `/m4 analyze` → command="m4", params=["analyze"]

### Step 2: Load Command Definition
Read `data/slash-commands.yml` and find command entry.

### Step 3: Apply Tier Signal
If command specifies `tier`, signal to llm_inference_auto:
- Set override_tier in tool call
- Example: `/m4` → override_tier="m4"

### Step 4: Apply Mode
If command specifies `mode`:
- "plan" → Use Plan mode behavior (explore, ask questions, present plan)
- null → Continue in current mode

### Step 5: Execute Workflow
Load workflow from workflow_ref and follow steps.

### Step 6: Use Template
If template_ref specified, load template for file creation.

### Common Commands
Run `/help` to list all available slash commands.
```

#### B. .kilocoderules Integration

Similar structure adapted for Kilocode's architecture.

### 3.4 Workflow Storage

**Location**: `data/memory-bank/workflows/`

**Example**: `ticket-workflow.md`

```markdown
# Ticket Creation Workflow

Triggered by: `/ticket`
Tier: auto (null)
Mode: current

## Steps

1. **Parse Input**
   - Check if title provided as parameter
   - If not: Ask "What's the ticket about? (brief title)"

2. **Generate Ticket Name**
   - Convert title to kebab-case
  - Example: "Fix Search Bug" → "fix-search-bug"

3. **Create Directory**
   ```bash
   mkdir -p ./open/tickets/{ticket-name}/
   ```

4. **Load Template**
   - Read: `data/memory-bank/templates/ticket-template.md`
   - Replace placeholders:
     - {TITLE} → user's title
     - {DATE} → current date
     - {STATUS} → "Open"

5. **Write File**
   ```bash
   ./open/tickets/{ticket-name}/{ticket-name}.md
   ```

6. **Confirm**
   - Show created ticket path
   - Ask if user wants to add notes or context

## Template Variables

- {TITLE} - Ticket title
- {DATE} - Creation date
- {STATUS} - Open/Closed
- {PRIORITY} - High/Medium/Low (prompt if needed)
```

### 3.5 Command Discovery

**List all commands**: `/help` or `/commands`

**Output**:
```
Available Slash Commands:

Tier Control:
  /m4, /deepthink - Maximum reasoning tier (claude-sonnet-4.5)
  /m3             - Deep reasoning tier
  /m2, /balanced  - Balanced tier
  /m1, /light     - Light tier
  /l0, /fast      - Speed tier

Workflows:
  /ticket       - Create new ticket
  /spec         - Create new spec
  /alert        - Create new alert
  /research     - Deep research mode
  /up <file>    - Update outdated information
  /plan         - Force Plan mode

Type /help <command> for details on a specific command.
Example: /help ticket
```

## 4. Implementation Plan

### Phase 1: Foundation (Week 1)
- [ ] Create `data/slash-commands.yml`
- [ ] Implement YAML parser
- [ ] Add direct matching logic
- [ ] Test with basic commands (/m4, /ticket)

### Phase 2: Enhanced Matching (Week 2)
- [ ] Add alias support
- [ ] Implement fuzzy matching (Levenshtein)
- [ ] Add suggestion system ("Did you mean X?")
- [ ] Test with typos and variations

### Phase 3: LLM Fallback (Week 3)
- [ ] Integrate evalLLM for unknown commands
- [ ] Add confidence scoring
- [ ] Implement learning from corrections
- [ ] Test with natural language inputs

### Phase 4: Workflow Integration (Week 4)
- [ ] Create workflow templates
- [ ] Implement workflow execution engine
- [ ] Add parameter parsing
- [ ] Test end-to-end workflows

### Phase 5: Polish (Week 5)
- [ ] Add `/help` command
- [ ] Create documentation
- [ ] Add telemetry (command usage tracking)
- [ ] Performance optimization

## 5. Testing Strategy

### Unit Tests
```python
def test_direct_match():
    assert parse_command("/ticket") == ("ticket", [])

def test_fuzzy_match():
    assert fuzzy_match("/tickt", commands) == "ticket"

def test_tier_signal():
    cmd = get_command("/m4")
    assert cmd.tier == "m4"
```

### Integration Tests
- Test full pipeline: input → command → workflow → output
- Test tier signaling to llm_inference_auto
- Test mode switching (Plan mode)

### Performance Tests
- Measure latency for 1000 direct matches (target: <1ms p99)
- Measure latency for 1000 fuzzy matches (target: <10ms p99)
- Measure LLM fallback latency (target: <500ms p95)

## 6. Success Metrics

### Performance
- **Direct match**: <1ms p99 ✅
- **Fuzzy match**: <10ms p99 ✅
- **LLM fallback**: <500ms p95 ✅

### Accuracy
- **Direct match**: 100% ✅
- **Fuzzy match**: >95% ✅
- **LLM fallback**: >90% ✅

### Cost
- **Average cost/command**: <$0.00001 ✅
- **99th percentile**: <$0.0001 ✅

### Adoption
- **Daily active commands**: >100
- **User satisfaction**: >4.5/5
- **Command coverage**: >80% of common workflows

## 7. Future Enhancements

### Command Composition
```
/m4 /research quantum computing
```
Combines m4 tier with research workflow.

### Parameter Templates
```
/ticket --priority high --assignee @john
```
Structured parameter parsing.

### Command History
```
/history
/repeat last
```
Track and repeat previous commands.

### Custom Commands
```yaml
# user-commands.yml (user-specific)
my_commands:
  deploy:
    description: "Deploy to staging"
    workflow: "run_tests && deploy_staging"
```

### Natural Language Expansion
```
User: "make me a ticket for this bug"
→ Recognized as: /ticket
```

## 8. Decision Log

### Why Hybrid Over Pure LLM?
- **Performance**: 250ms → <10ms for 95% of commands
- **Cost**: $0.10/1000 → $0.01/1000 (10x savings)
- **Determinism**: Critical for tier signaling
- **Offline**: Works without network for known commands

### Why YAML Over JSON?
- More human-readable
- Supports comments
- Better for version control diffs
- Easy to edit manually

### Why Not Code-Based Registry?
- YAML more accessible to non-developers
- Can be edited without development environment
- Clear separation of config vs. logic
- Easier to version and review

## 9. Open Questions

1. **Command Namespacing**: Flat (`/gitcommit`) vs hierarchical (`/git/commit`)?
   - Recommendation: Start flat, add hierarchy if needed

2. **Parameter Syntax**: Space-separated vs flag-based?
   - Recommendation: Space-separated by default, optional flags for complex cases

3. **Error Recovery**: Auto-correct vs suggest?
   - Recommendation: Suggest with confirmation for fuzzy matches

4. **Command Versioning**: How to handle breaking changes?
   - Recommendation: Semantic versioning in commands.yml, deprecation warnings

## 10. Appendix

### A. Example Commands YAML (Minimal)

```yaml
commands:
  ticket:
    tier: null
    workflow_ref: "workflows/ticket.md"
  m4:
    tier: "m4"
    workflow: "Use m4 tier for user request"

aliases:
  deepthink: "m4"
  bug: "ticket"

fuzzy_matching:
  enabled: true
  threshold: 0.8
```

### B. Fuzzy Matching Algorithm

```python
from difflib import SequenceMatcher

def fuzzy_match(input_cmd: str, commands: dict, threshold: float = 0.8) -> str:
    """Find best matching command using Levenshtein distance."""
    best_match = None
    best_ratio = 0

    for cmd in commands.keys():
        ratio = SequenceMatcher(None, input_cmd.lower(), cmd.lower()).ratio()
        if ratio > best_ratio and ratio >= threshold:
            best_ratio = ratio
            best_match = cmd

    return best_match
```

### C. Integration with llm_inference_auto

```python
# When /m4 detected
tool_call = {
    "name": "llm_inference_auto",
    "arguments": {
        "task_description": user_message,
        "messages": messages,
        "override_tier": "m4"  # ← Slash command tier signal
    }
}
```

## 11. References

- Fuzzy String Matching: https://github.com/seatgeek/fuzzywuzzy
- YAML Spec: https://yaml.org/spec/1.2.2/
- Command-line Interface Guidelines: https://clig.dev/
