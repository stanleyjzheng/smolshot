package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Candle struct {
	Date  int64   `json:"date"`
	Open  float64 `json:"open"`
	High  float64 `json:"high"`
	Low   float64 `json:"low"`
	Close float64 `json:"close"`
}

// PriceData represents raw API data from Mobula
type PriceData struct {
	Timestamp int64   `json:"timestamp"`
	Price     float64 `json:"price"`
}

type HistoryResponse struct {
	Data struct {
		PriceHistory [][]interface{} `json:"price_history"`
	} `json:"data"`
}

func fetchHistory(symbol, period string) ([]PriceData, error) {
	url := fmt.Sprintf("https://api.mobula.io/api/1/market/history?asset=%s&blockchain=solana&period=%s", symbol, period)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var history HistoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
		return nil, err
	}

	var prices []PriceData
	for _, entry := range history.Data.PriceHistory {
		if len(entry) < 2 {
			continue
		}
		prices = append(prices, PriceData{
			Timestamp: int64(entry[0].(float64)),
			Price:     entry[1].(float64),
		})
	}
	return prices, nil
}

func aggregateOHLC(data []PriceData, interval time.Duration) [][]interface{} {
	// Sort by timestamp
	sort.Slice(data, func(i, j int) bool {
		return data[i].Timestamp < data[j].Timestamp
	})

	candles := [][]interface{}{}
	start := data[0].Timestamp
	end := start + int64(interval.Seconds()*1000)
	var open, high, low, close float64

	for _, p := range data {
		if p.Timestamp >= start && p.Timestamp < end {
			if open == 0 {
				open = p.Price
			}
			if p.Price > high || high == 0 {
				high = p.Price
			}
			if p.Price < low || low == 0 {
				low = p.Price
			}
			close = p.Price
		} else {
			if open != 0 {
				candles = append(candles, []interface{}{start, high, low, open, close})
			}
			// Reset for the new interval
			start = end
			end = start + int64(interval.Seconds()*1000)
			open, high, low, close = p.Price, p.Price, p.Price, p.Price
		}
	}
	if open != 0 {
		candles = append(candles, []interface{}{start, high, low, open, close})
	}
	return candles
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v3/coins/")
	parts := strings.Split(path, "/ohlc")
	if len(parts) != 2 {
		http.Error(w, "Invalid endpoint", http.StatusBadRequest)
		return
	}

	symbol := parts[0]
	interval := r.URL.Query().Get("interval")
	if symbol == "" || interval == "" {
		http.Error(w, "Missing symbol or interval", http.StatusBadRequest)
		return
	}

	var duration time.Duration
	switch interval {
	case "15min":
		duration = 15 * time.Minute
	case "1hr":
		duration = time.Hour
	case "4h":
		duration = 4 * time.Hour
	case "1d":
		duration = 24 * time.Hour
	default:
		http.Error(w, "Unsupported interval", http.StatusBadRequest)
		return
	}

	data, err := fetchHistory(symbol, "5min")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch data: %v", err), http.StatusInternalServerError)
		return
	}

	candles := aggregateOHLC(data, duration)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(candles); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/api/v3/coins/", handler)
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
