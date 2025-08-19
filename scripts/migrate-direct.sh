#!/bin/bash

# BTC Trading Bot - Direct Database Migration Script
# This script runs database migrations directly using go run

set -e  # Exit on any error

echo "🔄 BTC Trading Bot - Running Database Migrations (Direct)"
echo "========================================================="
echo ""

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Error: go.mod not found. Please run this script from the project root directory."
    exit 1
fi

# Check if .env file exists
if [ ! -f ".env" ]; then
    echo "⚠️  Warning: .env file not found. Using system environment variables."
    echo "   Make sure your database environment variables are set:"
    echo "   - DB_HOST"
    echo "   - DB_PORT" 
    echo "   - DB_USER"
    echo "   - DB_PASSWORD"
    echo "   - DB_NAME"
    echo "   - DB_SSLMODE"
    echo ""
fi

# Run migrations directly
echo "🚀 Running database migrations..."
if ! go run cmd/migrate/main.go; then
    echo "❌ Migration failed"
    exit 1
fi

echo ""
echo "✅ Database migrations completed successfully!"
echo "🎉 All done!"
