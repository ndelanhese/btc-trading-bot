#!/bin/bash

# Test script for BTC Trading Bot
# This script tests the main API endpoints

BASE_URL="http://localhost:8080"
API_TOKEN=""

echo "🚀 Testing BTC Trading Bot API"
echo "================================"

# Test health endpoint
echo "1. Testing health endpoint..."
curl -s "$BASE_URL/health" | jq '.'

echo -e "\n2. Testing user registration..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "testpassword123"
  }')

echo $REGISTER_RESPONSE | jq '.'

echo -e "\n3. Testing user login..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "testpassword123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
echo "Token: $TOKEN"

if [ "$TOKEN" != "null" ] && [ "$TOKEN" != "" ]; then
    echo -e "\n4. Testing LN Markets configuration..."
    curl -s -X POST "$BASE_URL/api/lnmarkets/config" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "api_key": "your-api-key",
        "secret_key": "your-secret-key",
        "passphrase": "your-passphrase",
        "is_testnet": true
      }' | jq '.'

    echo -e "\n5. Testing margin protection configuration..."
    curl -s -X POST "$BASE_URL/api/trading/margin-protection" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "is_enabled": true,
        "activation_distance": 5.0,
        "new_liquidation_distance": 10.0
      }' | jq '.'

    echo -e "\n6. Testing take profit configuration..."
    curl -s -X POST "$BASE_URL/api/trading/take-profit" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "is_enabled": true,
        "daily_percentage": 1.0
      }' | jq '.'

    echo -e "\n7. Testing entry automation configuration..."
    curl -s -X POST "$BASE_URL/api/trading/entry-automation" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
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
      }' | jq '.'

    echo -e "\n8. Testing price alert configuration..."
    curl -s -X POST "$BASE_URL/api/trading/price-alert" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "is_enabled": true,
        "min_price": 100000.0,
        "max_price": 120000.0,
        "check_interval": 60
      }' | jq '.'

    echo -e "\n9. Testing bot status..."
    curl -s -X GET "$BASE_URL/api/trading/bot/status" \
      -H "Authorization: Bearer $TOKEN" | jq '.'

    echo -e "\n10. Testing bot start..."
    curl -s -X POST "$BASE_URL/api/trading/bot/start" \
      -H "Authorization: Bearer $TOKEN" | jq '.'

    echo -e "\n11. Testing bot status after start..."
    sleep 2
    curl -s -X GET "$BASE_URL/api/trading/bot/status" \
      -H "Authorization: Bearer $TOKEN" | jq '.'

    echo -e "\n12. Testing bot stop..."
    curl -s -X POST "$BASE_URL/api/trading/bot/stop" \
      -H "Authorization: Bearer $TOKEN" | jq '.'

    echo -e "\n13. Testing bot status after stop..."
    curl -s -X GET "$BASE_URL/api/trading/bot/status" \
      -H "Authorization: Bearer $TOKEN" | jq '.'

else
    echo "❌ Login failed, skipping protected endpoints"
fi

echo -e "\n✅ Test completed!"
