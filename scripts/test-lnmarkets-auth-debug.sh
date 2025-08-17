#!/bin/bash

# Script to test LN Markets authentication with debug enabled

echo "🔐 Testing LN Markets Authentication with Debug"
echo "==============================================="

# Set debug environment variable
export DEBUG_LNMARKETS=true

echo "✅ Debug mode enabled (DEBUG_LNMARKETS=true)"
echo ""

# Run the authentication test
./scripts/test-lnmarkets-auth.sh

echo ""
echo "🔍 Debug information should be visible in the bot logs above."
echo "Look for lines starting with 'DEBUG:' to see the signature generation details."


