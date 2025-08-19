package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"

	"btc-trading-bot/internal/database"
	"btc-trading-bot/internal/handlers"
	"btc-trading-bot/internal/services"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// CORS middleware function
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get allowed origin from environment variable, default to localhost:3000 for development
		allowedOrigin := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")

		// Check if the request origin matches the allowed origin
		origin := r.Header.Get("Origin")
		if origin == allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		allowedMethods := getEnv("CORS_ALLOWED_METHODS", "GET, POST, PUT, DELETE, OPTIONS")
		allowedHeaders := getEnv("CORS_ALLOWED_HEADERS", "Content-Type, Authorization, X-Requested-With")
		maxAge := getEnv("CORS_MAX_AGE", "86400")

		w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
		w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
		w.Header().Set("Access-Control-Max-Age", maxAge) // 24 hours

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	jwtSecret := getEnv("JWT_SECRET", randomString(32))
	authService := services.NewAuthService(db, jwtSecret)
	tradingService := services.NewTradingService(db)

	authHandler := handlers.NewAuthHandler(authService)
	tradingHandler := handlers.NewTradingHandler(db, tradingService)

	router := mux.NewRouter()

	// Apply CORS middleware to all routes
	router.Use(corsMiddleware)

	// Auth endpoints
	router.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")

	// Handle OPTIONS for all API routes
	router.HandleFunc("/api/{rest:.*}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.NotFound(w, r)
	}).Methods("OPTIONS")

	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHandler.AuthMiddleware(next.ServeHTTP)(w, r)
		})
	})

	protected.HandleFunc("/lnmarkets/config", tradingHandler.SetLNMarketsConfig).Methods("POST")
	protected.HandleFunc("/lnmarkets/config", tradingHandler.GetLNMarketsConfig).Methods("GET")

	protected.HandleFunc("/trading/margin-protection", tradingHandler.SetMarginProtection).Methods("POST")
	protected.HandleFunc("/trading/margin-protection", tradingHandler.GetMarginProtection).Methods("GET")

	protected.HandleFunc("/trading/take-profit", tradingHandler.SetTakeProfit).Methods("POST")
	protected.HandleFunc("/trading/take-profit", tradingHandler.GetTakeProfit).Methods("GET")

	protected.HandleFunc("/trading/entry-automation", tradingHandler.SetEntryAutomation).Methods("POST")
	protected.HandleFunc("/trading/entry-automation", tradingHandler.GetEntryAutomation).Methods("GET")

	protected.HandleFunc("/trading/price-alert", tradingHandler.SetPriceAlert).Methods("POST")
	protected.HandleFunc("/trading/price-alert", tradingHandler.GetPriceAlert).Methods("GET")

	protected.HandleFunc("/trading/orders", tradingHandler.GetOrders).Methods("GET")

	protected.HandleFunc("/trading/bot/start", tradingHandler.StartBot).Methods("POST")
	protected.HandleFunc("/trading/bot/stop", tradingHandler.StopBot).Methods("POST")
	protected.HandleFunc("/trading/bot/status", tradingHandler.GetBotStatus).Methods("GET")
	protected.HandleFunc("/trading/account/balance", tradingHandler.GetAccountBalance).Methods("GET")
	protected.HandleFunc("/trading/positions", tradingHandler.GetPositions).Methods("GET")
	protected.HandleFunc("/trading/positions/{id}", tradingHandler.GetPosition).Methods("GET")
	protected.HandleFunc("/trading/positions/{id}/close", tradingHandler.ClosePosition).Methods("POST")
	protected.HandleFunc("/trading/positions/{id}/take-profit", tradingHandler.UpdateTakeProfit).Methods("POST")
	protected.HandleFunc("/trading/positions/{id}/stop-loss", tradingHandler.UpdateStopLoss).Methods("POST")

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "message": "BTC Trading Bot API is running"}`))
	}).Methods("GET")

	port := getEnv("PORT", "8080")
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func randomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
