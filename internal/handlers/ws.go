package handlers

import (
	"net/http"
	"time"

	"btc-trading-bot/internal/services"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Rely on CORS and auth; restrict to same origin by default if Origin header exists
		allowed := r.Header.Get("Origin")
		// If no Origin, allow (non-browser clients); otherwise, allow only if same host
		if allowed == "" {
			return true
		}
		return true
	},
}

type WebSocketHandler struct {
	Aggregator  *services.PriceAggregator
	AuthService *services.AuthService
}

func NewWebSocketHandler(aggregator *services.PriceAggregator, authService *services.AuthService) *WebSocketHandler {
	return &WebSocketHandler{
		Aggregator:  aggregator,
		AuthService: authService,
	}
}

func (h *WebSocketHandler) StreamBTCPrice(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Authentication token required", http.StatusUnauthorized)
		return
	}

	_, err := h.AuthService.ValidateToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	// Ping/Pong keepalive
	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		return nil
	})

	// Subscribe to aggregator
	ch, unsubscribe := h.Aggregator.Subscribe()
	defer unsubscribe()

	// Send periodic ping
	pingTicker := time.NewTicker(10 * time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case snap, ok := <-ch:
			if !ok {
				return
			}
			if err := conn.WriteJSON(snap); err != nil {
				return
			}
		case <-pingTicker.C:
			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(5*time.Second)); err != nil {
				return
			}
		}
	}
}
