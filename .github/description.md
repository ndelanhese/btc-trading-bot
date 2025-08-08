# 🤖 BTC Trading Bot - Automated Bitcoin Trading on LN Markets

> **Automated Bitcoin trading bot with advanced risk management, DCA strategies, and real-time price monitoring for LN Markets platform.**

## 🚀 Features

### 🔐 **Authentication & Security**
- User registration and JWT authentication
- Secure password hashing with bcrypt
- Protected API endpoints
- Rate limiting and CORS support

### 🤖 **Trading Automations**

#### **🛡️ Margin Protection**
- Real-time liquidation monitoring
- Automatic activation at 5% distance from liquidation
- Dynamic adjustment to 10% liquidation distance
- Prevents margin calls automatically

#### **💰 Take Profit Automation**
- Daily 1% take profit updates
- Automatic fee calculation (entry/exit + daily rates)
- Immediate application to all open orders
- Smart profit capture strategy

#### **📈 DCA (Dollar Cost Averaging)**
- 9 strategic orders with $10 each
- $50 price intervals between entries
- Starting price: $116,000
- 0.25% take profit per order
- 10x leverage configuration

#### **🔔 Price Alerts**
- Real-time Bitcoin price monitoring
- Alerts when price exits $100k-$120k range
- Configurable check intervals
- WebSocket-based price updates

## 🛠️ Tech Stack

- **Backend**: Go 1.24+
- **Database**: PostgreSQL with SQLx ORM
- **API**: REST + WebSocket
- **Authentication**: JWT
- **Trading**: LN Markets API
- **Containerization**: Docker + Docker Compose
- **Real-time**: WebSocket connections

## 📊 Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   WebSocket     │    │   Trading       │    │   LN Markets    │
│   Price Feed    │───▶│   Service       │───▶│   API           │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   PostgreSQL    │    │   JWT Auth      │    │   Risk Mgmt     │
│   Database      │    │   Service       │    │   Engine        │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🚀 Quick Start

### 1. **Clone & Setup**
```bash
git clone https://github.com/yourusername/btc-trading-bot.git
cd btc-trading-bot
go mod tidy
```

### 2. **Start Database**
```bash
docker compose up -d postgres
```

### 3. **Run Application**
```bash
go run cmd/api/main.go
```

### 4. **Test API**
```bash
./scripts/test-api.sh
```

## 📋 API Endpoints

### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login

### Trading Configuration
- `POST /api/lnmarkets/config` - Configure LN Markets API
- `POST /api/trading/margin-protection` - Set margin protection
- `POST /api/trading/take-profit` - Configure take profit
- `POST /api/trading/entry-automation` - Set DCA strategy
- `POST /api/trading/price-alert` - Configure price alerts

### Bot Control
- `POST /api/trading/bot/start` - Start trading bot
- `POST /api/trading/bot/stop` - Stop trading bot
- `GET /api/trading/orders` - List trading orders

## 🔧 Configuration

### Environment Variables
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=btc_trading_bot
JWT_SECRET=your-secret-key
PORT=8080
```

### Default Trading Settings
- **Margin Protection**: 5% → 10% activation
- **Take Profit**: 1% daily increase
- **DCA Strategy**: 9 orders × $10, $50 intervals
- **Price Alerts**: $100k - $120k range
- **Leverage**: 10x

## 🎯 Use Cases

### **Conservative Trading**
- Low leverage (5x-10x)
- Wide margin protection
- Small DCA amounts
- Frequent take profits

### **Aggressive Trading**
- Higher leverage (10x-20x)
- Tighter margin protection
- Larger DCA amounts
- Longer-term holds

### **Risk Management**
- Automatic liquidation prevention
- Dynamic take profit adjustment
- Real-time price monitoring
- Configurable alert thresholds

## 📈 Performance Features

- **Real-time Processing**: WebSocket price feeds
- **Concurrent Operations**: Goroutine-based architecture
- **Database Optimization**: Indexed queries with SQLx
- **Error Handling**: Comprehensive error management
- **Logging**: Detailed operation logs

## 🔒 Security Features

- **Password Hashing**: bcrypt with salt
- **JWT Authentication**: Secure token-based auth
- **Input Validation**: Request sanitization
- **Rate Limiting**: API protection
- **CORS Support**: Cross-origin security

## 🐳 Docker Support

```bash
# Full stack with Docker Compose
docker compose up -d

# Build custom image
docker build -t btc-trading-bot .
docker run -p 8080:8080 btc-trading-bot
```

## 📝 Development

### Project Structure
```
btc-trading-bot/
├── cmd/api/           # Application entry point
├── internal/          # Internal application code
│   ├── database/     # Database configuration
│   ├── handlers/     # HTTP handlers
│   ├── models/       # Data models
│   └── services/     # Business logic
├── pkg/              # Reusable packages
│   ├── lnmarkets/    # LN Markets API client
│   └── websocket/    # WebSocket client
└── scripts/          # Utility scripts
```

### Testing
```bash
# Run API tests
./scripts/test-api.sh

# Manual testing
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"pass","email":"test@example.com"}'
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⚠️ Disclaimer

This software is for educational and research purposes. Trading cryptocurrencies involves substantial risk of loss. Use at your own risk. The authors are not responsible for any financial losses.

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/yourusername/btc-trading-bot/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/btc-trading-bot/discussions)
- **Documentation**: [Wiki](https://github.com/yourusername/btc-trading-bot/wiki)

---

**⭐ Star this repository if you find it useful!**

**🔗 Built with ❤️ using Go and LN Markets API**
