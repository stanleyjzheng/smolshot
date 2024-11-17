package ohlc

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/mux"
)

// Candle represents OHLC data for a given time period
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

// fetchHistory retrieves historical price data from the Mobula API
func fetchHistory(symbol, period string) ([]PriceData, error) {
	url := fmt.Sprintf("https://api.mobula.io/api/1/market/history?asset=%s&blockchain=solana&period=%s", symbol, period)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var history HistoryResponse
	if err := json.Unmarshal(body, &history); err != nil {
		log.Printf("Failed to decode JSON response: %v\nResponse body: %s", err, string(body))
		return nil, err
	}

	var prices []PriceData
	for _, entry := range history.Data.PriceHistory {
		if len(entry) < 2 {
			continue
		}
		timestamp, ok1 := entry[0].(float64)
		price, ok2 := entry[1].(float64)
		if !ok1 || !ok2 {
			log.Printf("Invalid data format: %v", entry)
			continue
		}
		prices = append(prices, PriceData{
			Timestamp: int64(timestamp),
			Price:     price,
		})
	}
	return prices, nil
}

func aggregateOHLC(data []PriceData, interval time.Duration) [][]interface{} {
	if len(data) == 0 {
		return [][]interface{}{}
	}

	// Ensure data is sorted by timestamp
	sort.Slice(data, func(i, j int) bool {
		return data[i].Timestamp < data[j].Timestamp
	})

	var candles [][]interface{}
	start := data[0].Timestamp
	end := start + int64(interval.Milliseconds())
	var open, high, low, close float64
	isFirstPrice := true

	for _, p := range data {
		if p.Timestamp >= start && p.Timestamp < end {
			// Aggregate within the current interval
			if isFirstPrice {
				open = p.Price
				high = p.Price
				low = p.Price
				isFirstPrice = false
			}
			if p.Price > high {
				high = p.Price
			}
			if p.Price < low {
				low = p.Price
			}
			close = p.Price
		} else {
			// Finalize the current candle
			if !isFirstPrice {
				candles = append(candles, []interface{}{start, open, high, low, close})
			}
			// Move to the next interval
			for p.Timestamp >= end {
				start = end
				end += int64(interval.Milliseconds())
			}
			// Initialize the new interval
			open = p.Price
			high = p.Price
			low = p.Price
			close = p.Price
			isFirstPrice = false
		}
	}

	// Finalize the last candle
	if !isFirstPrice {
		candles = append(candles, []interface{}{start, open, high, low, close})
	}

	return candles
}

// Handler processes API requests and responds with OHLC data
func Handler(w http.ResponseWriter, r *http.Request) {
	printf("received ohlc request")

	vars := mux.Vars(r)
	symbol := vars["symbol"]
	interval := r.URL.Query().Get("interval")

	if symbol == "" || interval == "" {
		http.Error(w, "Missing symbol or interval", http.StatusBadRequest)
		return
	}

	var duration time.Duration
	var singleTimestampDuration string
	switch interval {
	case "15min":
		duration = 15 * time.Minute
		singleTimestampDuration = "5min"
	case "1h":
		duration = time.Hour
		singleTimestampDuration = "5min"
	case "4h":
		duration = 4 * time.Hour
		singleTimestampDuration = "15min"
	case "1d":
		duration = 24 * time.Hour
		singleTimestampDuration = "1h"
	default:
		http.Error(w, "Unsupported interval", http.StatusBadRequest)
		return
	}

	data, err := fetchHistory(symbol, singleTimestampDuration)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch data: %v", err), http.StatusInternalServerError)
		return
	}

	if len(data) == 0 {
		http.Error(w, "No data available", http.StatusNotFound)
		return
	}

	candles := aggregateOHLC(data, duration)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(candles); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}
