package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func InitDB() {
	godotenv.Load()

	var err error
	db, err = sqlx.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database successfully")
}

func SetPrivateKeyHandler(w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		UserID     string `json:"user_id"`
		PrivateKey string `json:"private_key"`
	}

	type responsePayload struct {
		PublicKey string `json:"public_key"`
	}

	var payload requestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	keypair, err := types.AccountFromBase58(payload.PrivateKey)
	if err != nil {
		http.Error(w, "Invalid private key", http.StatusBadRequest)
		return
	}

	// Insert or update the user in the accounts table
	query := `INSERT INTO accounts (user_id, private_key, public_key) VALUES (:user_id, :private_key, :public_key)
		ON CONFLICT (user_id) DO UPDATE SET private_key = EXCLUDED.private_key, public_key = EXCLUDED.public_key`
	params := map[string]interface{}{
		"user_id":     payload.UserID,
		"private_key": payload.PrivateKey,
		"public_key":  keypair.PublicKey.ToBase58(),
	}
	_, err = db.NamedExec(query, params)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	response := responsePayload{PublicKey: keypair.PublicKey.ToBase58()}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}

	// Fetch the user's public key from the accounts table
	var publicKey string
	query := `SELECT public_key FROM accounts WHERE user_id = $1`
	err := db.Get(&publicKey, query, userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Connect to Solana RPC to fetch the balance
	c := client.NewClient(os.Getenv("SOLANA_RPC_URL"))
	resp, err := c.GetBalance(r.Context(), publicKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch balance: %v", err), http.StatusInternalServerError)
		return
	}

	balance := fmt.Sprintf("%f SOL", float64(resp)/1e9) // Convert lamports to SOL

	response := map[string]string{
		"user_id":    userID,
		"public_key": publicKey,
		"balance":    balance,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
