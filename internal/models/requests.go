package models

// LNMarketsConfigRequest representa a request para configurar LNMarkets
// sem os campos que são gerados automaticamente pelo servidor
type LNMarketsConfigRequest struct {
	APIKey     string `json:"api_key"`
	SecretKey  string `json:"secret_key"`
	Passphrase string `json:"passphrase"`
	IsTestnet  bool   `json:"is_testnet"`
}

// MarginProtectionRequest representa a request para configurar proteção de margem
// sem os campos que são gerados automaticamente pelo servidor
type MarginProtectionRequest struct {
	IsEnabled              bool    `json:"is_enabled"`
	ActivationDistance     float64 `json:"activation_distance"`
	NewLiquidationDistance float64 `json:"new_liquidation_distance"`
}

// TakeProfitRequest representa a request para configurar take profit
// sem os campos que são gerados automaticamente pelo servidor
type TakeProfitRequest struct {
	IsEnabled       bool    `json:"is_enabled"`
	DailyPercentage float64 `json:"daily_percentage"`
}

// EntryAutomationRequest representa a request para configurar automação de entrada
// sem os campos que são gerados automaticamente pelo servidor
type EntryAutomationRequest struct {
	IsEnabled          bool    `json:"is_enabled"`
	AmountPerOrder     float64 `json:"amount_per_order"`
	MarginPerOrder     int64   `json:"margin_per_order"`
	NumberOfOrders     int     `json:"number_of_orders"`
	PriceVariation     float64 `json:"price_variation"`
	InitialPrice       float64 `json:"initial_price"`
	TakeProfitPerOrder float64 `json:"take_profit_per_order"`
	OperationType      string  `json:"operation_type"`
	Leverage           float64 `json:"leverage"`
}

// PriceAlertRequest representa a request para configurar alerta de preço
// sem os campos que são gerados automaticamente pelo servidor
type PriceAlertRequest struct {
	IsEnabled     bool    `json:"is_enabled"`
	MinPrice      float64 `json:"min_price"`
	MaxPrice      float64 `json:"max_price"`
	CheckInterval int     `json:"check_interval"`
}
