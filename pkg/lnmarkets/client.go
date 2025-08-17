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
	"os"
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
	Leverage float64 `json:"leverage"`
}

type TradeResponse struct {
	ID       string  `json:"id"`
	Type     string  `json:"type"`
	Amount   float64 `json:"amount"`
	Price    float64 `json:"price"`
	Leverage float64 `json:"leverage"`
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
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Debug log (only if DEBUG_LNMARKETS is set)
	if os.Getenv("DEBUG_LNMARKETS") == "true" {
		fmt.Printf("DEBUG: Raw Message: %s\n", message)
		fmt.Printf("DEBUG: Creating signature for message: %s\n", message)
		fmt.Printf("DEBUG: Generated signature: %s\n", signature)
	}

	return signature
}

func (c *Client) makeRequest(method, path string, data interface{}) ([]byte, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	var params string
	var body io.Reader
	var cleanPath string
	var fullPath string

	// Parse the path to separate base path from query parameters
	parsedPath, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %v", err)
	}

	cleanPath = parsedPath.Path
	pathQueryParams := parsedPath.Query()

	if method == "GET" || method == "DELETE" {
		// Combine path query params with data query params
		queryParams := url.Values{}

		// Add query params from the path
		for k, v := range pathQueryParams {
			queryParams[k] = v
		}

		// Add query params from data
		if data != nil {
			if mapData, ok := data.(map[string]interface{}); ok {
				for k, v := range mapData {
					queryParams.Set(k, fmt.Sprintf("%v", v))
				}
			}
		}

		params = queryParams.Encode()
		if params != "" {
			fullPath = cleanPath + "?" + params
		} else {
			fullPath = cleanPath
		}
	} else {
		fullPath = cleanPath
		if data != nil {
			jsonData, err := json.Marshal(data)
			if err != nil {
				return nil, err
			}
			params = string(jsonData)
			body = bytes.NewBuffer(jsonData)
		}
	}

	baseURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %v", err)
	}
	signaturePath := baseURL.Path + cleanPath
	signature := c.createSignature(timestamp, method, signaturePath, params)

	req, err := http.NewRequest(method, c.BaseURL+fullPath, body)
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

	// Debug log (only if DEBUG_LNMARKETS is set)
	if os.Getenv("DEBUG_LNMARKETS") == "true" {
		fmt.Printf("DEBUG: Request URL: %s\n", c.BaseURL+fullPath)
		fmt.Printf("DEBUG: Request method: %s\n", method)
		fmt.Printf("DEBUG: Request headers:\n")
		fmt.Printf("  LNM-ACCESS-KEY: %s\n", c.APIKey)
		fmt.Printf("  LNM-ACCESS-SIGNATURE: %s\n", signature)
		fmt.Printf("  LNM-ACCESS-PASSPHRASE: %s\n", c.Passphrase)
		fmt.Printf("  LNM-ACCESS-TIMESTAMP: %s\n", timestamp)
		if method != "GET" && method != "DELETE" {
			fmt.Printf("  Content-Type: application/json\n")
		}
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

func (c *Client) GetPositions(positionType string) ([]TradeResponse, error) {
	var path string
	if positionType != "" {
		path = "/futures?type=" + positionType
	} else {
		path = "/futures?type=running"
	}

	resp, err := c.makeRequest("GET", path, nil)
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
	resp, err := c.makeRequest("GET", "/futures/trades/"+positionID, nil)
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
