package lnmarkets

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	APIKey     string
	SecretKey  string
	Passphrase string
	BaseURL    string
	HTTPClient *http.Client
}

type TradeRequest struct {
	Type     string  `json:"type"`
	Amount   float64 `json:"amount"`
	Price    float64 `json:"price"`
	Leverage int     `json:"leverage"`
}

type TradeResponse struct {
	ID       string  `json:"id"`
	Type     string  `json:"type"`
	Amount   float64 `json:"amount"`
	Price    float64 `json:"price"`
	Leverage int     `json:"leverage"`
	Status   string  `json:"status"`
}

type PriceData struct {
	Price float64 `json:"price"`
	Time  int64   `json:"time"`
}

type UserData struct {
	ID       string  `json:"id"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

func NewClient(apiKey, secretKey, passphrase string, isTestnet bool) *Client {
	baseURL := "https://api.lnmarkets.com/v2"
	if isTestnet {
		baseURL = "https://api.testnet4.lnmarkets.com/v2"
	}

	return &Client{
		APIKey:     apiKey,
		SecretKey:  secretKey,
		Passphrase: passphrase,
		BaseURL:    baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) createSignature(timestamp, method, path, params string) string {
	message := timestamp + method + path + params
	h := hmac.New(sha256.New, []byte(c.SecretKey))
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (c *Client) makeRequest(method, path string, data interface{}) ([]byte, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	var params string
	var body io.Reader

	if method == "GET" || method == "DELETE" {
		if data != nil {
			queryParams := url.Values{}
			if mapData, ok := data.(map[string]interface{}); ok {
				for k, v := range mapData {
					queryParams.Set(k, fmt.Sprintf("%v", v))
				}
			}
			params = queryParams.Encode()
		}
	} else {
		if data != nil {
			jsonData, err := json.Marshal(data)
			if err != nil {
				return nil, err
			}
			params = string(jsonData)
			body = bytes.NewBuffer(jsonData)
		}
	}

	signature := c.createSignature(timestamp, method, path, params)

	req, err := http.NewRequest(method, c.BaseURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("LNM-ACCESS-KEY", c.APIKey)
	req.Header.Set("LNM-ACCESS-SIGNATURE", signature)
	req.Header.Set("LNM-ACCESS-PASSPHRASE", c.Passphrase)
	req.Header.Set("LNM-ACCESS-TIMESTAMP", timestamp)

	if method != "GET" && method != "DELETE" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(respBody))
	}

	return respBody, nil
}

func (c *Client) GetUser() (*UserData, error) {
	resp, err := c.makeRequest("GET", "/user", nil)
	if err != nil {
		return nil, err
	}

	var user UserData
	if err := json.Unmarshal(resp, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Client) CreateTrade(trade *TradeRequest) (*TradeResponse, error) {
	resp, err := c.makeRequest("POST", "/futures/trade", trade)
	if err != nil {
		return nil, err
	}

	var tradeResp TradeResponse
	if err := json.Unmarshal(resp, &tradeResp); err != nil {
		return nil, err
	}

	return &tradeResp, nil
}

func (c *Client) GetTrades() ([]TradeResponse, error) {
	resp, err := c.makeRequest("GET", "/futures/trades", nil)
	if err != nil {
		return nil, err
	}

	var trades []TradeResponse
	if err := json.Unmarshal(resp, &trades); err != nil {
		return nil, err
	}

	return trades, nil
}

func (c *Client) GetTrade(tradeID string) (*TradeResponse, error) {
	resp, err := c.makeRequest("GET", "/futures/trade/"+tradeID, nil)
	if err != nil {
		return nil, err
	}

	var trade TradeResponse
	if err := json.Unmarshal(resp, &trade); err != nil {
		return nil, err
	}

	return &trade, nil
}

func (c *Client) CloseTrade(tradeID string) error {
	_, err := c.makeRequest("POST", "/futures/close", map[string]string{"id": tradeID})
	return err
}

func (c *Client) GetPrice() (*PriceData, error) {
	resp, err := c.makeRequest("GET", "/oracle/index", nil)
	if err != nil {
		return nil, err
	}

	var priceData PriceData
	if err := json.Unmarshal(resp, &priceData); err != nil {
		return nil, err
	}

	return &priceData, nil
}

func (c *Client) GetAccountBalance() (*UserData, error) {
	resp, err := c.makeRequest("GET", "/user", nil)
	if err != nil {
		return nil, err
	}

	var userData UserData
	if err := json.Unmarshal(resp, &userData); err != nil {
		return nil, err
	}

	return &userData, nil
}

func (c *Client) GetPositions() ([]TradeResponse, error) {
	resp, err := c.makeRequest("GET", "/futures/positions", nil)
	if err != nil {
		return nil, err
	}

	var positions []TradeResponse
	if err := json.Unmarshal(resp, &positions); err != nil {
		return nil, err
	}

	return positions, nil
}

func (c *Client) GetPosition(positionID string) (*TradeResponse, error) {
	resp, err := c.makeRequest("GET", "/futures/position/"+positionID, nil)
	if err != nil {
		return nil, err
	}

	var position TradeResponse
	if err := json.Unmarshal(resp, &position); err != nil {
		return nil, err
	}

	return &position, nil
}

func (c *Client) ClosePosition(positionID string) error {
	_, err := c.makeRequest("POST", "/futures/close", map[string]string{"id": positionID})
	return err
}

func (c *Client) UpdateTakeProfit(positionID string, takeProfitPrice float64) error {
	_, err := c.makeRequest("POST", "/futures/take-profit", map[string]interface{}{
		"id":    positionID,
		"price": takeProfitPrice,
	})
	return err
}

func (c *Client) UpdateStopLoss(positionID string, stopLossPrice float64) error {
	_, err := c.makeRequest("POST", "/futures/stop-loss", map[string]interface{}{
		"id":    positionID,
		"price": stopLossPrice,
	})
	return err
}
