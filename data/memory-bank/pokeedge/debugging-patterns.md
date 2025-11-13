# Debugging Patterns - PokeEdge Project

## Systematic Debugging Framework

### Phase 1: Root Cause Investigation
**What to look for:**
- Error messages and stack traces
- Recent code changes that might have introduced issues
- Boundary conditions and edge cases
- Data flow through the system

**Evidence gathering techniques:**
- Read full error messages and tracebacks
- Reproduce the issue in isolation
- Check git history for recent changes
- Instrument code boundaries to trace data flow
- Add logging at key decision points

### Phase 2: Pattern Analysis
**Finding working examples:**
- Locate similar functionality that works correctly
- Compare implementation differences
- Understand dependency relationships
- Identify what makes the working version different

**Common patterns in PokeEdge:**
- API integration issues often stem from authentication or rate limiting
- TUI rendering problems usually involve state management or event handling
- CLI command failures typically relate to argument parsing or dependency injection
- Search functionality issues often involve data validation or transformation

### Phase 3: Hypothesis & Testing
**Forming hypotheses:**
- Single, testable explanation for the issue
- Based on evidence from phases 1-2
- Specific enough to verify or falsify

**Testing approach:**
- Minimal changes to test the hypothesis
- Isolate the suspected problematic code
- Verify results with automated tests where possible

### Phase 4: Implementation
**Creating failing tests first:**
- Ensures the fix addresses the actual issue
- Prevents regression
- Documents expected behavior

**Implementation principles:**
- Minimal change addressing root cause
- One fix per issue
- Verify globally that the fix doesn't break other functionality

## Red Flag Phrases (Stop & Return to Phase 1)

- "This should work..."
- "It worked before, I didn't change anything"
- "Let me try this quick fix"
- "We'll just add a workaround for now"
- "It seems like it's working..."
- "Probably just a temporary glitch"

## Partner Signals

When working with other developers:
- "Is that not happening?" → Unverified assumption
- "Will it show us...?" → Missing evidence gathering
- "Stop guessing" → Need systematic analysis
- "Have you tried reproducing it?" → Missing Phase 1

## Integration with PokeEdge Architecture

### API Layer Debugging
- Check dependency injection in api/dependencies.py
- Verify middleware configuration in api/middleware/
- Validate model schemas in api/models/
- Review error handling in api/exceptions.py

### CLI Debugging
- Command parsing in cli/commands/
- Presenter logic in cli/presenters/
- Integration with config system in config/

### TUI Debugging
- State management in models/ui_state.py
- Event handling in tui/ directory
- Presenter logic integration

### Client Integration Debugging
- Authentication issues in clients/base_client.py
- Rate limiting in clients/retry/
- Error propagation through the response chain

## Success Verification

- Automated tests pass
- Broader test suite remains green
- Issue can no longer be reproduced
- Root cause documented
- Minimal, focused change applied