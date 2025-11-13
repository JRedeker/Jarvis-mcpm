# Spec & Ticket Folder Structure Pattern

## Pattern

When creating or updating specifications that have related tickets, use this folder structure:

```
specs/
├── {spec-name}/
│   ├── README.md              # Overview, navigation, quick links
│   ├── tickets/               # Related tickets and requirements
│   │   ├── ticket-1.md
│   │   ├── ticket-2.md
│   │   └── analysis-*.md     # Gap analysis, implementation summaries
│   ├── notes/                 # Implementation notes, decisions, research
│   │   ├── decision-log.md
│   │   └── research-*.md
│   └── assets/                # Diagrams, mockups, examples
│       ├── architecture.png
│       └── mockup.png
└── {spec-name}-spec.md        # Main specification document
```

## Example: VS Code LM Provider

```
specs/
├── vscode-lm-provider-cipher/
│   ├── README.md
│   ├── tickets/
│   │   ├── implement-response-metadata-enrichment.md
│   │   ├── metadata-display-verification-report.md
│   │   ├── metadata-display-implementation-summary.md
│   │   └── metadata-display-architecture-issue.md
│   ├── notes/
│   └── assets/
└── vscode-lm-provider-cipher-spec.md
```

## Benefits

1. **Organization**: All related materials in one place
2. **Navigation**: README.md provides clear entry point
3. **Hierarchy**: Tickets live under specs that solve them
4. **Discoverability**: Easy to find related context
5. **Maintenance**: Clear ownership and status tracking

## When to Use

- Creating a new specification
- A ticket requires significant implementation (spec-worthy)
- Multiple tickets relate to same architectural solution
- Need to organize research, decisions, and assets

## README.md Template

Each spec folder should have a README with:
- Overview of the spec
- Directory structure explanation
- Link to main spec document
- Key use cases
- Implementation status
- Dependencies
- Quick links to tickets
- Decision log
- Contact/ownership

## Cross-referencing

- Tickets reference parent spec: `See: specs/{spec-name}/{spec-name}-spec.md`
- Spec references tickets: `See: tickets/ticket-name.md` (relative path)
- README links everything together

## Created: 2025-11-12
## Example: VS Code LM Provider Cipher spec
## Related: vs-code-lm-provider, ticket-management, project-organization