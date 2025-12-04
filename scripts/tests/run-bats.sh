#!/bin/bash
# run-bats.sh - Run all bats tests
#
# Usage: ./scripts/tests/run-bats.sh
#
# Requires: bats (https://github.com/bats-core/bats-core)
# Install: npm install -g bats or apt install bats

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Check if bats is installed
if ! command -v bats &> /dev/null; then
    echo "Error: bats is not installed"
    echo "Install with: npm install -g bats"
    echo "         or: apt install bats"
    exit 1
fi

echo "Running bats tests..."
echo "====================="

# Run all .bats files
bats "$SCRIPT_DIR"/*.bats

echo ""
echo "All bats tests passed!"
