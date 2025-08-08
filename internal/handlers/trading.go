package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"btc-trading-bot/internal/database"
	"btc-trading-bot/internal/models"
	"btc-trading-bot/internal/services"
)

type TradingHandler struct {
	db             *database.Database
	tradingService *services.TradingService
}

func NewTradingHandler(db *database.Database, tradingService *services.TradingService) *TradingHandler {
	return &TradingHandler{
		db:             db,
		tradingService: tradingService,
	}
}

func (h *TradingHandler) SetLNMarketsConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id").(int)

	var config models.LNMarketsConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	config.UserID = userID
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	var existingConfig models.LNMarketsConfig
	err := h.db.Get(&existingConfig, "SELECT id FROM ln_markets_config WHERE user_id = $1", userID)
	if err == nil {
		_, err = h.db.Exec(`
			UPDATE ln_markets_config 
			SET api_key = $1, secret_key = $2, passphrase = $3, is_testnet = $4, updated_at = $5
			WHERE user_id = $6
		`, config.APIKey, config.SecretKey, config.Passphrase, config.IsTestnet, config.UpdatedAt, userID)
	} else {
		_, err = h.db.Exec(`
			INSERT INTO ln_markets_config (user_id, api_key, secret_key, passphrase, is_testnet, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, config.UserID, config.APIKey, config.SecretKey, config.Passphrase, config.IsTestnet, config.CreatedAt, config.UpdatedAt)
	}

	if err != nil {
		http.Error(w, "Failed to save configuration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Configuration saved successfully"})
}

func (h *TradingHandler) GetLNMarketsConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id").(int)

	var config models.LNMarketsConfig
	err := h.db.Get(&config, "SELECT * FROM ln_markets_config WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, "Configuration not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

func (h *TradingHandler) SetMarginProtection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id").(int)

	var config models.MarginProtection
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	config.UserID = userID
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	var existingConfig models.MarginProtection
	err := h.db.Get(&existingConfig, "SELECT id FROM margin_protection WHERE user_id = $1", userID)
	if err == nil {
		_, err = h.db.Exec(`
			UPDATE margin_protection 
			SET is_enabled = $1, activation_distance = $2, new_liquidation_distance = $3, updated_at = $4
			WHERE user_id = $5
		`, config.IsEnabled, config.ActivationDistance, config.NewLiquidationDistance, config.UpdatedAt, userID)
	} else {
		_, err = h.db.Exec(`
			INSERT INTO margin_protection (user_id, is_enabled, activation_distance, new_liquidation_distance, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, config.UserID, config.IsEnabled, config.ActivationDistance, config.NewLiquidationDistance, config.CreatedAt, config.UpdatedAt)
	}

	if err != nil {
		http.Error(w, "Failed to save configuration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Margin protection configuration saved"})
}

func (h *TradingHandler) GetMarginProtection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id").(int)

	var config models.MarginProtection
	err := h.db.Get(&config, "SELECT * FROM margin_protection WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, "Configuration not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

func (h *TradingHandler) SetTakeProfit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id").(int)

	var config models.TakeProfit
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	config.UserID = userID
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()
	config.LastUpdate = time.Now()

	var existingConfig models.TakeProfit
	err := h.db.Get(&existingConfig, "SELECT id FROM take_profit WHERE user_id = $1", userID)
	if err == nil {
		_, err = h.db.Exec(`
			UPDATE take_profit 
			SET is_enabled = $1, daily_percentage = $2, updated_at = $3
			WHERE user_id = $4
		`, config.IsEnabled, config.DailyPercentage, config.UpdatedAt, userID)
	} else {
		_, err = h.db.Exec(`
			INSERT INTO take_profit (user_id, is_enabled, daily_percentage, last_update, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, config.UserID, config.IsEnabled, config.DailyPercentage, config.LastUpdate, config.CreatedAt, config.UpdatedAt)
	}

	if err != nil {
		http.Error(w, "Failed to save configuration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Take profit configuration saved"})
}

func (h *TradingHandler) GetTakeProfit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id").(int)

	var config models.TakeProfit
	err := h.db.Get(&config, "SELECT * FROM take_profit WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, "Configuration not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

func (h *TradingHandler) SetEntryAutomation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id").(int)

	var config models.EntryAutomation
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	config.UserID = userID
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	var existingConfig models.EntryAutomation
	err := h.db.Get(&existingConfig, "SELECT id FROM entry_automation WHERE user_id = $1", userID)
	if err == nil {
		_, err = h.db.Exec(`
			UPDATE entry_automation 
			SET is_enabled = $1, amount_per_order = $2, margin_per_order = $3, number_of_orders = $4,
				price_variation = $5, initial_price = $6, take_profit_per_order = $7, operation_type = $8, leverage = $9, updated_at = $10
			WHERE user_id = $11
		`, config.IsEnabled, config.AmountPerOrder, config.MarginPerOrder, config.NumberOfOrders,
			config.PriceVariation, config.InitialPrice, config.TakeProfitPerOrder, config.OperationType, config.Leverage, config.UpdatedAt, userID)
	} else {
		_, err = h.db.Exec(`
			INSERT INTO entry_automation (user_id, is_enabled, amount_per_order, margin_per_order, number_of_orders,
				price_variation, initial_price, take_profit_per_order, operation_type, leverage, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`, config.UserID, config.IsEnabled, config.AmountPerOrder, config.MarginPerOrder, config.NumberOfOrders,
			config.PriceVariation, config.InitialPrice, config.TakeProfitPerOrder, config.OperationType, config.Leverage, config.CreatedAt, config.UpdatedAt)
	}

	if err != nil {
		http.Error(w, "Failed to save configuration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Entry automation configuration saved"})
}

func (h *TradingHandler) GetEntryAutomation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id").(int)

	var config models.EntryAutomation
	err := h.db.Get(&config, "SELECT * FROM entry_automation WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, "Configuration not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

func (h *TradingHandler) SetPriceAlert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id").(int)

	var config models.PriceAlert
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	config.UserID = userID
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()
	config.LastAlert = time.Now()

	var existingConfig models.PriceAlert
	err := h.db.Get(&existingConfig, "SELECT id FROM price_alert WHERE user_id = $1", userID)
	if err == nil {
		_, err = h.db.Exec(`
			UPDATE price_alert 
			SET is_enabled = $1, min_price = $2, max_price = $3, check_interval = $4, updated_at = $5
			WHERE user_id = $6
		`, config.IsEnabled, config.MinPrice, config.MaxPrice, config.CheckInterval, config.UpdatedAt, userID)
	} else {
		_, err = h.db.Exec(`
			INSERT INTO price_alert (user_id, is_enabled, min_price, max_price, check_interval, last_alert, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, config.UserID, config.IsEnabled, config.MinPrice, config.MaxPrice, config.CheckInterval, config.LastAlert, config.CreatedAt, config.UpdatedAt)
	}

	if err != nil {
		http.Error(w, "Failed to save configuration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Price alert configuration saved"})
}

func (h *TradingHandler) GetPriceAlert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id").(int)

	var config models.PriceAlert
	err := h.db.Get(&config, "SELECT * FROM price_alert WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, "Configuration not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

func (h *TradingHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id").(int)

	var orders []models.TradingOrder
	err := h.db.Select(&orders, "SELECT * FROM trading_orders WHERE user_id = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *TradingHandler) StartBot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id").(int)

	err := h.tradingService.InitializeClient(userID)
	if err != nil {
		http.Error(w, "Failed to start bot: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Trading bot started successfully"})
}

func (h *TradingHandler) StopBot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.tradingService.Stop()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Trading bot stopped successfully"})
}
