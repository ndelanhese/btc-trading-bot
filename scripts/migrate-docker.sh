#!/bin/bash

# BTC Trading Bot - Docker Migration Script
# This script runs database migrations in Docker environment

set -e  # Exit on any error

echo "🔄 BTC Trading Bot - Running Database Migrations (Docker)"
echo "========================================================"
echo ""

# Check if we're in the right directory
if [ ! -f "btc-trading-bot" ]; then
    echo "❌ Error: btc-trading-bot binary not found. Please run this script from the correct directory."
    exit 1
fi

# Check if migration binary exists
if [ ! -f "migrate" ]; then
    echo "❌ Error: migrate binary not found."
    exit 1
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
