#!/usr/bin/env bats
# Tests for generate-api-docs.sh script
# Updated for Jarvis v3.0 consolidated tools

SCRIPT_DIR="$(cd "$(dirname "$BATS_TEST_FILENAME")/.." && pwd)"
GEN_SCRIPT="$SCRIPT_DIR/generate-api-docs.sh"
OUTPUT_FILE="/tmp/test-api-docs-$$.md"

# Setup and teardown
setup() {
    ORIG_DIR="$(pwd)"
    cd "$SCRIPT_DIR/.."
}

teardown() {
    cd "$ORIG_DIR"
    rm -f "$OUTPUT_FILE"
}

# Basic functionality tests
@test "generate-api-docs.sh exists and is executable" {
    [ -x "$GEN_SCRIPT" ]
}

@test "generate-api-docs.sh has valid bash syntax" {
    run bash -n "$GEN_SCRIPT"
    [ "$status" -eq 0 ]
}

@test "generate-api-docs.sh generates output file" {
    run "$GEN_SCRIPT" "$OUTPUT_FILE"
    [ "$status" -eq 0 ]
    [ -f "$OUTPUT_FILE" ]
}

@test "generated docs contain header" {
    "$GEN_SCRIPT" "$OUTPUT_FILE"
    run grep -q "Jarvis API Reference" "$OUTPUT_FILE"
    [ "$status" -eq 0 ]
}

# v3.0 consolidated tool names
@test "generated docs contain jarvis_check_status tool" {
    "$GEN_SCRIPT" "$OUTPUT_FILE"
    run grep -q "jarvis_check_status" "$OUTPUT_FILE"
    [ "$status" -eq 0 ]
}

@test "generated docs contain jarvis_server tool" {
    "$GEN_SCRIPT" "$OUTPUT_FILE"
    run grep -q "jarvis_server" "$OUTPUT_FILE"
    [ "$status" -eq 0 ]
}

@test "generated docs contain jarvis_profile tool" {
    "$GEN_SCRIPT" "$OUTPUT_FILE"
    run grep -q "jarvis_profile" "$OUTPUT_FILE"
    [ "$status" -eq 0 ]
}

@test "generated docs contain jarvis_system tool" {
    "$GEN_SCRIPT" "$OUTPUT_FILE"
    run grep -q "jarvis_system" "$OUTPUT_FILE"
    [ "$status" -eq 0 ]
}

@test "generated docs contain category sections" {
    "$GEN_SCRIPT" "$OUTPUT_FILE"
    run grep -c "^## " "$OUTPUT_FILE"
    [ "$status" -eq 0 ]
    [ "$output" -ge 1 ]  # Should have at least 1 category section
}

@test "generated docs have tool descriptions" {
    "$GEN_SCRIPT" "$OUTPUT_FILE"
    # Check for descriptions (lines after tool headers)
    run grep -c "health\|server\|profile\|system" "$OUTPUT_FILE"
    [ "$status" -eq 0 ]
    [ "$output" -ge 3 ]
}
