#!/bin/bash

# Default values
DEFAULT_URL="https://example.com"
DEFAULT_CONCURRENCY=10
DEFAULT_MAX_PAGES=100

# Get arguments or use defaults
URL=${1:-$DEFAULT_URL}
CONCURRENCY=${2:-$DEFAULT_CONCURRENCY}
MAX_PAGES=${3:-$DEFAULT_MAX_PAGES}

# Build and run
echo "Crawling $URL with concurrency $CONCURRENCY and max pages $MAX_PAGES"
make run ARGS="$URL $CONCURRENCY $MAX_PAGES" 