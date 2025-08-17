# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview

A sophisticated Bitcoin trading bot that integrates with LN Markets API for automated trading with advanced features like margin protection, take-profit automation, entry automation (DCA), and price alerts. Built in Go with PostgreSQL and JWT authentication.

## Essential Development Commands

### Environment Setup
```bash
# Initial setup
cp .env.example .env
# Edit .env with your database and JWT credentials
go mod download

# Automated environment setup (interactive)
./scripts/setup-env.sh

# Database setup (PostgreSQL must be running)
createdb btc_trading_bot
```

### Running the Application
```bash
# Run the API server (main command)
go run cmd/api/main.go

# Build and run executable
go build -o btc-trading-bot main.go
./btc-trading-bot

# Run with Docker
docker-compose up --build
```

### Testing
```bash
# Test all API endpoints (requires running server)
./scripts/test-bot.sh

# Test basic API functionality
./scripts/test-api.sh

# Run health check
curl http://localhost:8080/health
```

### Development Commands
```bash
# Format code
go fmt ./...

# Run tests (if any exist)
go test ./...

# Check dependencies
go mod tidy
go mod verify

# Build for production
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o btc-trading-bot cmd/api/main.go
```

## Architecture Overview

### Directory Structure
```
├── cmd/api/main.go          # Application entry point, server setup
├── internal/                # Private application code
│   ├── database/           # Database connection, migrations
│   ├── handlers/           # HTTP request handlers (auth, trading)
│   ├── models/            # Data structures for DB and API
│   └── services/          # Business logic (auth, trading)
├── pkg/                   # Reusable packages
│   ├── lnmarkets/        # LN Markets API client
│   └── websocket/        # WebSocket connection handling
└── scripts/              # Utility scripts for testing and setup
```

### Core Architecture Components

**API Layer**: RESTful HTTP API built with Gorilla Mux, JWT authentication middleware, CORS enabled for cross-origin requests.

**Service Layer**: Business logic separation between `AuthService` (user management, JWT tokens) and `TradingService` (bot management, LN Markets integration).

**Data Layer**: PostgreSQL with sqlx for database operations, automatic migrations on startup, separate models for each trading configuration type.

**External Integration**: LN Markets API client with HMAC signature authentication, supports both testnet and mainnet, WebSocket client for real-time price feeds.

### Key Models and Configuration Types
- **LNMarketsConfig**: API credentials and testnet flag
- **MarginProtection**: Automatic liquidation distance adjustment
- **TakeProfit**: Daily percentage-based profit taking
- **EntryAutomation**: DCA with configurable orders, price variations, leverage
- **PriceAlert**: Custom price range monitoring with intervals
- **TradingOrder**: Order tracking and status management

## Development Patterns

### Authentication Flow
All protected endpoints require JWT Bearer token. Use `/api/auth/register` and `/api/auth/login` to obtain tokens. Middleware extracts user_id from token context for database operations.

### Configuration Management
Each user has separate configurations stored in database tables. Handlers follow upsert pattern (UPDATE if exists, INSERT if new). All configurations include user_id, timestamps, and enable/disable flags.

### Error Handling
HTTP handlers return appropriate status codes (400 for bad requests, 401 for auth errors, 404 for not found, 500 for server errors). Database errors and LN Markets API errors are wrapped and returned as HTTP errors.

### Database Operations
Uses sqlx for SQL operations with struct scanning. Migrations run automatically on startup. Connection pooling configured with max 25 open connections, 5 idle connections.

## Key Integration Points

### LN Markets API
Client supports both testnet and mainnet endpoints. Authentication uses HMAC-SHA256 signatures with timestamp, method, path, and parameters. All trading operations (positions, orders, balance) proxy through to LN Markets API.

### WebSocket Integration
Real-time price monitoring through WebSocket connections. Used for price alerts and margin protection monitoring. Client handles reconnection and heartbeat functionality.

### Trading Bot Lifecycle
Bot can be started/stopped per user. When started, initializes LN Markets client with user's credentials, begins monitoring based on enabled configurations (margin protection, take profit, price alerts, entry automation).

## Environment Requirements

- Go 1.19+ (currently using Go 1.24.5)
- PostgreSQL database (any recent version)
- LN Markets API credentials (testnet or mainnet)
- JWT secret for token signing

## Production Considerations

Database migrations run automatically on startup. Use strong JWT secrets and secure database credentials. LN Markets API keys should be configured per user through the API, not environment variables. Docker setup included with PostgreSQL container and multi-stage build for optimized binary.
