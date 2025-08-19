#!/bin/bash

# BTC Trading Bot - Docker Migration Script
# This script runs database migrations in Docker environment

set -e  # Exit on any error

echo "🔄 BTC Trading Bot - Running Database Migrations (Docker)"
echo "========================================================"
echo ""

# Debug: Show current directory and contents
echo "📁 Current directory: $(pwd)"
echo "📁 Directory contents:"
ls -la
echo ""

# Find the script's directory and navigate to the parent directory where binaries are located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
echo "📁 Script directory: $SCRIPT_DIR"

# Navigate to the parent directory of scripts (where binaries should be)
cd "$SCRIPT_DIR/.."
echo "📁 Working directory after navigation: $(pwd)"
echo "📁 Directory contents after navigation:"
ls -la
echo ""

# Check if we're in the right directory (look for binaries in current directory)
if [ ! -f "./btc-trading-bot" ] && [ ! -f "./trading-bot" ]; then
    echo "❌ Error: btc-trading-bot binary not found in current directory."
    echo "   Expected location: ./btc-trading-bot or ./trading-bot"
    echo "   Current directory: $(pwd)"
    exit 1
fi

# Check if migration binary exists, if not build it
if [ ! -f "./migrate" ]; then
    echo "📦 migrate binary not found, building it..."
    
    # Check if Go is available
    if ! command -v go &> /dev/null; then
        echo "❌ Error: Go is not installed or not in PATH"
        echo "   Please install Go to build the migrate binary"
        exit 1
    fi
    
    # Build the migrate binary
    echo "🔨 Building migrate binary from cmd/migrate/main.go..."
    if ! go build -o migrate cmd/migrate/main.go; then
        echo "❌ Failed to build migrate binary"
        exit 1
    fi
    
    echo "✅ migrate binary built successfully"
else
    echo "✅ migrate binary found"
fi

# Run migrations using the built binary
echo "🚀 Running database migrations..."
if ! ./migrate; then
    echo "❌ Migration failed"
    exit 1
fi

echo ""
echo "✅ Database migrations completed successfully!"
echo "🎉 All done!"
