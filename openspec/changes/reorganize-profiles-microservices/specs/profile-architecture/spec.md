# Profile Architecture Specification

## ADDED Requirements

### Requirement: Profile Segmentation
The system MUST support granular "Micro-Profiles" to isolate failure domains.

#### Scenario: Segregated Research Tools
Given the `research` profile contains high-latency tools like `brave-search`
When `brave-search` fails to start or times out
Then the `essentials` profile (containing `time`) MUST remain accessible
And the `dev-core` profile (containing `morph-fast-apply`) MUST remain accessible.

#### Scenario: Parallel Startup
Given multiple profiles are configured in the client
When the system starts
Then the profiles SHOULD start independently (allocated to distinct ports/processes)
So that a slow startup in one does not block the others.

### Requirement: Standard Profile Set
The system MUST provide the following standard profiles by default or convention:

1.  **Essentials:** Minimal, fast, robust tools.
2.  **Memory:** Persistence layer.
3.  **Dev-Core:** Coding intelligence.
4.  **Research:** External/Web access (High Risk).
5.  **Data:** Heavy storage/Vector DB.

#### Scenario: Standard Profiles Exist
Given a fresh installation
When the user lists profiles
Then they SHOULD see `essentials`, `memory`, `dev-core`, `research`, and `data` (or be able to create them easily).

### Requirement: Client Configuration
The Client Configuration MUST support loading multiple profiles as a stack.

#### Scenario: Stack Composition
Given a client like OpenCode
When configured with `add_profile="essentials,memory,dev-core,research"`
Then it SHOULD have access to tools from ALL those profiles simultaneously.
