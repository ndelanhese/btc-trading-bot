# LN Markets Authentication

Este documento explica como a autenticação com o LN Markets está implementada e como usar corretamente.

## Headers de Segurança

O LN Markets requer os seguintes headers de segurança em todas as requisições autenticadas:

- `LNM-ACCESS-KEY`: API key como string
- `LNM-ACCESS-SIGNATURE`: Assinatura da mensagem
- `LNM-ACCESS-PASSPHRASE`: Passphrase da API key
- `LNM-ACCESS-TIMESTAMP`: Timestamp da requisição (em milissegundos desde Unix Epoch)

## Geração da Assinatura

A assinatura é gerada usando HMAC-SHA256 com a seguinte fórmula:

```
message = timestamp + method + path + params
signature = base64(hmac_sha256(secret_key, message))
```

Onde:
- `timestamp`: Valor do header LNM-ACCESS-TIMESTAMP
- `method`: Método HTTP em maiúsculas (GET, POST, etc.)
- `path`: Caminho da URL incluindo query parameters (ex: `/v2/user?param=value`)
- `params`: Body da requisição como JSON string (sem espaços, sem quebras de linha) ou query parameters URL encoded

## Implementação

A implementação está no arquivo `pkg/lnmarkets/client.go`:

### Método `createSignature`
```go
func (c *Client) createSignature(timestamp, method, path, params string) string {
    message := timestamp + method + path + params
    h := hmac.New(sha256.New, []byte(c.SecretKey))
    h.Write([]byte(message))
    return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
```

### Método `makeRequest`
O método `makeRequest` constrói a requisição HTTP com todos os headers necessários:

1. Gera o timestamp em milissegundos
2. Constrói o path completo (incluindo query parameters para GET/DELETE)
3. Serializa o body para JSON (para POST/PUT)
4. Gera a assinatura
5. Adiciona todos os headers de segurança

## Debug

Para ativar logs de debug e ver como a assinatura está sendo gerada:

```bash
export DEBUG_LNMARKETS=true
```

Isso irá mostrar:
- A mensagem usada para gerar a assinatura
- A assinatura gerada
- Os headers enviados na requisição

## Testando

Use os scripts de teste para verificar se a autenticação está funcionando:

```bash
# Teste normal
./scripts/test-lnmarkets-auth.sh

# Teste com debug ativado
./scripts/test-lnmarkets-auth-debug.sh
```

## Exemplos de Uso

### Requisição GET
```
GET /v2/user
Headers:
  LNM-ACCESS-KEY: your_api_key
  LNM-ACCESS-SIGNATURE: base64_hmac_sha256
  LNM-ACCESS-PASSPHRASE: your_passphrase
  LNM-ACCESS-TIMESTAMP: 1640995200000
```

### Requisição POST
```
POST /v2/futures/trade
Headers:
  LNM-ACCESS-KEY: your_api_key
  LNM-ACCESS-SIGNATURE: base64_hmac_sha256
  LNM-ACCESS-PASSPHRASE: your_passphrase
  LNM-ACCESS-TIMESTAMP: 1640995200000
  Content-Type: application/json
Body:
  {"type":"buy","amount":100,"price":50000,"leverage":10}
```

## Troubleshooting

### Erro de Autenticação
1. Verifique se as credenciais estão corretas
2. Ative o debug para ver os detalhes da assinatura
3. Verifique se o timestamp está sincronizado (deve estar dentro de 30 segundos)
4. Verifique se o path está correto (incluindo query parameters)

### Erro de Timestamp
- O timestamp deve estar em milissegundos (não segundos)
- Deve estar dentro de 30 segundos do tempo do servidor
- Use `time.Now().UnixMilli()` para gerar o timestamp

### Erro de Assinatura
- Verifique se o método está em maiúsculas
- Verifique se o path inclui query parameters para requisições GET
- Verifique se o body JSON não tem espaços extras ou quebras de linha


