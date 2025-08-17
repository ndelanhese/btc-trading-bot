# BTC Trading Bot

A sophisticated Bitcoin trading bot that integrates with LN Markets API for automated trading with advanced features like margin protection, take-profit automation, entry automation, and price alerts.

## 🚀 Features

### Core Trading Features
- **Real-time Price Monitoring**: WebSocket connection to LN Markets for live price updates
- **Margin Protection**: Automatically adjusts take-profit when positions are close to liquidation
- **Take Profit Automation**: Daily percentage-based take-profit adjustments
- **Entry Automation**: DCA (Dollar Cost Averaging) with configurable parameters
- **Price Alerts**: Custom price range monitoring with configurable intervals

### API Features
- **User Authentication**: JWT-based authentication system
- **Configuration Management**: Store and manage LN Markets API credentials
- **Bot Management**: Start, stop, and monitor trading bots
- **Position Management**: View, close, and update positions
- **Account Balance**: Real-time account balance monitoring

## 🛠️ Setup

### Prerequisites
- Go 1.19 or higher
- PostgreSQL database
- LN Markets API credentials

### Installation

1. **Clone the repository**
```bash
git clone <repository-url>
cd btc-trading-bot
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up environment variables**
```bash
cp .env.example .env
```

Edit `.env` file:
```env
# Database
DATABASE_URL=postgres://username:password@localhost:5432/btc_trading_bot?sslmode=disable

# JWT
JWT_SECRET=your-secret-key-change-this

# Server
PORT=8080

# LN Markets (optional - can be set via API)
LN_MARKETS_API_KEY=your-api-key
LN_MARKETS_SECRET_KEY=your-secret-key
LN_MARKETS_PASSPHRASE=your-passphrase
LN_MARKETS_TESTNET=true
```

4. **Set up database**
```bash
# Create database
createdb btc_trading_bot

# Run migrations (automatic on startup)
go run main.go
```

5. **Run the application**
```bash
go run main.go

# Or

# Use air for dev mode with watch (air installation: https://github.com/air-verse/air)
air
```

## 📚 API Documentation

### Authentication

#### Register User
```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "your_username",
  "email": "your_email@example.com",
  "password": "your_password"
}
```

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "your_username",
  "password": "your_password"
}
```

### LN Markets Configuration

#### Set Configuration
```http
POST /api/lnmarkets/config
Authorization: Bearer <token>
Content-Type: application/json

{
  "api_key": "your-api-key",
  "secret_key": "your-secret-key",
  "passphrase": "your-passphrase",
  "is_testnet": true
}
```

#### Get Configuration
```http
GET /api/lnmarkets/config
Authorization: Bearer <token>
```

### Trading Configuration

#### Margin Protection
```http
POST /api/trading/margin-protection
Authorization: Bearer <token>
Content-Type: application/json

{
  "is_enabled": true,
  "activation_distance": 5.0,
  "new_liquidation_distance": 10.0
}
```

#### Take Profit
```http
POST /api/trading/take-profit
Authorization: Bearer <token>
Content-Type: application/json

{
  "is_enabled": true,
  "daily_percentage": 1.0
}
```

#### Entry Automation
```http
POST /api/trading/entry-automation
Authorization: Bearer <token>
Content-Type: application/json

{
  "is_enabled": true,
  "amount_per_order": 10.0,
  "margin_per_order": 855,
  "number_of_orders": 9,
  "price_variation": 50.0,
  "initial_price": 116000.0,
  "take_profit_per_order": 0.25,
  "operation_type": "buy",
  "leverage": 10
}
```

#### Price Alert
```http
POST /api/trading/price-alert
Authorization: Bearer <token>
Content-Type: application/json

{
  "is_enabled": true,
  "min_price": 100000.0,
  "max_price": 120000.0,
  "check_interval": 60
}
```

### Bot Management

#### Start Bot
```http
POST /api/trading/bot/start
Authorization: Bearer <token>
```

#### Stop Bot
```http
POST /api/trading/bot/stop
Authorization: Bearer <token>
```

#### Get Bot Status
```http
GET /api/trading/bot/status
Authorization: Bearer <token>
```

### Trading Operations

#### Get Account Balance
```http
GET /api/trading/account/balance
Authorization: Bearer <token>
```

#### Get Positions
```http
GET /api/trading/positions
Authorization: Bearer <token>
```

#### Get Position
```http
GET /api/trading/positions/{id}
Authorization: Bearer <token>
```

#### Close Position
```http
POST /api/trading/positions/{id}/close
Authorization: Bearer <token>
```

#### Update Take Profit
```http
POST /api/trading/positions/{id}/take-profit
Authorization: Bearer <token>
Content-Type: application/json

{
  "price": 120000.0
}
```

#### Update Stop Loss
```http
POST /api/trading/positions/{id}/stop-loss
Authorization: Bearer <token>
Content-Type: application/json

{
  "price": 110000.0
}
```

## 🧪 Testing

Run the test script to verify all endpoints:

```bash
./scripts/test-bot.sh
```

## 🔧 Configuration

### Entry Automation Parameters
- `amount_per_order`: Amount in USD per order
- `margin_per_order`: Margin in sats per order
- `number_of_orders`: Total number of orders to place
- `price_variation`: Price difference between orders
- `initial_price`: Starting price for the first order
- `take_profit_per_order`: Take profit percentage per order
- `operation_type`: "buy" or "sell"
- `leverage`: Leverage for the positions

### Margin Protection Parameters
- `activation_distance`: Distance to liquidation to trigger protection (%)
- `new_liquidation_distance`: New distance to liquidation after protection (%)

### Take Profit Parameters
- `daily_percentage`: Daily percentage increase for take profit (%)

### Price Alert Parameters
- `min_price`: Minimum price threshold
- `max_price`: Maximum price threshold
- `check_interval`: Interval between checks (seconds)

## 🚨 Security Considerations

1. **API Keys**: Store LN Markets API credentials securely
2. **JWT Secret**: Use a strong, unique JWT secret
3. **Database**: Use strong passwords and enable SSL
4. **Environment**: Run in production with proper firewall rules
5. **Monitoring**: Implement logging and monitoring for production use

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⚠️ Disclaimer

This software is for educational and research purposes. Trading cryptocurrencies involves significant risk. Use at your own risk and never invest more than you can afford to lose.
