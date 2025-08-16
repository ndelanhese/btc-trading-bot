# Swagger Documentation for BTC Trading Bot API

This document provides information about the Swagger/OpenAPI documentation for the BTC Trading Bot API.

## Overview

The API documentation is available in OpenAPI 3.0 format in the `swagger.yaml` file. This documentation covers all endpoints, request/response schemas, authentication methods, and error handling for the BTC Trading Bot backend.

## API Endpoints Covered

### Authentication
- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login user and get JWT token

### Configuration Management
- `GET/POST /api/lnmarkets/config` - LN Markets configuration
- `GET/POST /api/trading/margin-protection` - Margin protection settings
- `GET/POST /api/trading/take-profit` - Take profit configuration
- `GET/POST /api/trading/entry-automation` - Entry automation settings
- `GET/POST /api/trading/price-alert` - Price alert configuration

### Trading Operations
- `GET /api/trading/orders` - Get trading orders
- `GET /api/trading/account/balance` - Get account balance
- `GET /api/trading/positions` - Get all positions
- `GET /api/trading/positions/{id}` - Get specific position
- `POST /api/trading/positions/{id}/close` - Close position
- `POST /api/trading/positions/{id}/take-profit` - Update take profit
- `POST /api/trading/positions/{id}/stop-loss` - Update stop loss

### Bot Control
- `POST /api/trading/bot/start` - Start trading bot
- `POST /api/trading/bot/stop` - Stop trading bot
- `GET /api/trading/bot/status` - Get bot status

### Health Check
- `GET /health` - API health check

## How to View the Documentation

### Option 1: Swagger UI (Recommended)

1. **Online Swagger Editor:**
   - Go to [Swagger Editor](https://editor.swagger.io/)
   - Copy the contents of `swagger.yaml` and paste it into the editor
   - The documentation will be rendered automatically

2. **Local Swagger UI:**
   ```bash
   # Install swagger-ui-cli globally
   npm install -g swagger-ui-cli
   
   # Serve the documentation locally
   swagger-ui-cli serve swagger.yaml
   ```

3. **Docker Swagger UI:**
   ```bash
   # Run Swagger UI in Docker
   docker run -p 8080:8080 -e SWAGGER_JSON=/swagger.yaml -v $(pwd):/swagger swaggerapi/swagger-ui
   ```

### Option 2: Redoc

1. **Online Redoc:**
   - Go to [Redoc](https://redocly.github.io/redoc/)
   - Copy the contents of `swagger.yaml` and paste it into the editor

2. **Local Redoc:**
   ```bash
   # Install redoc-cli globally
   npm install -g redoc-cli
   
   # Generate HTML documentation
   redoc-cli serve swagger.yaml
   ```

### Option 3: Postman

1. Import the `swagger.yaml` file into Postman
2. Postman will automatically generate a collection with all the endpoints
3. You can then test the API directly from Postman

## Authentication

The API uses JWT (JSON Web Token) authentication for protected endpoints. To authenticate:

1. Register a new user using `POST /api/auth/register`
2. Login using `POST /api/auth/login` to get a JWT token
3. Include the token in the `Authorization` header for subsequent requests:
   ```
   Authorization: Bearer <your_jwt_token>
   ```

## Data Models

The documentation includes comprehensive schemas for all data models:

- **User**: User account information
- **UserRegister/UserLogin**: Authentication request models
- **LNMarketsConfig**: LN Markets API configuration
- **MarginProtection**: Margin protection settings
- **TakeProfit**: Take profit configuration
- **EntryAutomation**: Entry automation settings
- **PriceAlert**: Price alert configuration
- **TradingOrder**: Trading order information
- **Position**: Position information
- **AccountBalance**: Account balance information
- **BotStatus**: Bot status information

## Error Handling

All endpoints include proper error responses with appropriate HTTP status codes:

- `400` - Bad Request (invalid input)
- `401` - Unauthorized (missing or invalid token)
- `404` - Not Found (resource not found)
- `500` - Internal Server Error (server error)

## Testing the API

### Using curl

```bash
# Health check
curl http://localhost:8080/health

# Register a user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123","email":"test@example.com"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'

# Use the token for authenticated requests
curl -X GET http://localhost:8080/api/trading/bot/status \
  -H "Authorization: Bearer <your_jwt_token>"
```

### Using Swagger UI

1. Open the Swagger UI
2. Click on any endpoint
3. Click "Try it out"
4. Fill in the required parameters
5. Click "Execute" to test the endpoint

## Development

### Adding New Endpoints

When adding new endpoints to the API:

1. Update the `swagger.yaml` file with the new endpoint definition
2. Include proper request/response schemas
3. Add appropriate error responses
4. Update this README if necessary

### Schema Updates

When modifying data models:

1. Update the corresponding schema in `swagger.yaml`
2. Ensure all examples are current and accurate
3. Update any related endpoint documentation

## Notes

- The API runs on port 8080 by default (configurable via `PORT` environment variable)
- All timestamps are in ISO 8601 format
- All monetary values are in sats (satoshis) unless otherwise specified
- The API supports CORS for cross-origin requests

## Support

For questions or issues with the API documentation, please refer to the main project documentation or create an issue in the project repository.
