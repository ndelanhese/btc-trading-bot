# 🤖 BTC Trading Bot

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-blue.svg)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![API Status](https://img.shields.io/badge/API-REST%20%2B%20WebSocket-orange.svg)](https://github.com/yourusername/btc-trading-bot)
[![LN Markets](https://img.shields.io/badge/LN%20Markets-API%20Ready-green.svg)](https://docs.lnmarkets.com/api/)
[![Trading Bot](https://img.shields.io/badge/Trading%20Bot-Automated%20BTC-blue.svg)](https://github.com/yourusername/btc-trading-bot)

> **Automated Bitcoin trading bot with advanced risk management, DCA strategies, and real-time price monitoring for LN Markets platform.**

## Features

### 🔐 Authentication
- User registration and login
- JWT authentication
- Route protection

### 🤖 Trading Automations

#### 1. Margin Protection
- Automatic liquidation monitoring
- Activation when price reaches 5% of liquidation
- Automatic increase of liquidation distance to 10%

#### 2. Automatic Take Profit
- Daily take profit updates (1%)
- Automatic calculation considering fees
- Immediate application to all orders

#### 3. Entry Automation (DCA)
- Dollar Cost Averaging system
- 9 orders with $10 each
- $50 intervals between orders
- Initial price: $116,000
- Take profit of 0.25% per order
- 10x leverage

#### 4. Price Alert
- Real-time price monitoring
- Alerts when price exits $100,000 - $120,000 range
- Check every 1 minute

## Technologies

- **Backend**: Go 1.24+
- **Database**: PostgreSQL
- **ORM**: SQLx
- **API**: REST + WebSocket
- **Authentication**: JWT
- **Trading**: LN Markets API

## Project Structure

```
btc-trading-bot/
├── cmd/api/           # Application entry point
├── internal/          # Internal application code
│   ├── database/     # Database configuration and migrations
│   ├── handlers/     # HTTP handlers
│   ├── models/       # Data models
│   └── services/     # Business logic
├── pkg/              # Reusable packages
│   ├── lnmarkets/    # LN Markets API client
│   └── websocket/    # WebSocket client
├── go.mod
├── go.sum
└── README.md
```

## Configuration

### 1. Database

Create a PostgreSQL database and configure environment variables:

#### Option A: Interactive Setup
```bash
./scripts/setup-env.sh
```

#### Option B: Manual Setup
```bash
# Copy the example file
cp .env.example .env

# Edit the .env file with your settings
nano .env
```

#### Required Environment Variables:
```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=btc_trading_bot
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-secret-key-change-this

# Server Configuration
PORT=8080
```

### 2. Installation

```bash
# Clone repository
git clone <repository-url>
cd btc-trading-bot

# Install dependencies
go mod tidy

# Setup environment (interactive)
./scripts/setup-env.sh

# Or manually copy and configure environment
cp env.example .env
# Edit .env with your settings
```

### 3. Run

```bash
go run cmd/api/main.go
```

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register user
- `POST /api/auth/login` - Login

### Configurations (Protected)
- `POST /api/lnmarkets/config` - Configure LN Markets API
- `GET /api/lnmarkets/config` - Get LN Markets configuration

### Trading (Protected)
- `POST /api/trading/margin-protection` - Configure margin protection
- `GET /api/trading/margin-protection` - Get configuration
- `POST /api/trading/take-profit` - Configure take profit
- `GET /api/trading/take-profit` - Get configuration
- `POST /api/trading/entry-automation` - Configure entry automation
- `GET /api/trading/entry-automation` - Get configuration
- `POST /api/trading/price-alert` - Configure price alert
- `GET /api/trading/price-alert` - Get configuration
- `GET /api/trading/orders` - List orders
- `POST /api/trading/bot/start` - Start bot
- `POST /api/trading/bot/stop` - Stop bot

### Health Check
- `GET /health` - API status

## Usage Examples

### 1. Register User

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "trader1",
    "password": "password123",
    "email": "trader@example.com"
  }'
```

### 2. Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "trader1",
    "password": "password123"
  }'
```

### 3. Configure LN Markets

```bash
curl -X POST http://localhost:8080/api/lnmarkets/config \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "api_key": "your_api_key",
    "secret_key": "your_secret_key",
    "passphrase": "your_passphrase",
    "is_testnet": true
  }'
```

### 4. Configure Margin Protection

```bash
curl -X POST http://localhost:8080/api/trading/margin-protection \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "is_enabled": true,
    "activation_distance": 5.0,
    "new_liquidation_distance": 10.0
  }'
```

### 5. Configure Entry Automation

```bash
curl -X POST http://localhost:8080/api/trading/entry-automation \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
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
  }'
```

### 6. Start Bot

```bash
curl -X POST http://localhost:8080/api/trading/bot/start \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Security

- Passwords are hashed with bcrypt
- JWT tokens for authentication
- Rate limiting implemented
- Input validation on all endpoints
- HTTPS connections recommended in production

## Monitoring

The bot logs detailed information about:
- Price updates
- Protection activations
- Order creation
- Price alerts
- API errors

## Development

### Data Structure

The database includes tables for:
- Users
- LN Markets configurations
- Margin protection
- Take profit
- Entry automation
- Price alerts
- Trading orders

### Data Flow

1. **WebSocket** receives real-time price updates
2. **Trading Service** processes each update
3. **Automations** are executed based on configurations
4. **LN Markets API** executes orders when needed
5. **Database** stores history and configurations

## Contributing

1. Fork the project
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## License

This project is licensed under the MIT License.
