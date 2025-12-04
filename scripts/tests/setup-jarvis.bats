#!/usr/bin/env bats
# Tests for setup-jarvis.sh script

SCRIPT_DIR="$(cd "$(dirname "$BATS_TEST_FILENAME")/.." && pwd)"
SETUP_SCRIPT="$SCRIPT_DIR/setup-jarvis.sh"

# Setup and teardown
setup() {
    ORIG_DIR="$(pwd)"
    cd "$SCRIPT_DIR/.."
}

teardown() {
    cd "$ORIG_DIR"
}

# Basic functionality tests
@test "setup-jarvis.sh exists and is executable" {
    [ -x "$SETUP_SCRIPT" ]
}

@test "setup-jarvis.sh has valid bash syntax" {
    run bash -n "$SETUP_SCRIPT"
    [ "$status" -eq 0 ]
}

@test "setup-jarvis.sh detects Go availability" {
    if ! command -v go &> /dev/null; then
        skip "Go not installed"
    fi
    # Script should not fail immediately when Go is available
    run timeout 10 bash -c "echo 'n' | $SETUP_SCRIPT 2>&1" || true
    [[ ! "$output" =~ "Go is required" ]] || [ "$status" -eq 0 ]
}

# Directory structure tests
@test "Jarvis directory exists" {
    [ -d "$SCRIPT_DIR/../Jarvis" ]
}

@test "Jarvis go.mod exists" {
    [ -f "$SCRIPT_DIR/../Jarvis/go.mod" ]
}

@test "MCPM directory exists" {
    [ -d "$SCRIPT_DIR/../MCPM" ]
}

@test "MCPM package.json exists" {
    [ -f "$SCRIPT_DIR/../MCPM/package.json" ]
}
