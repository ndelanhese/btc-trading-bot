# BTC Trading Bot

Bot automatizado de trading para Bitcoin na plataforma LN Markets, desenvolvido em Go.

## Funcionalidades

### 🔐 Autenticação
- Registro e login de usuários
- Autenticação JWT
- Proteção de rotas

### 🤖 Automações de Trading

#### 1. Proteção de Margem
- Monitoramento automático de liquidação
- Ativação quando preço chega a 5% da liquidação
- Aumento automático da distância de liquidação para 10%

#### 2. Take Profit Automático
- Atualização diária de take profit (1%)
- Cálculo automático considerando taxas
- Aplicação imediata em todas as ordens

#### 3. Automação de Entradas (DCA)
- Sistema de Dollar Cost Averaging
- 9 ordens com $10 cada
- Intervalos de $50 entre ordens
- Preço inicial: $116,000
- Take profit de 0.25% por ordem
- Alavancagem 10x

#### 4. Alerta de Preço
- Monitoramento de preços em tempo real
- Alertas quando preço sai do intervalo $100,000 - $120,000
- Verificação a cada 1 minuto

## Tecnologias

- **Backend**: Go 1.24+
- **Database**: PostgreSQL
- **ORM**: SQLx
- **API**: REST + WebSocket
- **Autenticação**: JWT
- **Trading**: LN Markets API

## Estrutura do Projeto

```
btc-trading-bot/
├── cmd/api/           # Ponto de entrada da aplicação
├── internal/          # Código interno da aplicação
│   ├── database/     # Configuração e migrações do banco
│   ├── handlers/     # Handlers HTTP
│   ├── models/       # Modelos de dados
│   └── services/     # Lógica de negócio
├── pkg/              # Pacotes reutilizáveis
│   ├── lnmarkets/    # Cliente da API LN Markets
│   └── websocket/    # Cliente WebSocket
├── go.mod
├── go.sum
└── README.md
```

## Configuração

### 1. Banco de Dados

Crie um banco PostgreSQL e configure as variáveis de ambiente:

```bash
# .env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=btc_trading_bot
DB_SSLMODE=disable

JWT_SECRET=your-secret-key-change-this
PORT=8080
```

### 2. Instalação

```bash
# Clonar o repositório
git clone <repository-url>
cd btc-trading-bot

# Instalar dependências
go mod tidy

# Executar migrações (automático na primeira execução)
go run cmd/api/main.go
```

### 3. Executar

```bash
go run cmd/api/main.go
```

## API Endpoints

### Autenticação
- `POST /api/auth/register` - Registrar usuário
- `POST /api/auth/login` - Login

### Configurações (Protegidas)
- `POST /api/lnmarkets/config` - Configurar API LN Markets
- `GET /api/lnmarkets/config` - Obter configuração LN Markets

### Trading (Protegidas)
- `POST /api/trading/margin-protection` - Configurar proteção de margem
- `GET /api/trading/margin-protection` - Obter configuração
- `POST /api/trading/take-profit` - Configurar take profit
- `GET /api/trading/take-profit` - Obter configuração
- `POST /api/trading/entry-automation` - Configurar automação de entrada
- `GET /api/trading/entry-automation` - Obter configuração
- `POST /api/trading/price-alert` - Configurar alerta de preço
- `GET /api/trading/price-alert` - Obter configuração
- `GET /api/trading/orders` - Listar ordens
- `POST /api/trading/bot/start` - Iniciar bot
- `POST /api/trading/bot/stop` - Parar bot

### Health Check
- `GET /health` - Status da API

## Exemplos de Uso

### 1. Registrar Usuário

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

### 3. Configurar LN Markets

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

### 4. Configurar Proteção de Margem

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

### 5. Configurar Automação de Entrada

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

### 6. Iniciar Bot

```bash
curl -X POST http://localhost:8080/api/trading/bot/start \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Segurança

- Senhas são hasheadas com bcrypt
- Tokens JWT para autenticação
- Rate limiting implementado
- Validação de entrada em todos os endpoints
- Conexões HTTPS recomendadas em produção

## Monitoramento

O bot registra logs detalhados de:
- Atualizações de preço
- Ativação de proteções
- Criação de ordens
- Alertas de preço
- Erros de API

## Desenvolvimento

### Estrutura de Dados

O banco de dados inclui tabelas para:
- Usuários
- Configurações LN Markets
- Proteção de margem
- Take profit
- Automação de entrada
- Alertas de preço
- Ordens de trading

### Fluxo de Dados

1. **WebSocket** recebe atualizações de preço em tempo real
2. **Trading Service** processa cada atualização
3. **Automações** são executadas baseadas nas configurações
4. **LN Markets API** executa ordens quando necessário
5. **Database** armazena histórico e configurações

## Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature
3. Commit suas mudanças
4. Push para a branch
5. Abra um Pull Request

## Licença

Este projeto está sob a licença MIT.
