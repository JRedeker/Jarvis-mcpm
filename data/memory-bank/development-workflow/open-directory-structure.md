# ./open/ Directory Structure

## Overview

The `./open/` directory contains all active work items organized by type: tickets (user requests), specs (architectural solutions), and alerts (system-detected issues).

## Structure

```
./open/
├── tickets/           # User-requested features, bugs, improvements
│   └── {ticket-name}/
│       ├── {ticket-name}.md      # Main ticket description
│       └── notes/                 # Optional research, context
├── specs/            # Architectural specifications and solutions
│   └── {spec-name}/
│       ├── {spec-name}.md        # Spec index/overview (named after directory)
│       ├── spec.md                # Main specification document
│       │                          # OR: proposal.md, plan.md, tasks.md (workflow dependent)
│       ├── tickets/               # Tickets this spec solves
│       │   └── {ticket-name}/
│       │       └── {ticket-name}.md
│       ├── alerts/                # Alerts this spec solves
│       │   └── {alert-name}/
│       │       └── {alert-name}.md
│       ├── notes/                 # Implementation notes, decisions
│       └── assets/                # Diagrams, mockups, examples
└── alerts/           # System-detected issues with tooling/setup/internals
    └── {alert-name}/
        ├── {alert-name}.md       # Alert description
        ├── detection-log.md       # When/how detected
        └── context/               # Logs, screenshots, evidence
```

## Lifecycle & Migration

**New Items:**
- Tickets start in `./open/tickets/`
- Alerts start in `./open/alerts/`
- Specs start in `./open/specs/`

**When Solutions Identified:**
- If a ticket/alert would be solved by a spec, move its directory to that spec:
  - `./open/tickets/problem-x/` → `./open/specs/solution-y/tickets/problem-x/`
  - `./open/alerts/alert-z/` → `./open/specs/solution-y/alerts/alert-z/`

**Completion:**
- Delete from `./open/` when work is complete
- Git history preserves all artifacts
- No archive directory needed

## Naming Conventions

**Directory Names:**
- Use kebab-case: `metadata-display`, `performance-optimization`
- Be descriptive but concise

**File Names:**
For specs:
- Index file: `{spec-name}.md` (matches parent directory name)
  - Example: `./open/specs/vscode-lm-provider-cipher/vscode-lm-provider-cipher.md`
- Detailed spec: `spec.md`, `proposal.md`, `plan.md`, or `tasks.md` (workflow dependent)
  - Example: `./open/specs/vscode-lm-provider-cipher/spec.md`
- DO NOT create `README.md` - the index file serves this purpose

For tickets/alerts:
- Main document: `{name}.md` (matches directory name)
  - Example: `./open/tickets/metadata-display/metadata-display.md`

**Cross-References:**
- Tickets reference specs: `Spec: ../../specs/vscode-lm-provider/vscode-lm-provider.md`
- Specs reference tickets: `Tickets: ./tickets/metadata-display/metadata-display.md`

## Item Types

### Tickets
User-requested work items:
- New features
- Bug fixes
- Improvements
- Refactoring requests
- User-reported issues

### Specs
Architectural solutions:
- Design specifications
- Implementation plans
- Architecture documents
- Integration specifications
- System designs

### Alerts
System-detected issues:
- MCP server connection failures
- Missing dependencies
- Configuration drift
- Performance degradation
- Setup/tooling problems
- Internal inconsistencies

**Alert Severity:**
Prefix alerts with severity for priority:
- `critical-{issue}` - System down, blocking work
- `warning-{issue}` - Degraded performance, needs attention
- `info-{issue}` - Configuration drift, non-urgent

## Benefits

1. **Visibility**: All active work in one place
2. **Organization**: Clear categorization by type
3. **Relationships**: Tickets/alerts linked to solving specs
4. **Cleanliness**: Completed work removed, not archived
5. **Self-Service**: System can create alerts automatically

## Examples

**New Ticket:**
```
./open/tickets/add-cost-tracking/
└── add-cost-tracking.md
```

**Spec with Ticket:**
```
./open/specs/llm-cost-monitoring/
├── llm-cost-monitoring.md     # Index/overview
├── spec.md                     # Detailed specification
└── tickets/
    └── add-cost-tracking/
        └── add-cost-tracking.md
```

**Spec with Multiple Supporting Files (SpecKit workflow example):**
```
./open/specs/vscode-lm-provider-cipher/
├── vscode-lm-provider-cipher.md     # Index/overview
├── spec.md                          # Main specification
├── proposal.md                      # Initial proposal
├── tasks.md                         # Implementation tasks
├── tickets/
│   ├── implement-response-metadata-enrichment/
│   │   └── implement-response-metadata-enrichment.md
│   └── metadata-display-verification-report/
│       └── metadata-display-verification-report.md
├── notes/
│   └── decision-log.md
└── assets/
    └── architecture-diagram.png
```

**System Alert:**
```
./open/alerts/warning-mcp-timeout/
├── warning-mcp-timeout.md
├── detection-log.md
└── context/
    └── error-logs.txt
```

## Migration from Old Structure

When migrating from `specs/` and `tickets/` at root:

1. Move spec directories: `specs/{name}/ → ./open/specs/{name}/`
2. Move spec files: `specs/{name}-spec.md → ./open/specs/{name}/{name}.md`
3. Move tickets: `tickets/{name}.md → ./open/tickets/{name}/{name}.md`
4. Update all cross-references
5. Check for nested tickets in specs and preserve structure

## File Naming Convention Details

**Spec Index File**:
- **Purpose**: Navigation, quick links, status overview
- **Name**: Must match directory name (e.g., `vscode-lm-provider-cipher.md`)
- **Content**: Directory structure, key use cases, quick links, status
- **Common Mistake**: Creating `README.md` instead (don't do this)

**Supporting Spec Files**:
- **spec.md**: Detailed technical specification
- **proposal.md**: Initial proposal or RFC
- **plan.md**: Implementation plan
- **tasks.md**: Task breakdown (SpecKit workflow)
- Choose names based on your workflow/methodology

**Why Not README.md?**:
- The index file name should match the spec directory name
- This creates clear association between directory and its main file
- README.md is generic and doesn't convey what the spec is about
- Consistency: tickets use `{ticket-name}.md`, specs use `{spec-name}.md`

## Created: 2025-11-12
## Updated: 2025-11-12 (Corrected file naming conventions)
## Replaces: spec-ticket-folder-structure.md (evolves that pattern)
## Related: project-organization, ticket-management, spec-management
