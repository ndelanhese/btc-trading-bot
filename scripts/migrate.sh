#!/bin/bash

# BTC Trading Bot - Database Migration Script
# This script runs database migrations for production environment

set -e  # Exit on any error

echo "🔄 BTC Trading Bot - Running Database Migrations"
echo "================================================"
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

# Build the migration binary
echo "🔨 Building migration binary..."
if ! go build -o bin/migrate cmd/migrate/main.go; then
    echo "❌ Failed to build migration binary"
    exit 1
fi

echo "✅ Migration binary built successfully"
echo ""

# Run migrations
echo "🚀 Running database migrations..."
if ! ./bin/migrate; then
    echo "❌ Migration failed"
    exit 1
fi

echo ""
echo "✅ Database migrations completed successfully!"
echo ""

# Clean up
echo "🧹 Cleaning up..."
rm -f bin/migrate

echo "🎉 All done!"
