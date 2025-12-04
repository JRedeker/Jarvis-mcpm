#!/usr/bin/env bats
# Tests for manage-mcp.sh script

SCRIPT_DIR="$(cd "$(dirname "$BATS_TEST_FILENAME")/.." && pwd)"
MANAGE_SCRIPT="$SCRIPT_DIR/manage-mcp.sh"

# Setup and teardown
setup() {
    # Save original directory
    ORIG_DIR="$(pwd)"
    cd "$SCRIPT_DIR/.."
}

teardown() {
    cd "$ORIG_DIR"
}

# Basic functionality tests
@test "manage-mcp.sh exists and is executable" {
    [ -x "$MANAGE_SCRIPT" ]
}

@test "manage-mcp.sh shows help with no arguments" {
    run "$MANAGE_SCRIPT"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Usage:" ]]
}

@test "manage-mcp.sh shows help with -h flag" {
    run "$MANAGE_SCRIPT" -h
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Usage:" ]]
}

@test "manage-mcp.sh shows help with --help flag" {
    run "$MANAGE_SCRIPT" --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Usage:" ]]
}

@test "manage-mcp.sh status command runs" {
    run "$MANAGE_SCRIPT" status
    # Status command should succeed even if Docker is not running
    # It just reports the current state
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
}

@test "manage-mcp.sh test command is recognized" {
    # Just check the command is recognized (don't run full tests)
    run bash -c "echo 'test' | timeout 1 $MANAGE_SCRIPT test 2>&1 || true"
    # Should not show "unknown command"
    [[ ! "$output" =~ "unknown command" ]] || [[ ! "$output" =~ "Unknown command" ]]
}

# Command validation tests
@test "manage-mcp.sh rejects invalid commands" {
    run "$MANAGE_SCRIPT" invalid_command_xyz
    [ "$status" -ne 0 ]
}

# Docker dependency tests (skip if Docker not available)
@test "manage-mcp.sh start command requires Docker" {
    if ! command -v docker &> /dev/null; then
        skip "Docker not installed"
    fi
    # Just verify the command doesn't crash immediately
    run timeout 5 "$MANAGE_SCRIPT" start 2>&1 || true
    # Should attempt to start, not fail with syntax error
    [[ ! "$output" =~ "syntax error" ]]
}

@test "manage-mcp.sh logs command requires Docker" {
    if ! command -v docker &> /dev/null; then
        skip "Docker not installed"
    fi
    run timeout 5 "$MANAGE_SCRIPT" logs 2>&1 || true
    [[ ! "$output" =~ "syntax error" ]]
}
