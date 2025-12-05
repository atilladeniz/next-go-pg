#!/bin/bash

# Fetch LLM-friendly documentation from any URL
#
# Priority:
#   1. llms-full.txt / llms.txt (if available)
#   2. sitefetch (crawls entire site) - requires: bun/npm install -g sitefetch
#   3. Jina Reader (single page fallback)
#
# Features:
# - Auto-detects latest version for versioned docs
# - Supports llms.txt standard
# - Can crawl entire documentation sites
#
# Usage:
#   ./scripts/fetch-docs.sh <url> [output-name] [--single]
#
# Examples:
#   ./scripts/fetch-docs.sh https://tanstack.com/query/latest/docs
#   ./scripts/fetch-docs.sh https://nextjs.org/docs nextjs
#   ./scripts/fetch-docs.sh https://example.com/docs example --single  # Only single page

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
DIM='\033[2m'
NC='\033[0m'

# Config
DOCS_DIR="$(git rev-parse --show-toplevel)/.docs"
JINA_BASE="https://r.jina.ai"
USER_AGENT="Mozilla/5.0 (compatible; DocsBot/1.0)"

# Ensure .docs directory exists
mkdir -p "$DOCS_DIR"

# Parse arguments
URL="$1"
OUTPUT_NAME=""
SINGLE_PAGE=false

# Parse remaining arguments
shift || true
while [[ $# -gt 0 ]]; do
    case $1 in
        --single|-s)
            SINGLE_PAGE=true
            shift
            ;;
        *)
            if [ -z "$OUTPUT_NAME" ]; then
                OUTPUT_NAME="$1"
            fi
            shift
            ;;
    esac
done

if [ -z "$URL" ]; then
    echo -e "${RED}Error: URL required${NC}"
    echo ""
    echo "Usage: $0 <url> [output-name] [--single]"
    echo ""
    echo "Options:"
    echo "  --single, -s    Only fetch single page (skip site crawling)"
    echo ""
    echo "Examples:"
    echo "  $0 https://tanstack.com/query/latest/docs"
    echo "  $0 https://nextjs.org/docs nextjs"
    echo "  $0 https://example.com/page --single"
    echo ""
    echo "Requirements for full site crawling:"
    echo "  bun install -g sitefetch"
    exit 1
fi

# Extract domain and path for naming
extract_name() {
    local url="$1"
    local clean="${url#https://}"
    clean="${clean#http://}"
    clean="${clean#www.}"
    local domain="${clean%%/*}"
    local path="${clean#*/}"
    path="${path%%/*}"

    if [ "$path" != "$clean" ] && [ -n "$path" ] && [ "$path" != "docs" ] && [ "$path" != "documentation" ]; then
        echo "${domain%%.*}-${path}" | tr '[:upper:]' '[:lower:]' | tr -cd 'a-z0-9-'
    else
        echo "${domain%%.*}" | tr '[:upper:]' '[:lower:]' | tr -cd 'a-z0-9-'
    fi
}

# Detect base URL (handle versioned docs)
detect_base_url() {
    local url="$1"

    if [[ "$url" =~ /(latest|stable|current)(/|$) ]]; then
        echo "$url"
        return
    fi

    if [[ "$url" =~ /v?[0-9]+(\.[0-9x]+)*(/|$) ]]; then
        local base=$(echo "$url" | sed -E 's|/v?[0-9]+(\.[0-9x]+)*(/\|$)|/latest\2|')
        local status=$(curl -s -o /dev/null -w "%{http_code}" -L "$base" 2>/dev/null || echo "000")
        if [ "$status" = "200" ]; then
            echo -e "${YELLOW}ℹ Detected versioned docs, using /latest${NC}" >&2
            echo "$base"
            return
        fi
    fi

    echo "$url"
}

# Extract docs match pattern from URL
get_docs_pattern() {
    local url="$1"
    local clean="${url#https://}"
    clean="${clean#http://}"
    local path="/${clean#*/}"

    # Remove trailing slash
    path="${path%/}"

    if [ "$path" = "/$clean" ] || [ -z "$path" ] || [ "$path" = "/" ]; then
        echo "/docs/**"
    else
        echo "${path}/**"
    fi
}

# Try to fetch llms.txt or llms-full.txt
fetch_llms_txt() {
    local url="$1"
    local protocol="${url%%://*}"
    local rest="${url#*://}"
    local domain="${rest%%/*}"
    local base_url="${protocol}://${domain}"
    local docs_base=$(echo "$url" | sed -E 's|(/docs?(/.*)?)?$||')

    local urls_to_try=(
        "${base_url}/llms-full.txt"
        "${base_url}/llms.txt"
        "${docs_base}/llms-full.txt"
        "${docs_base}/llms.txt"
        "${url}/llms-full.txt"
        "${url}/llms.txt"
    )

    local unique_urls=($(printf '%s\n' "${urls_to_try[@]}" | sort -u))

    for try_url in "${unique_urls[@]}"; do
        echo -e "${DIM}   Checking: ${try_url}${NC}" >&2

        local response=$(curl -s -w "\n%{http_code}" -L -A "$USER_AGENT" "$try_url" 2>/dev/null)
        local status=$(echo "$response" | tail -n1)
        local content=$(echo "$response" | sed '$d')

        if [ "$status" = "200" ] && [ -n "$content" ]; then
            if echo "$content" | head -5 | grep -qiE '^#|^\*|^-|^[a-zA-Z]'; then
                echo -e "${GREEN}✓ Found: ${try_url}${NC}" >&2
                echo "$content"
                return 0
            fi
        fi
    done

    return 1
}

# Fetch entire site with sitefetch
fetch_with_sitefetch() {
    local url="$1"
    local output_file="$2"
    local pattern=$(get_docs_pattern "$url")

    echo -e "${BLUE}→ Crawling site with sitefetch...${NC}" >&2
    echo -e "${DIM}   Pattern: ${pattern}${NC}" >&2
    echo -e "${DIM}   This may take a while...${NC}" >&2

    # Create temp file
    local temp_file=$(mktemp)

    # Run sitefetch with match pattern (suppress output)
    if sitefetch "$url" -m "$pattern" -o "$temp_file" --concurrency 5 >/dev/null 2>&1; then
        if [ -s "$temp_file" ]; then
            # Remove sitefetch INFO lines from content
            grep -v "^INFO " "$temp_file" | grep -v "^WARN " | grep -v "^ERROR "
            rm -f "$temp_file"
            return 0
        fi
    fi

    rm -f "$temp_file"
    return 1
}

# Fetch via Jina Reader API (single page)
fetch_via_jina() {
    local url="$1"

    echo -e "${BLUE}→ Fetching single page via Jina Reader...${NC}" >&2

    local jina_url="${JINA_BASE}/${url}"
    local response=$(curl -s -w "\n%{http_code}" -L -A "$USER_AGENT" \
        -H "Accept: text/markdown" \
        "$jina_url" 2>/dev/null)

    local status=$(echo "$response" | tail -n1)
    local content=$(echo "$response" | sed '$d')

    if [ "$status" = "200" ] && [ -n "$content" ]; then
        echo -e "${GREEN}✓ Fetched via Jina Reader${NC}" >&2
        echo "$content"
        return 0
    else
        echo -e "${RED}✗ Jina Reader failed (HTTP $status)${NC}" >&2
        return 1
    fi
}

# Add metadata header to content
add_metadata() {
    local content="$1"
    local source_url="$2"
    local fetch_method="$3"

    cat << EOF
---
source: ${source_url}
fetched: $(date -u +"%Y-%m-%dT%H:%M:%SZ")
method: ${fetch_method}
---

${content}
EOF
}

# Main
echo -e "${YELLOW}📚 Fetching documentation...${NC}"
echo ""

# Detect and normalize URL (handle versions)
NORMALIZED_URL=$(detect_base_url "$URL")
if [ "$NORMALIZED_URL" != "$URL" ]; then
    echo -e "${YELLOW}ℹ Normalized URL: ${NORMALIZED_URL}${NC}"
fi

# Generate output filename
if [ -z "$OUTPUT_NAME" ]; then
    OUTPUT_NAME=$(extract_name "$NORMALIZED_URL")
fi
OUTPUT_FILE="${DOCS_DIR}/${OUTPUT_NAME}.md"

echo -e "Target: ${BLUE}${OUTPUT_FILE}${NC}"
echo ""

# Try fetching methods in order
CONTENT=""
METHOD=""

# 1. Try llms.txt
echo -e "${YELLOW}Step 1: Looking for llms.txt...${NC}"
if CONTENT=$(fetch_llms_txt "$NORMALIZED_URL"); then
    METHOD="llms.txt"
else
    echo -e "${DIM}   No llms.txt found${NC}"
    echo ""

    # 2. Try sitefetch for full site crawling (unless --single)
    if [ "$SINGLE_PAGE" = false ] && command -v sitefetch &> /dev/null; then
        echo -e "${YELLOW}Step 2: Crawling entire documentation site...${NC}"
        if CONTENT=$(fetch_with_sitefetch "$NORMALIZED_URL" "$OUTPUT_FILE"); then
            METHOD="sitefetch"
        else
            echo -e "${DIM}   sitefetch failed, trying Jina Reader...${NC}"
            echo ""
        fi
    elif [ "$SINGLE_PAGE" = false ]; then
        echo -e "${YELLOW}Step 2: sitefetch not installed, skipping full crawl${NC}"
        echo -e "${DIM}   Install with: bun install -g sitefetch${NC}"
        echo ""
    fi

    # 3. Fall back to Jina Reader (single page)
    if [ -z "$METHOD" ]; then
        echo -e "${YELLOW}Step 3: Fetching single page via Jina Reader...${NC}"
        if CONTENT=$(fetch_via_jina "$NORMALIZED_URL"); then
            METHOD="jina-reader"
            echo ""
            echo -e "${YELLOW}⚠ Note: Only fetched single page. For full docs:${NC}"
            echo -e "${DIM}   1. Install sitefetch: bun install -g sitefetch${NC}"
            echo -e "${DIM}   2. Re-run this command${NC}"
        else
            echo ""
            echo -e "${RED}❌ Failed to fetch documentation${NC}"
            echo ""
            echo "Possible issues:"
            echo "  - URL might be incorrect"
            echo "  - Site might block automated requests"
            echo ""
            echo "Manual options:"
            echo "  1. Install sitefetch: bun install -g sitefetch"
            echo "  2. Check if site has llms.txt: ${NORMALIZED_URL}/llms.txt"
            exit 1
        fi
    fi
fi

# Save with metadata
echo ""
echo -e "${YELLOW}Saving documentation...${NC}"
add_metadata "$CONTENT" "$NORMALIZED_URL" "$METHOD" > "$OUTPUT_FILE"

# Show stats
LINES=$(wc -l < "$OUTPUT_FILE" | tr -d ' ')
SIZE=$(du -h "$OUTPUT_FILE" | cut -f1)

echo ""
echo -e "${GREEN}✅ Documentation saved!${NC}"
echo ""
echo "  File:   ${OUTPUT_FILE}"
echo "  Size:   ${SIZE}"
echo "  Lines:  ${LINES}"
echo "  Method: ${METHOD}"
echo ""

if [ "$METHOD" = "jina-reader" ]; then
    echo -e "${YELLOW}Tip: For complete docs, install sitefetch and re-run${NC}"
elif [ "$METHOD" = "sitefetch" ]; then
    echo -e "${BLUE}Tip: Review and clean up the file if needed${NC}"
fi
