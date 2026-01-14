#!/bin/bash
# Version Comparison Script
# Compares findings between the current release (bearer) and dev version (go run)
#
# Usage: ./scripts/version_comparison.sh <repo_url> [repo_name]
# Example: ./scripts/version_comparison.sh https://github.com/Bearer/railsgoat railsgoat

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
OUTPUT_DIR="${PROJECT_ROOT}/tmp/version_comparison"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

usage() {
    echo "Usage: $0 <repo_url> [repo_name]"
    echo ""
    echo "Arguments:"
    echo "  repo_url   - Git URL of the repository to scan"
    echo "  repo_name  - Optional name for the repo (defaults to basename of URL)"
    echo ""
    echo "Examples:"
    echo "  $0 https://github.com/Bearer/railsgoat"
    echo "  $0 https://github.com/Bearer/railsgoat railsgoat"
    echo ""
    echo "Output will be written to: ${OUTPUT_DIR}/<repo_name>/"
    exit 1
}

if [ -z "$1" ]; then
    usage
fi

REPO_URL="$1"
REPO_NAME="${2:-$(basename "$REPO_URL" .git)}"

echo -e "${GREEN}=== Version Comparison Script ===${NC}"
echo "Repository: $REPO_URL"
echo "Name: $REPO_NAME"
echo ""

# Create output directory
REPO_OUTPUT_DIR="${OUTPUT_DIR}/${REPO_NAME}"
mkdir -p "$REPO_OUTPUT_DIR"

# Clone repository (skip if already exists)
REPO_DIR="${OUTPUT_DIR}/repos/${REPO_NAME}"
echo -e "${YELLOW}[1/5] Cloning repository...${NC}"
if [ -d "$REPO_DIR/.git" ]; then
    echo "  Repository already exists, skipping clone."
else
    mkdir -p "$(dirname "$REPO_DIR")"
    git clone --single-branch --depth 1 --no-tags "$REPO_URL" "$REPO_DIR" 2>&1 | grep -v "^$" || true
    echo "  Done."
fi

# Run base scan (current release)
echo -e "${YELLOW}[2/5] Running base scan (current release)...${NC}"
BASE_JSON="${REPO_OUTPUT_DIR}/base.json"
if command -v bearer &> /dev/null; then
    bearer scan "$REPO_DIR" \
        --format jsonv2 \
        --exit-code 0 \
        --force \
        --disable-version-check \
        --quiet \
        --hide-progress-bar \
        2>/dev/null > "$BASE_JSON" || true
    BASE_VERSION=$(jq -r '.version // "unknown"' "$BASE_JSON")
    echo "  Base version: $BASE_VERSION"
else
    echo -e "${RED}  Error: 'bearer' command not found. Please install the release version.${NC}"
    exit 1
fi

# Run test scan (dev version)
echo -e "${YELLOW}[3/5] Running test scan (dev version)...${NC}"
TEST_JSON="${REPO_OUTPUT_DIR}/test.json"
cd "$PROJECT_ROOT"
go run cmd/bearer/bearer.go scan "$REPO_DIR" \
    --format jsonv2 \
    --exit-code 0 \
    --force \
    --disable-version-check \
    --quiet \
    --hide-progress-bar \
    2>/dev/null > "$TEST_JSON" || true
TEST_VERSION=$(jq -r '.version // "unknown"' "$TEST_JSON")
echo "  Test version: $TEST_VERSION"

# Generate comparison files
echo -e "${YELLOW}[4/5] Generating comparison files...${NC}"

# Extract key finding info for comparison
jq -S '.findings[] | {id: .id, filename: .filename, line: .line_number, fingerprint: .fingerprint, severity: .severity}' "$BASE_JSON" 2>/dev/null | sort > "${REPO_OUTPUT_DIR}/base_findings.txt"
jq -S '.findings[] | {id: .id, filename: .filename, line: .line_number, fingerprint: .fingerprint, severity: .severity}' "$TEST_JSON" 2>/dev/null | sort > "${REPO_OUTPUT_DIR}/test_findings.txt"

# Generate diff
diff -u "${REPO_OUTPUT_DIR}/base_findings.txt" "${REPO_OUTPUT_DIR}/test_findings.txt" > "${REPO_OUTPUT_DIR}/findings_diff.txt" 2>&1 || true

# Extract errors
jq -S '.errors // []' "$BASE_JSON" > "${REPO_OUTPUT_DIR}/base_errors.json"
jq -S '.errors // []' "$TEST_JSON" > "${REPO_OUTPUT_DIR}/test_errors.json"

# Generate summary
echo -e "${YELLOW}[5/5] Generating summary...${NC}"
SUMMARY_FILE="${REPO_OUTPUT_DIR}/summary.txt"

BASE_FINDINGS=$(jq '.findings | length' "$BASE_JSON")
TEST_FINDINGS=$(jq '.findings | length' "$TEST_JSON")
BASE_ERRORS=$(jq 'length' "${REPO_OUTPUT_DIR}/base_errors.json")
TEST_ERRORS=$(jq 'length' "${REPO_OUTPUT_DIR}/test_errors.json")

# Count findings only in base (removed/missed)
ONLY_IN_BASE_RAW=$(comm -23 "${REPO_OUTPUT_DIR}/base_findings.txt" "${REPO_OUTPUT_DIR}/test_findings.txt" 2>/dev/null | grep -c "fingerprint" 2>/dev/null) || true
ONLY_IN_BASE=${ONLY_IN_BASE_RAW:-0}
# Count findings only in test (new)
ONLY_IN_TEST_RAW=$(comm -13 "${REPO_OUTPUT_DIR}/base_findings.txt" "${REPO_OUTPUT_DIR}/test_findings.txt" 2>/dev/null | grep -c "fingerprint" 2>/dev/null) || true
ONLY_IN_TEST=${ONLY_IN_TEST_RAW:-0}

cat > "$SUMMARY_FILE" << EOF
Version Comparison Summary: ${REPO_NAME}
==========================================
Repository: ${REPO_URL}
Base Version: ${BASE_VERSION}
Test Version: ${TEST_VERSION}
Generated: $(date)

Findings:
  Base: ${BASE_FINDINGS}
  Test: ${TEST_FINDINGS}
  Difference: $((TEST_FINDINGS - BASE_FINDINGS))

  Only in Base (potentially missed): ${ONLY_IN_BASE}
  Only in Test (new findings): ${ONLY_IN_TEST}

Errors:
  Base: ${BASE_ERRORS}
  Test: ${TEST_ERRORS}
  New errors: $((TEST_ERRORS - BASE_ERRORS))

Files Generated:
  - base.json: Full scan output from release version
  - test.json: Full scan output from dev version
  - base_findings.txt: Extracted findings (sorted)
  - test_findings.txt: Extracted findings (sorted)
  - findings_diff.txt: Diff between findings
  - base_errors.json: Errors from base scan
  - test_errors.json: Errors from test scan
  - summary.txt: This file
EOF

# Print summary
echo ""
echo -e "${GREEN}=== Summary ===${NC}"
cat "$SUMMARY_FILE"

# Highlight potential issues
echo ""
if [ "$ONLY_IN_BASE" -gt 0 ]; then
    echo -e "${RED}⚠️  WARNING: ${ONLY_IN_BASE} findings in base are NOT in test (potentially missed)${NC}"
    echo "   Review: ${REPO_OUTPUT_DIR}/findings_diff.txt"
fi

if [ "$ONLY_IN_TEST" -gt 0 ]; then
    echo -e "${GREEN}✓  ${ONLY_IN_TEST} new findings in test version${NC}"
fi

if [ "$TEST_ERRORS" -gt "$BASE_ERRORS" ]; then
    echo -e "${YELLOW}ℹ️  ${TEST_ERRORS} errors in test vs ${BASE_ERRORS} in base${NC}"
    echo "   Review: ${REPO_OUTPUT_DIR}/test_errors.json"
fi

echo ""
echo -e "${GREEN}Output directory: ${REPO_OUTPUT_DIR}${NC}"

