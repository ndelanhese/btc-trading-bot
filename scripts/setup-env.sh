#!/bin/bash

# BTC Trading Bot - Environment Setup Script
# This script helps you set up your environment variables

echo "🤖 BTC Trading Bot - Environment Setup"
echo "======================================"
echo ""

# Check if .env file already exists
if [ -f ".env" ]; then
    echo "⚠️  .env file already exists!"
    read -p "Do you want to overwrite it? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Setup cancelled."
        exit 1
    fi
fi

# Copy the example file
if [ -f ".env.example" ]; then
    cp .env.example .env
    echo "✅ Created .env file from .env.example"
else
    echo "❌ .env.example file not found!"
    exit 1
fi

echo ""
echo "🔧 Now let's configure your environment variables:"
echo ""

# Database Configuration
echo "📊 Database Configuration:"
read -p "Database Host (default: localhost): " db_host
db_host=${db_host:-localhost}

read -p "Database Port (default: 5432): " db_port
db_port=${db_port:-5432}

read -p "Database User (default: postgres): " db_user
db_user=${db_user:-postgres}

read -s -p "Database Password: " db_password
echo

read -p "Database Name (default: btc_trading_bot): " db_name
db_name=${db_name:-btc_trading_bot}

# JWT Secret
echo ""
echo "🔐 JWT Configuration:"
read -s -p "JWT Secret (leave empty for auto-generation): " jwt_secret
echo

if [ -z "$jwt_secret" ]; then
    jwt_secret=$(openssl rand -base64 32)
    echo "✅ Auto-generated JWT secret"
fi

# Server Configuration
echo ""
echo "🌐 Server Configuration:"
read -p "Server Port (default: 8080): " port
port=${port:-8080}

# Update .env file
echo ""
echo "📝 Updating .env file..."

# Use sed to update the .env file
sed -i "s/DB_HOST=.*/DB_HOST=$db_host/" .env
sed -i "s/DB_PORT=.*/DB_PORT=$db_port/" .env
sed -i "s/DB_USER=.*/DB_USER=$db_user/" .env
sed -i "s/DB_PASSWORD=.*/DB_PASSWORD=$db_password/" .env
sed -i "s/DB_NAME=.*/DB_NAME=$db_name/" .env
sed -i "s/JWT_SECRET=.*/JWT_SECRET=$jwt_secret/" .env
sed -i "s/PORT=.*/PORT=$port/" .env

echo "✅ Environment variables configured!"
echo ""

# Show summary
echo "📋 Configuration Summary:"
echo "   Database Host: $db_host"
echo "   Database Port: $db_port"
echo "   Database User: $db_user"
echo "   Database Name: $db_name"
echo "   Server Port: $port"
echo "   JWT Secret: [HIDDEN]"
echo ""

echo "🚀 Next steps:"
echo "1. Start your PostgreSQL database"
echo "2. Run: go run cmd/api/main.go"
echo "3. Test with: ./scripts/test-api.sh"
echo ""

echo "⚠️  Important Security Notes:"
echo "- Keep your .env file secure and never commit it to version control"
echo "- Use strong passwords in production"
echo "- Consider using a secrets manager for production environments"
echo ""

echo "✅ Environment setup complete!"
