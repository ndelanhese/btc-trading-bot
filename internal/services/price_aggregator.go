package services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// PriceSnapshot represents a single aggregated price observation.
type PriceSnapshot struct {
	Price     float64            `json:"price"`
	Timestamp int64              `json:"timestamp"`
	Sources   map[string]float64 `json:"sources"`
}

// PriceAggregator periodically fetches BTC/USD prices from multiple exchanges
// and maintains an averaged price for subscribers.
type PriceAggregator struct {
	httpClient    *http.Client
	pollInterval  time.Duration
	stopChan      chan struct{}
	subscribersMu sync.Mutex
	subscribers   map[chan PriceSnapshot]struct{}
	latestAtomic  atomic.Value // stores PriceSnapshot
}

func NewPriceAggregator() *PriceAggregator {
	timeoutStr := os.Getenv("AGG_HTTP_TIMEOUT_MS")
	timeout := 1200 * time.Millisecond
	if timeoutStr != "" {
		if v, err := strconv.Atoi(timeoutStr); err == nil && v > 0 {
			timeout = time.Duration(v) * time.Millisecond
		}
	}

	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     30 * time.Second,
			TLSHandshakeTimeout: 3 * time.Second,
		},
	}

	agg := &PriceAggregator{
		httpClient:   client,
		pollInterval: time.Second,
		stopChan:     make(chan struct{}),
		subscribers:  make(map[chan PriceSnapshot]struct{}),
	}
	agg.latestAtomic.Store(PriceSnapshot{})
	return agg
}

func (p *PriceAggregator) Start() {
	ticker := time.NewTicker(p.pollInterval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-p.stopChan:
				return
			case <-ticker.C:
				snapshot, ok := p.fetchAndAggregate()
				if ok {
					p.latestAtomic.Store(snapshot)
					p.broadcast(snapshot)
				}
			}
		}
	}()
}

func (p *PriceAggregator) Stop() {
	close(p.stopChan)
}

// Latest returns the most recent snapshot.
func (p *PriceAggregator) Latest() PriceSnapshot {
	val := p.latestAtomic.Load()
	if val == nil {
		return PriceSnapshot{}
	}
	return val.(PriceSnapshot)
}

// Subscribe returns a channel that will receive future snapshots and an unsubscribe function.
func (p *PriceAggregator) Subscribe() (<-chan PriceSnapshot, func()) {
	ch := make(chan PriceSnapshot, 1)
	p.subscribersMu.Lock()
	p.subscribers[ch] = struct{}{}
	p.subscribersMu.Unlock()

	// Send latest immediately if present
	if snap := p.Latest(); snap.Timestamp != 0 {
		select {
		case ch <- snap:
		default:
		}
	}

	unsubscribe := func() {
		p.subscribersMu.Lock()
		if _, ok := p.subscribers[ch]; ok {
			delete(p.subscribers, ch)
			close(ch)
		}
		p.subscribersMu.Unlock()
	}
	return ch, unsubscribe
}

func (p *PriceAggregator) broadcast(snap PriceSnapshot) {
	p.subscribersMu.Lock()
	defer p.subscribersMu.Unlock()
	for ch := range p.subscribers {
		select {
		case ch <- snap:
		default:
			// drop if slow consumer
		}
	}
}

func (p *PriceAggregator) fetchAndAggregate() (PriceSnapshot, bool) {
	// Context to bound each poll cycle
	ctx, cancel := context.WithTimeout(context.Background(), p.httpClient.Timeout)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(3)

	type result struct {
		name  string
		price float64
		err   error
	}

	results := make(chan result, 3)

	go func() {
		defer wg.Done()
		price, err := p.fetchBinance(ctx)
		results <- result{name: "binance", price: price, err: err}
	}()
	go func() {
		defer wg.Done()
		price, err := p.fetchCoinbase(ctx)
		results <- result{name: "coinbase", price: price, err: err}
	}()
	go func() {
		defer wg.Done()
		price, err := p.fetchKraken(ctx)
		results <- result{name: "kraken", price: price, err: err}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	sum := 0.0
	count := 0
	sources := make(map[string]float64, 3)

	for r := range results {
		if r.err == nil && r.price > 0 {
			sum += r.price
			count++
			sources[r.name] = r.price
		}
	}

	if count == 0 {
		return PriceSnapshot{}, false
	}

	avg := sum / float64(count)
	return PriceSnapshot{
		Price:     avg,
		Timestamp: time.Now().UnixMilli(),
		Sources:   sources,
	}, true
}

func (p *PriceAggregator) fetchBinance(ctx context.Context) (float64, error) {
	// https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT", nil)
	if err != nil {
		return 0, err
	}
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("binance non-200")
	}
	var payload struct {
		Price string `json:"price"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0, err
	}
	return strconv.ParseFloat(payload.Price, 64)
}

func (p *PriceAggregator) fetchCoinbase(ctx context.Context) (float64, error) {
	// https://api.coinbase.com/v2/prices/spot?currency=USD
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.coinbase.com/v2/prices/spot?currency=USD", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", "btc-trading-bot/1.0")
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("coinbase non-200")
	}
	var payload struct {
		Data struct {
			Amount string `json:"amount"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0, err
	}
	return strconv.ParseFloat(payload.Data.Amount, 64)
}

func (p *PriceAggregator) fetchKraken(ctx context.Context) (float64, error) {
	// https://api.kraken.com/0/public/Ticker?pair=XBTUSD
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.kraken.com/0/public/Ticker?pair=XBTUSD", nil)
	if err != nil {
		return 0, err
	}
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("kraken non-200")
	}
	var payload struct {
		Result map[string]struct {
			C []string `json:"c"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0, err
	}
	for _, v := range payload.Result {
		if len(v.C) > 0 {
			return strconv.ParseFloat(v.C[0], 64)
		}
	}
	return 0, errors.New("kraken malformed")
}
