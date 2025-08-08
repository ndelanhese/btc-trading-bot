package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	url      string
	mu       sync.Mutex
	handlers map[string]func([]byte)
	done     chan struct{}
}

type Message struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	ID      string      `json:"id"`
	Params  interface{} `json:"params,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

type PriceUpdate struct {
	Price float64 `json:"price"`
	Time  int64   `json:"time"`
}

func NewClient(url string) *Client {
	return &Client{
		url:      url,
		handlers: make(map[string]func([]byte)),
		done:     make(chan struct{}),
	}
}

func (c *Client) Connect() error {
	var err error
	c.conn, _, err = websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to websocket: %v", err)
	}

	log.Printf("Connected to WebSocket: %s", c.url)

	go c.readMessages()

	go c.heartbeat()

	return nil
}

func (c *Client) Disconnect() error {
	close(c.done)
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) readMessages() {
	for {
		select {
		case <-c.done:
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				return
			}

			var msg Message
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			c.handleMessage(msg)
		}
	}
}

func (c *Client) handleMessage(msg Message) {
	switch msg.Method {
	case "futures:btc_usd:last-price":
		if handler, exists := c.handlers["price"]; exists {
			if data, err := json.Marshal(msg.Result); err == nil {
				handler(data)
			}
		}
	case "futures:btc_usd:index":
		if handler, exists := c.handlers["index"]; exists {
			if data, err := json.Marshal(msg.Result); err == nil {
				handler(data)
			}
		}
	default:
		log.Printf("Received message: %s", msg.Method)
	}
}

func (c *Client) heartbeat() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.done:
			return
		case <-ticker.C:
			c.mu.Lock()
			if c.conn != nil {
				if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("Error sending ping: %v", err)
				}
			}
			c.mu.Unlock()
		}
	}
}

func (c *Client) Subscribe(event string) error {
	msg := Message{
		JSONRPC: "2.0",
		Method:  "subscribe",
		ID:      fmt.Sprintf("sub_%d", time.Now().Unix()),
		Params: map[string]string{
			"event": event,
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return c.conn.WriteMessage(websocket.TextMessage, data)
	}

	return fmt.Errorf("not connected")
}

func (c *Client) OnPriceUpdate(handler func([]byte)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers["price"] = handler
}

func (c *Client) OnIndexUpdate(handler func([]byte)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers["index"] = handler
}

func (c *Client) SendEcho(message string) error {
	msg := Message{
		JSONRPC: "2.0",
		Method:  "debug/echo",
		ID:      fmt.Sprintf("echo_%d", time.Now().Unix()),
		Params: map[string]string{
			"hello": message,
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return c.conn.WriteMessage(websocket.TextMessage, data)
	}

	return fmt.Errorf("not connected")
}
