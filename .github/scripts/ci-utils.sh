#!/usr/bin/env bash
# CI utility functions using gum for beautiful output
# Source this file: source .github/scripts/ci-utils.sh
# Adapted for Go/Node.js project (Jarvis/MCPM)

set -eo pipefail

# Colors (ANSI 256)
_BLUE=212
_GREEN=82
_RED=196
_ORANGE=208
_YELLOW=220
_PURPLE=141
_CYAN=75
_GRAY=240

# Core display functions (internal)
_header() { gum style --border double --border-foreground "$_BLUE" --padding "1 2" "$@"; }
_success() { gum style --border double --border-foreground "$_GREEN" --padding "1 2" "$@"; }
_error() { gum style --border double --border-foreground "$_RED" --padding "1 2" "$@"; }
_warn() { gum style --border double --border-foreground "$_ORANGE" --padding "1 2" "$@"; }
_critical() { gum style --border thick --border-foreground "$_RED" --padding "1 2" "$@"; }
_high() { gum style --border thick --border-foreground "$_ORANGE" --padding "1 2" "$@"; }
_medium() { gum style --border normal --border-foreground "$_CYAN" --padding "1 2" "$@"; }
_low() { gum style --border normal --border-foreground "$_YELLOW" --padding "1 2" "$@"; }
_section() { local c="$1"; shift; gum style --border double --border-foreground "$c" --padding "1 2" "$@"; }

# =============================================================================
# CI CHECK FUNCTIONS - Just call these, output is automatic
# =============================================================================

ci_complexity() {
    _header "📊 COMPLEXITY SCAN" "" "Go: CCN ≤ 15 | NLOC ≤ 100"

    local exit_code=0

    # Run lizard on Go files
    if [ -d "Jarvis" ]; then
        _section "$_CYAN" "Scanning Jarvis (Go)..."

        # lizard outputs warnings for functions exceeding thresholds
        # CCN (Cyclomatic Complexity) <= 15, NLOC (Lines) <= 100, Args <= 6
        if ! lizard Jarvis/ \
            --CCN 15 \
            --length 100 \
            --arguments 6 \
            --warnings_only \
            --exclude "*/testing/*" \
            --exclude "*_test.go"; then
            _warn "⚠️ Some functions exceed complexity thresholds"
            exit_code=1
        fi
    fi

    # Scan Node.js if present
    if [ -d "MCPM" ] && [ -f "MCPM/package.json" ]; then
        _section "$_CYAN" "Scanning MCPM (JavaScript)..."

        if ! lizard MCPM/ \
            --CCN 15 \
            --length 100 \
            --arguments 6 \
            --warnings_only \
            --exclude "*/node_modules/*" \
            --exclude "*.test.js"; then
            _warn "⚠️ Some functions exceed complexity thresholds"
            exit_code=1
        fi
    fi

    if [ $exit_code -eq 0 ]; then
        _success "✅ COMPLEXITY PASSED"
    else
        _warn "⚠️ COMPLEXITY WARNINGS" "" "Review flagged functions for refactoring"
    fi

    return $exit_code
}

ci_review() {
    local pr_number="$1" pr_title="$2" pr_author="$3" head_ref="$4" base_ref="$5"

    local changed file_count additions deletions stats

    changed=$(git diff --name-only "origin/${base_ref}...HEAD" | grep -E '\.(go|js|ts|md|yml|yaml|toml|json)$' || true)
    file_count=$(echo "$changed" | grep -c . || echo "0")
    stats=$(git diff --stat "origin/${base_ref}...HEAD" | tail -1)
    additions=$(echo "$stats" | grep -oE '[0-9]+ insertion' | grep -oE '[0-9]+' || echo "0")
    deletions=$(echo "$stats" | grep -oE '[0-9]+ deletion' | grep -oE '[0-9]+' || echo "0")

    _section 99 "🤖 AI CODE REVIEW" "" \
        "PR #${pr_number}: ${pr_title}" \
        "Branch: ${head_ref} → ${base_ref}" \
        "Author: ${pr_author}" "" \
        "Files: ${file_count} | +${additions} -${deletions}"

    echo ""
    if [ -n "$changed" ]; then
        echo "$changed" | while read -r file; do echo "  • $file"; done
    fi

    _critical "🚨 CRITICAL (P01-P04) — Blocks merge" "" \
        "☐ No hardcoded secrets/credentials" \
        "☐ Least privilege enforced" \
        "☐ Commands have timeouts" \
        "☐ Behavior obvious from local context"

    _high "⚠️ HIGH (P05-P06) — Should fix" "" \
        "☐ No partial changes without tests" \
        "☐ Atomic, logically grouped commits"

    _medium "📋 MEDIUM (P07-P14)" "" \
        "☐ Tests for new behavior" \
        "☐ Dependencies verified" \
        "☐ Reduces technical debt" \
        "☐ Structured logs/traces"

    _low "💡 STANDARD (P15-P22)" "" \
        "☐ Fails fast with clear errors" \
        "☐ Simple, clear, well-named" \
        "☐ Docs current"

    _high "🔧 GO PROJECT RULES" "" \
        "☐ No //nolint without reason" \
        "☐ All exports documented" \
        "☐ Errors wrapped with context" \
        "☐ No panic in library code" \
        "☐ Interfaces accepted, structs returned"

    _section 141 "📝 REPORT FORMAT" "" \
        "### [SEVERITY] P## — file:line" \
        "**Issue:** Description" \
        "**Fix:** Suggestion" "" \
        "🚨 BLOCKING | ⚠️ MAJOR | 📋 MINOR | 💡 SUGGEST" "" \
        "Verdict: ✅ APPROVE | ❌ REQUEST_CHANGES | 💬 COMMENT"

    _section 99 "📄 DIFF (first 50KB)"
    git diff "origin/${base_ref}...HEAD" -- '*.go' '*.js' '*.md' '*.yml' '*.yaml' '*.toml' '*.json' | head -c 50000

    _section 99 "📄 END"
}

# =============================================================================
# ROADMAP CONSISTENCY CHECK
# =============================================================================

ci_roadmap() {
    _header "📋 ROADMAP CHECK" "" "Validates OpenSpec changes against archived specs"

    local errors=0
    local warnings=0

    # Check 1: All archived specs should be in CHANGELOG.md (if it exists)
    _section "$_CYAN" "Checking archived specs..."

    local archive_dir="openspec/changes/archive"
    if [ -d "$archive_dir" ]; then
        for spec_dir in "$archive_dir"/*/; do
            if [ -d "$spec_dir" ]; then
                local spec_name
                spec_name=$(basename "$spec_dir" | sed 's/^[0-9-]*//' | sed 's/^-//')

                if [ -f "CHANGELOG.md" ]; then
                    if ! grep -qi "$spec_name" CHANGELOG.md 2>/dev/null; then
                        echo "  ⚠️ Archived spec not in CHANGELOG: $spec_name"
                        ((warnings++))
                    else
                        echo "  ✓ $spec_name"
                    fi
                else
                    echo "  ✓ $spec_name (no CHANGELOG.md)"
                fi
            fi
        done
    else
        echo "  ℹ️ No archived specs found"
    fi

    # Check 2: Active specs should have required files
    _section "$_CYAN" "Checking active specs structure..."

    local changes_dir="openspec/changes"
    if [ -d "$changes_dir" ]; then
        for spec_dir in "$changes_dir"/*/; do
            if [ -d "$spec_dir" ] && [[ "$spec_dir" != *"/archive/"* ]]; then
                local spec_name
                spec_name=$(basename "$spec_dir")

                local has_proposal=false
                local has_design=false
                local has_tasks=false

                [ -f "$spec_dir/proposal.md" ] && has_proposal=true
                [ -f "$spec_dir/design.md" ] && has_design=true
                [ -f "$spec_dir/tasks.md" ] && has_tasks=true

                if $has_proposal || $has_design; then
                    if ! $has_tasks; then
                        echo "  ⚠️ Missing tasks.md: $spec_name"
                        ((warnings++))
                    else
                        echo "  ✓ $spec_name"
                    fi
                fi
            fi
        done
    else
        echo "  ℹ️ No active specs found"
    fi

    # Check 3: Verify openspec/project.md exists
    _section "$_CYAN" "Checking project metadata..."

    if [ -f "openspec/project.md" ]; then
        echo "  ✓ openspec/project.md exists"
    else
        echo "  ⚠️ Missing openspec/project.md"
        ((warnings++))
    fi

    # Summary
    echo ""
    if [ $errors -gt 0 ]; then
        _error "❌ ROADMAP CHECK FAILED" "" "$errors error(s), $warnings warning(s)"
        return 1
    elif [ $warnings -gt 0 ]; then
        _warn "⚠️ ROADMAP CHECK: $warnings warning(s)" "" \
            "Update CHANGELOG.md with archived specs" \
            "Ensure active specs have tasks.md"
        return 0  # Warnings don't fail CI, but are visible
    else
        _success "✅ ROADMAP CONSISTENT"
    fi
}

# =============================================================================
# INSTALL HELPER
# =============================================================================

ci_install_gum() {
    if command -v gum &> /dev/null; then return 0; fi

    sudo mkdir -p /etc/apt/keyrings
    curl -fsSL https://repo.charm.sh/apt/gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/charm.gpg
    echo "deb [signed-by=/etc/apt/keyrings/charm.gpg] https://repo.charm.sh/apt/ * *" | sudo tee /etc/apt/sources.list.d/charm.list
    sudo apt update && sudo apt install -y gum
}
