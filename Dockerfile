FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Compilar a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o btc-trading-bot main.go

# Compilar o comando de migração
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate cmd/migrate/main.go

# Imagem final
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar binários compilados
COPY --from=builder /app/btc-trading-bot .
COPY --from=builder /app/migrate .

# Copiar scripts directory
COPY --from=builder /app/scripts ./scripts

# Expor porta
EXPOSE 8080

# Comando para executar
CMD ["./btc-trading-bot"]
