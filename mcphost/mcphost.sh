#!/bin/bash

set -e

cd /Users/cuenobi/Documents/Dev/mcp-platform/mcphost || exit 1
if [ ! -f "./mcphost" ]; then
  echo "🔧 Building mcphost binary..."
  go build -o mcphost
fi

./mcphost jira sync --project "$1"