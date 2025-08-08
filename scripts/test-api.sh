#!/bin/bash

# Script to test the BTC Trading Bot API

API_URL="http://localhost:8080"

echo "🚀 Testing BTC Trading Bot API"
echo "=================================="

# Testing health check
echo "1. Testing health check..."
curl -s "$API_URL/health" | jq .
echo ""

# Registering user
echo "2. Registering user..."
REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "email": "test@example.com"
  }')

echo "$REGISTER_RESPONSE" | jq .
echo ""

# Logging in
echo "3. Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }')

echo "$LOGIN_RESPONSE" | jq .

# Extracting token
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')

if [ "$TOKEN" = "null" ] || [ "$TOKEN" = "" ]; then
    echo "❌ Login failed"
    exit 1
fi

echo "✅ Login successful"
echo ""

# Configuring LN Markets (example)
echo "4. Configuring LN Markets..."
curl -s -X POST "$API_URL/api/lnmarkets/config" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "api_key": "test_api_key",
    "secret_key": "test_secret_key",
    "passphrase": "test_passphrase",
    "is_testnet": true
  }' | jq .
echo ""

# Configuring margin protection
echo "5. Configuring margin protection..."
curl -s -X POST "$API_URL/api/trading/margin-protection" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "is_enabled": true,
    "activation_distance": 5.0,
    "new_liquidation_distance": 10.0
  }' | jq .
echo ""

# Configuring take profit
echo "6. Configuring take profit..."
curl -s -X POST "$API_URL/api/trading/take-profit" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "is_enabled": true,
    "daily_percentage": 1.0
  }' | jq .
echo ""

# Configuring entry automation
echo "7. Configuring entry automation..."
curl -s -X POST "$API_URL/api/trading/entry-automation" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "is_enabled": true,
    "amount_per_order": 10.0,
    "margin_per_order": 855,
    "number_of_orders": 9,
    "price_variation": 50.0,
    "initial_price": 116000.0,
    "take_profit_per_order": 0.25,
    "operation_type": "buy",
    "leverage": 10
  }' | jq .
echo ""

# Configuring price alert
echo "8. Configuring price alert..."
curl -s -X POST "$API_URL/api/trading/price-alert" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "is_enabled": true,
    "min_price": 100000.0,
    "max_price": 120000.0,
    "check_interval": 60
  }' | jq .
echo ""

# Listing configurations
echo "9. Listing configurations..."
echo "Margin protection:"
curl -s -X GET "$API_URL/api/trading/margin-protection" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

echo "Take profit:"
curl -s -X GET "$API_URL/api/trading/take-profit" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

echo "Entry automation:"
curl -s -X GET "$API_URL/api/trading/entry-automation" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

echo "Price alert:"
curl -s -X GET "$API_URL/api/trading/price-alert" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

echo "✅ All tests completed successfully!"
echo ""
echo "To start the trading bot, execute:"
echo "curl -X POST $API_URL/api/trading/bot/start -H \"Authorization: Bearer $TOKEN\""
