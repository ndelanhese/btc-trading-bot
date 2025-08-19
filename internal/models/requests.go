package models

// LNMarketsConfigRequest representa a request para configurar LNMarkets
// sem os campos que são gerados automaticamente pelo servidor
type LNMarketsConfigRequest struct {
	APIKey     string      `json:"api_key"`
	SecretKey  string      `json:"secret_key"`
	Passphrase string      `json:"passphrase"`
	IsTestnet  interface{} `json:"is_testnet"` // Aceita bool ou string
}

// GetIsTestnet converte o IsTestnet para boolean
func (r *LNMarketsConfigRequest) GetIsTestnet() bool {
	if r.IsTestnet == nil {
		return false
	}

	switch v := r.IsTestnet.(type) {
	case bool:
		return v
	case string:
		return v == "true" || v == "on" || v == "1" || v == "yes"
	case float64:
		return v != 0
	case int:
		return v != 0
	default:
		return false
	}
}

// MarginProtectionRequest representa a request para configurar proteção de margem
// sem os campos que são gerados automaticamente pelo servidor
type MarginProtectionRequest struct {
	IsEnabled              interface{} `json:"is_enabled"` // Aceita bool ou string
	ActivationDistance     float64     `json:"activation_distance"`
	NewLiquidationDistance float64     `json:"new_liquidation_distance"`
}

// GetIsEnabled converte o IsEnabled para boolean
func (r *MarginProtectionRequest) GetIsEnabled() bool {
	if r.IsEnabled == nil {
		return false
	}

	switch v := r.IsEnabled.(type) {
	case bool:
		return v
	case string:
		return v == "true" || v == "on" || v == "1" || v == "yes"
	case float64:
		return v != 0
	case int:
		return v != 0
	default:
		return false
	}
}

// TakeProfitRequest representa a request para configurar take profit
// sem os campos que são gerados automaticamente pelo servidor
type TakeProfitRequest struct {
	IsEnabled       interface{} `json:"is_enabled"` // Aceita bool ou string
	DailyPercentage float64     `json:"daily_percentage"`
}

// GetIsEnabled converte o IsEnabled para boolean
func (r *TakeProfitRequest) GetIsEnabled() bool {
	if r.IsEnabled == nil {
		return false
	}

	switch v := r.IsEnabled.(type) {
	case bool:
		return v
	case string:
		return v == "true" || v == "on" || v == "1" || v == "yes"
	case float64:
		return v != 0
	case int:
		return v != 0
	default:
		return false
	}
}

// EntryAutomationRequest representa a request para configurar automação de entrada
// sem os campos que são gerados automaticamente pelo servidor
type EntryAutomationRequest struct {
	IsEnabled          interface{} `json:"is_enabled"` // Aceita bool ou string
	AmountPerOrder     float64     `json:"amount_per_order"`
	MarginPerOrder     int64       `json:"margin_per_order"`
	NumberOfOrders     int         `json:"number_of_orders"`
	PriceVariation     float64     `json:"price_variation"`
	InitialPrice       float64     `json:"initial_price"`
	TakeProfitPerOrder float64     `json:"take_profit_per_order"`
	OperationType      string      `json:"operation_type"`
	Leverage           float64     `json:"leverage"`
}

// GetIsEnabled converte o IsEnabled para boolean
func (r *EntryAutomationRequest) GetIsEnabled() bool {
	if r.IsEnabled == nil {
		return false
	}

	switch v := r.IsEnabled.(type) {
	case bool:
		return v
	case string:
		return v == "true" || v == "on" || v == "1" || v == "yes"
	case float64:
		return v != 0
	case int:
		return v != 0
	default:
		return false
	}
}

// PriceAlertRequest representa a request para configurar alerta de preço
// sem os campos que são gerados automaticamente pelo servidor
type PriceAlertRequest struct {
	IsEnabled     interface{} `json:"is_enabled"` // Aceita bool ou string
	MinPrice      float64     `json:"min_price"`
	MaxPrice      float64     `json:"max_price"`
	CheckInterval int         `json:"check_interval"`
}

// GetIsEnabled converte o IsEnabled para boolean
func (r *PriceAlertRequest) GetIsEnabled() bool {
	if r.IsEnabled == nil {
		return false
	}

	switch v := r.IsEnabled.(type) {
	case bool:
		return v
	case string:
		return v == "true" || v == "on" || v == "1" || v == "yes"
	case float64:
		return v != 0
	case int:
		return v != 0
	default:
		return false
	}
}
