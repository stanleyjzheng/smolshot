package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// OHLCVResponse represents the response structure from the CoinGecko API
type OHLCVResponse struct {
	Data struct {
		Attributes struct {
			OHLCVList [][]interface{} `json:"ohlcv_list"`
		} `json:"attributes"`
	} `json:"data"`
}

// fetchOHLCV retrieves OHLCV data from the CoinGecko API
func fetchOHLCV(pool string, day string, aggregate string) ([][]interface{}, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	apiKey := os.Getenv("CG_PRO_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API key not found in .env file")
	}
	url := fmt.Sprintf("https://pro-api.coingecko.com/api/v3/onchain/networks/solana/pools/%s/ohlcv/%s?aggregate=%s&currency=usd&limit=1000", pool, day, aggregate)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-cg-pro-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ohlcvResponse OHLCVResponse
	if err := json.Unmarshal(body, &ohlcvResponse); err != nil {
		log.Printf("Failed to decode JSON response: %v\nResponse body: %s", err, string(body))
		return nil, err
	}

	return ohlcvResponse.Data.Attributes.OHLCVList, nil
}

// Handler processes API requests and responds with OHLCV data
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("OHLCV Handler got request to", r.URL)

	vars := mux.Vars(r)
	pool := vars["pool"]
	day := r.URL.Query().Get("period")
	aggregate := r.URL.Query().Get("aggregate")

	if pool == "" || day == "" || aggregate == "" {
		http.Error(w, "Missing pool, day, aggregate", http.StatusBadRequest)
		return
	}

	ohlcvList, err := fetchOHLCV(pool, day, aggregate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch data: %v", err), http.StatusInternalServerError)
		return
	}

	if len(ohlcvList) == 0 {
		http.Error(w, "No data available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ohlcvList); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}
