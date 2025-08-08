package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("BTC Trading Bot")
	fmt.Println("Para executar a aplicação, use:")
	fmt.Println("go run cmd/api/main.go")
	fmt.Println("")
	fmt.Println("Ou compile e execute:")
	fmt.Println("go build -o btc-trading-bot cmd/api/main.go")
	fmt.Println("./btc-trading-bot")

	os.Exit(1)
}
