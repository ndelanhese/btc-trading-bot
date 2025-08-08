package services

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"btc-trading-bot/internal/database"
	"btc-trading-bot/internal/models"
	"btc-trading-bot/pkg/lnmarkets"
	"btc-trading-bot/pkg/websocket"
)

type TradingService struct {
	db           *database.Database
	lnClient     *lnmarkets.Client
	wsClient     *websocket.Client
	priceUpdates chan float64
	stopChan     chan struct{}
}

type TradingConfig struct {
	UserID           int
	MarginProtection *models.MarginProtection
	TakeProfit       *models.TakeProfit
	EntryAutomation  *models.EntryAutomation
	PriceAlert       *models.PriceAlert
	LNMarketsConfig  *models.LNMarketsConfig
}

func NewTradingService(db *database.Database) *TradingService {
	return &TradingService{
		db:           db,
		priceUpdates: make(chan float64, 100),
		stopChan:     make(chan struct{}),
	}
}

func (s *TradingService) InitializeClient(userID int) error {
	var config models.LNMarketsConfig
	err := s.db.Get(&config, "SELECT * FROM ln_markets_config WHERE user_id = $1", userID)
	if err != nil {
		return fmt.Errorf("LN Markets config not found: %v", err)
	}

	s.lnClient = lnmarkets.NewClient(config.APIKey, config.SecretKey, config.Passphrase, config.IsTestnet)

	wsURL := "wss://api.lnmarkets.com"
	if config.IsTestnet {
		wsURL = "wss://api.testnet4.lnmarkets.com"
	}

	s.wsClient = websocket.NewClient(wsURL)

	if err := s.wsClient.Connect(); err != nil {
		return fmt.Errorf("failed to connect to websocket: %v", err)
	}

	s.wsClient.OnPriceUpdate(func(data []byte) {
		var priceData lnmarkets.PriceData
		if err := json.Unmarshal(data, &priceData); err == nil {
			s.priceUpdates <- priceData.Price
		}
	})

	if err := s.wsClient.Subscribe("futures:btc_usd:last-price"); err != nil {
		return fmt.Errorf("failed to subscribe to price updates: %v", err)
	}

	go s.processPriceUpdates(userID)

	return nil
}

func (s *TradingService) processPriceUpdates(userID int) {
	for {
		select {
		case price := <-s.priceUpdates:
			s.handlePriceUpdate(userID, price)
		case <-s.stopChan:
			return
		}
	}
}

func (s *TradingService) handlePriceUpdate(userID int, price float64) {
	log.Printf("Price update: $%.2f", price)

	config, err := s.getTradingConfig(userID)
	if err != nil {
		log.Printf("Error getting trading config: %v", err)
		return
	}

	go s.checkMarginProtection(config, price)
	go s.checkTakeProfit(config, price)
	go s.checkEntryAutomation(config, price)
	go s.checkPriceAlert(config, price)
}

func (s *TradingService) getTradingConfig(userID int) (*TradingConfig, error) {
	config := &TradingConfig{UserID: userID}

	var marginProtection models.MarginProtection
	err := s.db.Get(&marginProtection, "SELECT * FROM margin_protection WHERE user_id = $1", userID)
	if err == nil {
		config.MarginProtection = &marginProtection
	}

	var takeProfit models.TakeProfit
	err = s.db.Get(&takeProfit, "SELECT * FROM take_profit WHERE user_id = $1", userID)
	if err == nil {
		config.TakeProfit = &takeProfit
	}

	var entryAutomation models.EntryAutomation
	err = s.db.Get(&entryAutomation, "SELECT * FROM entry_automation WHERE user_id = $1", userID)
	if err == nil {
		config.EntryAutomation = &entryAutomation
	}

	var priceAlert models.PriceAlert
	err = s.db.Get(&priceAlert, "SELECT * FROM price_alert WHERE user_id = $1", userID)
	if err == nil {
		config.PriceAlert = &priceAlert
	}

	var lnConfig models.LNMarketsConfig
	err = s.db.Get(&lnConfig, "SELECT * FROM ln_markets_config WHERE user_id = $1", userID)
	if err == nil {
		config.LNMarketsConfig = &lnConfig
	}

	return config, nil
}

func (s *TradingService) checkMarginProtection(config *TradingConfig, currentPrice float64) {
	if config.MarginProtection == nil || !config.MarginProtection.IsEnabled {
		return
	}

	var orders []models.TradingOrder
	err := s.db.Select(&orders, "SELECT * FROM trading_orders WHERE user_id = $1 AND status = 'open'", config.UserID)
	if err != nil {
		log.Printf("Error getting open orders: %v", err)
		return
	}

	for _, order := range orders {
		liquidationPrice := order.Price * (1 - float64(order.Leverage)/100)

		distanceToLiquidation := math.Abs(currentPrice-liquidationPrice) / liquidationPrice * 100

		if distanceToLiquidation <= config.MarginProtection.ActivationDistance {
			log.Printf("Margin protection activated for order %s", order.OrderID)

			newTakeProfitPrice := liquidationPrice * (1 + config.MarginProtection.NewLiquidationDistance/100)

			if s.lnClient != nil {
				log.Printf("Updating take profit to $%.2f", newTakeProfitPrice)
			}
			_, err := s.db.Exec("UPDATE trading_orders SET take_profit_price = $1, updated_at = $2 WHERE id = $3",
				newTakeProfitPrice, time.Now(), order.ID)
			if err != nil {
				log.Printf("Error updating order: %v", err)
			}
		}
	}
}

func (s *TradingService) checkTakeProfit(config *TradingConfig, currentPrice float64) {
	if config.TakeProfit == nil || !config.TakeProfit.IsEnabled {
		return
	}

	if time.Since(config.TakeProfit.LastUpdate) < 24*time.Hour {
		return
	}

	var orders []models.TradingOrder
	err := s.db.Select(&orders, "SELECT * FROM trading_orders WHERE user_id = $1 AND status = 'open'", config.UserID)
	if err != nil {
		log.Printf("Error getting open orders: %v", err)
		return
	}

	for _, order := range orders {
		newTakeProfitPrice := order.Price * (1 + config.TakeProfit.DailyPercentage/100)

		if s.lnClient != nil {
			log.Printf("Updating take profit for order %s to $%.2f", order.OrderID, newTakeProfitPrice)
		}

		_, err := s.db.Exec("UPDATE trading_orders SET take_profit_price = $1, updated_at = $2 WHERE id = $3",
			newTakeProfitPrice, time.Now(), order.ID)
		if err != nil {
			log.Printf("Error updating order: %v", err)
		}
	}

	_, err = s.db.Exec("UPDATE take_profit SET last_update = $1 WHERE user_id = $2",
		time.Now(), config.UserID)
	if err != nil {
		log.Printf("Error updating take profit last_update: %v", err)
	}
}

func (s *TradingService) checkEntryAutomation(config *TradingConfig, currentPrice float64) {
	if config.EntryAutomation == nil || !config.EntryAutomation.IsEnabled {
		return
	}

	if config.EntryAutomation.FilledSlots >= config.EntryAutomation.NumberOfOrders {
		return
	}

	targetPrice := config.EntryAutomation.InitialPrice +
		float64(config.EntryAutomation.FilledSlots)*config.EntryAutomation.PriceVariation

	if math.Abs(currentPrice-targetPrice) <= config.EntryAutomation.PriceVariation/2 {
		trade := &lnmarkets.TradeRequest{
			Type:     config.EntryAutomation.OperationType,
			Amount:   config.EntryAutomation.AmountPerOrder,
			Price:    currentPrice,
			Leverage: config.EntryAutomation.Leverage,
		}

		if s.lnClient != nil {
			tradeResp, err := s.lnClient.CreateTrade(trade)
			if err != nil {
				log.Printf("Error creating trade: %v", err)
				return
			}

			order := &models.TradingOrder{
				UserID:          config.UserID,
				OrderID:         tradeResp.ID,
				Type:            tradeResp.Type,
				Amount:          tradeResp.Amount,
				Price:           tradeResp.Price,
				Leverage:        tradeResp.Leverage,
				Status:          tradeResp.Status,
				TakeProfitPrice: currentPrice * (1 + config.EntryAutomation.TakeProfitPerOrder/100),
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}

			_, err = s.db.NamedExec(`
				INSERT INTO trading_orders (user_id, order_id, type, amount, price, leverage, status, take_profit_price, created_at, updated_at)
				VALUES (:user_id, :order_id, :type, :amount, :price, :leverage, :status, :take_profit_price, :created_at, :updated_at)
			`, order)
			if err != nil {
				log.Printf("Error saving order: %v", err)
				return
			}

			_, err = s.db.Exec("UPDATE entry_automation SET filled_slots = filled_slots + 1 WHERE user_id = $1",
				config.UserID)
			if err != nil {
				log.Printf("Error updating filled slots: %v", err)
			}

			log.Printf("Created new order: %s at price $%.2f", tradeResp.ID, currentPrice)
		}
	}
}

func (s *TradingService) checkPriceAlert(config *TradingConfig, currentPrice float64) {
	if config.PriceAlert == nil || !config.PriceAlert.IsEnabled {
		return
	}

	if currentPrice < config.PriceAlert.MinPrice || currentPrice > config.PriceAlert.MaxPrice {
		if time.Since(config.PriceAlert.LastAlert) >= time.Duration(config.PriceAlert.CheckInterval)*time.Second {
			log.Printf("PRICE ALERT: Bitcoin price $%.2f is outside range $%.2f - $%.2f",
				currentPrice, config.PriceAlert.MinPrice, config.PriceAlert.MaxPrice)

			_, err := s.db.Exec("UPDATE price_alert SET last_alert = $1 WHERE user_id = $2",
				time.Now(), config.UserID)
			if err != nil {
				log.Printf("Error updating last alert: %v", err)
			}
		}
	}
}

func (s *TradingService) Stop() {
	close(s.stopChan)
	if s.wsClient != nil {
		s.wsClient.Disconnect()
	}
}
