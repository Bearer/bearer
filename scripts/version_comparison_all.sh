#!/bin/bash
# Version Comparison Script - Run All Repositories
# Runs version comparison on all repositories from the version_comparison list
#
# Usage: ./scripts/version_comparison_all.sh [--limit N]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
REPOS_FILE="${PROJECT_ROOT}/.github/workflows/version_comparison/repositories.json5"
OUTPUT_DIR="${PROJECT_ROOT}/tmp/version_comparison"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

LIMIT=""
if [ "$1" == "--limit" ] && [ -n "$2" ]; then
    LIMIT="$2"
fi

# Check for json5 parser (npx json5)
if ! command -v npx &> /dev/null; then
    echo -e "${RED}Error: npx not found. Please install Node.js.${NC}"
    exit 1
fi

echo -e "${GREEN}=== Version Comparison - All Repositories ===${NC}"
echo ""

# Parse repositories from json5
REPOS=$(npx --yes json5 "$REPOS_FILE" 2>/dev/null | jq -r '.include[] | "\(.name)|\(.repository_url)"')

COUNT=0
TOTAL=$(echo "$REPOS" | wc -l | tr -d ' ')

if [ -n "$LIMIT" ]; then
    echo "Running on first $LIMIT of $TOTAL repositories"
    REPOS=$(echo "$REPOS" | head -n "$LIMIT")
    TOTAL="$LIMIT"
fi

echo "Total repositories: $TOTAL"
echo ""

# Create master summary file
MASTER_SUMMARY="${OUTPUT_DIR}/master_summary.txt"
mkdir -p "$OUTPUT_DIR"

cat > "$MASTER_SUMMARY" << EOF
Version Comparison Master Summary
=================================
Generated: $(date)
Total Repositories: $TOTAL

EOF

while IFS='|' read -r NAME URL; do
    COUNT=$((COUNT + 1))
    echo -e "${YELLOW}[$COUNT/$TOTAL] Processing: $NAME${NC}"

    # Run comparison script
    if "$SCRIPT_DIR/version_comparison.sh" "$URL" "$NAME" > "${OUTPUT_DIR}/${NAME}/run.log" 2>&1; then
        # Extract key stats from summary
        SUMMARY_FILE="${OUTPUT_DIR}/${NAME}/summary.txt"
        if [ -f "$SUMMARY_FILE" ]; then
            BASE_FINDINGS=$(grep "Base:" "$SUMMARY_FILE" | head -1 | awk '{print $2}')
            TEST_FINDINGS=$(grep "Test:" "$SUMMARY_FILE" | head -1 | awk '{print $2}')
            ONLY_IN_BASE=$(grep "Only in Base" "$SUMMARY_FILE" | awk '{print $NF}')
            ONLY_IN_TEST=$(grep "Only in Test" "$SUMMARY_FILE" | awk '{print $NF}')

            STATUS="✓"
            if [ "$ONLY_IN_BASE" -gt 0 ]; then
                STATUS="⚠️ MISSING"
            fi

            echo "  Base: $BASE_FINDINGS, Test: $TEST_FINDINGS, Missing: $ONLY_IN_BASE, New: $ONLY_IN_TEST $STATUS"

            printf "%-30s Base: %4s  Test: %4s  Missing: %3s  New: %3s  %s\n" \
                "$NAME" "$BASE_FINDINGS" "$TEST_FINDINGS" "$ONLY_IN_BASE" "$ONLY_IN_TEST" "$STATUS" >> "$MASTER_SUMMARY"
        fi
    else
        echo -e "${RED}  Failed - check ${OUTPUT_DIR}/${NAME}/run.log${NC}"
        printf "%-30s FAILED\n" "$NAME" >> "$MASTER_SUMMARY"
    fi
done <<< "$REPOS"

echo ""
echo -e "${GREEN}=== Master Summary ===${NC}"
cat "$MASTER_SUMMARY"

echo ""
echo -e "${GREEN}Results saved to: ${OUTPUT_DIR}${NC}"
echo -e "${GREEN}Master summary: ${MASTER_SUMMARY}${NC}"

# Report repos with potential issues
echo ""
echo -e "${YELLOW}=== Repositories Needing Review ===${NC}"
grep "MISSING\|FAILED" "$MASTER_SUMMARY" || echo "None - all good!"

