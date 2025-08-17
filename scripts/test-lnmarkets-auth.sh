#!/bin/bash

# Script to test LN Markets authentication

API_URL="http://localhost:8080"

echo "🔐 Testing LN Markets Authentication"
echo "====================================="

# Check if the bot is running
echo "1. Checking if the bot is running..."
if ! curl -s "$API_URL/health" > /dev/null; then
    echo "❌ Bot is not running. Please start the bot first."
    exit 1
fi
echo "✅ Bot is running"
echo ""

# Register user if not exists
echo "2. Registering user..."
REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "lnmarkets_test",
    "password": "password123",
    "email": "lnmarkets@test.com"
  }')

echo "$REGISTER_RESPONSE" | jq .
echo ""

# Login
echo "3. Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "lnmarkets_test",
    "password": "password123"
  }')

echo "$LOGIN_RESPONSE" | jq .

# Extract token
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')

if [ "$TOKEN" = "null" ] || [ "$TOKEN" = "" ]; then
    echo "❌ Login failed"
    exit 1
fi

echo "✅ Login successful"
echo ""

# Configure LN Markets with test credentials
echo "4. Configuring LN Markets..."
curl -s -X POST "$API_URL/api/lnmarkets/config" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "api_key": "your_test_api_key_here",
    "secret_key": "your_test_secret_key_here",
    "passphrase": "your_test_passphrase_here",
    "is_testnet": true
  }' | jq .
echo ""

echo "⚠️  IMPORTANT: Replace the API credentials above with your actual test credentials"
echo ""

# Test account balance (this will show the debug logs)
echo "5. Testing account balance (with debug logs)..."
curl -s -X GET "$API_URL/api/trading/account-balance" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# Test positions
echo "6. Testing positions..."
curl -s -X GET "$API_URL/api/trading/positions" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

echo "✅ Authentication test completed!"
echo ""
echo "To enable debug logs, set the environment variable:"
echo "export DEBUG_LNMARKETS=true"
echo ""
echo "Then restart the bot and run this test again to see detailed signature information."
echo ""
echo "Check the bot logs for debug information about the signature generation."
echo "If you see authentication errors, verify your API credentials and check the debug logs."
