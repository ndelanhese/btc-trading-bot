package models

import (
	"time"
)

type LNMarketsConfig struct {
	ID         int       `db:"id" json:"id"`
	UserID     int       `db:"user_id" json:"user_id"`
	APIKey     string    `db:"api_key" json:"api_key"`
	SecretKey  string    `db:"secret_key" json:"secret_key"`
	Passphrase string    `db:"passphrase" json:"passphrase"`
	IsTestnet  bool      `db:"is_testnet" json:"is_testnet"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

type MarginProtection struct {
	ID                     int       `db:"id" json:"id"`
	UserID                 int       `db:"user_id" json:"user_id"`
	IsEnabled              bool      `db:"is_enabled" json:"is_enabled"`
	ActivationDistance     float64   `db:"activation_distance" json:"activation_distance"`
	NewLiquidationDistance float64   `db:"new_liquidation_distance" json:"new_liquidation_distance"`
	CreatedAt              time.Time `db:"created_at" json:"created_at"`
	UpdatedAt              time.Time `db:"updated_at" json:"updated_at"`
}

type TakeProfit struct {
	ID              int       `db:"id" json:"id"`
	UserID          int       `db:"user_id" json:"user_id"`
	IsEnabled       bool      `db:"is_enabled" json:"is_enabled"`
	DailyPercentage float64   `db:"daily_percentage" json:"daily_percentage"` // 1%
	LastUpdate      time.Time `db:"last_update" json:"last_update"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

type EntryAutomation struct {
	ID                 int       `db:"id" json:"id"`
	UserID             int       `db:"user_id" json:"user_id"`
	IsEnabled          bool      `db:"is_enabled" json:"is_enabled"`
	AmountPerOrder     float64   `db:"amount_per_order" json:"amount_per_order"`
	MarginPerOrder     int64     `db:"margin_per_order" json:"margin_per_order"`
	NumberOfOrders     int       `db:"number_of_orders" json:"number_of_orders"`
	FilledSlots        int       `db:"filled_slots" json:"filled_slots"`
	PriceVariation     float64   `db:"price_variation" json:"price_variation"`
	InitialPrice       float64   `db:"initial_price" json:"initial_price"`
	TakeProfitPerOrder float64   `db:"take_profit_per_order" json:"take_profit_per_order"`
	OperationType      string    `db:"operation_type" json:"operation_type"`
	Leverage           int       `db:"leverage" json:"leverage"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
}

type PriceAlert struct {
	ID            int       `db:"id" json:"id"`
	UserID        int       `db:"user_id" json:"user_id"`
	IsEnabled     bool      `db:"is_enabled" json:"is_enabled"`
	MinPrice      float64   `db:"min_price" json:"min_price"`
	MaxPrice      float64   `db:"max_price" json:"max_price"`
	CheckInterval int       `db:"check_interval" json:"check_interval"`
	LastAlert     time.Time `db:"last_alert" json:"last_alert"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type TradingOrder struct {
	ID              int       `db:"id" json:"id"`
	UserID          int       `db:"user_id" json:"user_id"`
	OrderID         string    `db:"order_id" json:"order_id"`
	Type            string    `db:"type" json:"type"`
	Amount          float64   `db:"amount" json:"amount"`
	Price           float64   `db:"price" json:"price"`
	Leverage        int       `db:"leverage" json:"leverage"`
	Status          string    `db:"status" json:"status"`
	TakeProfitPrice float64   `db:"take_profit_price" json:"take_profit_price"`
	StopLossPrice   float64   `db:"stop_loss_price" json:"stop_loss_price"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}
