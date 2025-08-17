# API Requests - Estruturas Simplificadas

Este documento descreve as estruturas de request simplificadas para a API do trading bot. Os campos que são gerados automaticamente pelo servidor (ID, UserID, CreatedAt, UpdatedAt, etc.) não precisam mais ser enviados pelo cliente.

## Configuração LNMarkets

**Endpoint:** `POST /api/lnmarkets/config`

**Request Body:**
```json
{
  "api_key": "sua_api_key_aqui",
  "secret_key": "sua_secret_key_aqui", 
  "passphrase": "sua_passphrase_aqui",
  "is_testnet": true
}
```

## Proteção de Margem

**Endpoint:** `POST /api/margin-protection`

**Request Body:**
```json
{
  "is_enabled": true,
  "activation_distance": 0.5,
  "new_liquidation_distance": 0.3
}
```

## Take Profit

**Endpoint:** `POST /api/take-profit`

**Request Body:**
```json
{
  "is_enabled": true,
  "daily_percentage": 1.0
}
```

## Automação de Entrada

**Endpoint:** `POST /api/entry-automation`

**Request Body:**
```json
{
  "is_enabled": true,
  "amount_per_order": 100.0,
  "margin_per_order": 1000,
  "number_of_orders": 5,
  "price_variation": 0.1,
  "initial_price": 50000.0,
  "take_profit_per_order": 0.5,
  "operation_type": "long",
  "leverage": 10
}
```

## Alerta de Preço

**Endpoint:** `POST /api/price-alert`

**Request Body:**
```json
{
  "is_enabled": true,
  "min_price": 45000.0,
  "max_price": 55000.0,
  "check_interval": 300
}
```

## Atualizar Take Profit de Posição

**Endpoint:** `POST /api/positions/{id}/take-profit`

**Request Body:**
```json
{
  "price": 52000.0
}
```

## Atualizar Stop Loss de Posição

**Endpoint:** `POST /api/positions/{id}/stop-loss`

**Request Body:**
```json
{
  "price": 48000.0
}
```

## Mudanças Implementadas

### Antes:
- O cliente precisava enviar campos como `id`, `user_id`, `created_at`, `updated_at`
- Esses campos eram ignorados ou sobrescritos pelo servidor
- A API era mais confusa e propensa a erros

### Depois:
- O cliente só envia os campos realmente necessários
- Campos como `user_id`, `created_at`, `updated_at` são gerados automaticamente pelo servidor
- A API é mais limpa e intuitiva
- Menor chance de erros por campos desnecessários

## Benefícios

1. **Simplicidade**: Requests mais simples e diretas
2. **Segurança**: O servidor controla campos sensíveis como `user_id`
3. **Consistência**: Datas são sempre geradas no momento correto
4. **Manutenibilidade**: Código mais limpo e fácil de manter
5. **Documentação**: API mais fácil de documentar e usar
