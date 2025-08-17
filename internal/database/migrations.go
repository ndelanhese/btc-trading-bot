package database

import (
	"fmt"
	"log"
)

func (db *Database) RunMigrations() error {
	log.Println("Running database migrations...")

	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS ln_markets_config (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			api_key VARCHAR(255) NOT NULL,
			secret_key VARCHAR(255) NOT NULL,
			passphrase VARCHAR(255) NOT NULL,
			is_testnet BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS margin_protection (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			is_enabled BOOLEAN DEFAULT false,
			activation_distance DECIMAL(5,2) DEFAULT 5.0,
			new_liquidation_distance DECIMAL(5,2) DEFAULT 10.0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS take_profit (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			is_enabled BOOLEAN DEFAULT false,
			daily_percentage DECIMAL(5,2) DEFAULT 1.0,
			last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS entry_automation (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			is_enabled BOOLEAN DEFAULT false,
			amount_per_order DECIMAL(10,2) DEFAULT 10.0,
			margin_per_order BIGINT DEFAULT 855,
			number_of_orders INTEGER DEFAULT 9,
			filled_slots INTEGER DEFAULT 0,
			price_variation DECIMAL(10,2) DEFAULT 50.0,
			initial_price DECIMAL(15,2) DEFAULT 116000.0,
			take_profit_per_order DECIMAL(5,2) DEFAULT 0.25,
			operation_type VARCHAR(10) DEFAULT 'buy',
			leverage DECIMAL(5,2) DEFAULT 10.0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS price_alert (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			is_enabled BOOLEAN DEFAULT false,
			min_price DECIMAL(15,2) DEFAULT 100000.0,
			max_price DECIMAL(15,2) DEFAULT 120000.0,
			check_interval INTEGER DEFAULT 60,
			last_alert TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS trading_orders (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			order_id VARCHAR(100) UNIQUE NOT NULL,
			type VARCHAR(10) NOT NULL,
			amount DECIMAL(15,2) NOT NULL,
			price DECIMAL(15,2) NOT NULL,
			leverage DECIMAL(5,2) NOT NULL,
			status VARCHAR(20) DEFAULT 'open',
			take_profit_price DECIMAL(15,2),
			stop_loss_price DECIMAL(15,2),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,
		`CREATE INDEX IF NOT EXISTS idx_ln_markets_config_user_id ON ln_markets_config(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_margin_protection_user_id ON margin_protection(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_take_profit_user_id ON take_profit(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_entry_automation_user_id ON entry_automation(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_price_alert_user_id ON price_alert(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_trading_orders_user_id ON trading_orders(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_trading_orders_status ON trading_orders(status)`,
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration %d failed: %v", i+1, err)
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}
