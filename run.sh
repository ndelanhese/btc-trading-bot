#!/bin/bash

echo "Starting BTC Trading Bot with Air (hot reloading)..."
echo "Press Ctrl+C to stop"

# Check if .env file exists, if not create a basic one
if [ ! -f .env ]; then
    echo "Creating basic .env file..."
    cat > .env << EOF
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=btc_trading_bot

# JWT Configuration
JWT_SECRET=your-secret-key-here

# Server Configuration
PORT=8080

# LN Markets Configuration (optional)
LN_MARKETS_API_KEY=
LN_MARKETS_SECRET_KEY=
LN_MARKETS_PASSPHRASE=
LN_MARKETS_IS_TESTNET=true
EOF
    echo ".env file created with default values. Please update with your actual configuration."
fi

# Run with Air
air
