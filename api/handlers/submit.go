package handlers

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/blocto/solana-go-sdk/types"
	"github.com/ilkamo/jupiter-go/jupiter"
	"github.com/ilkamo/jupiter-go/solana"
	"github.com/mr-tron/base58"
)

// Assuming 'db' is already initialized and accessible in this package
// var db *sqlx.DB

func SwapTokenHandler(w http.ResponseWriter, r *http.Request) {
	type SwapTokenRequest struct {
		UserID      string `json:"user_id"`
		InputMint   string `json:"input_mint"`
		OutputMint  string `json:"output_mint"`
		Amount      int    `json:"amount"` // Changed from uint64 to int
		SlippageBps int    `json:"slippage_bps"`
	}

	type SwapTokenResponse struct {
		TransactionSignature string `json:"transaction_signature"`
	}

	var req SwapTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Basic input validation
	if req.UserID == "" || req.InputMint == "" || req.OutputMint == "" || req.Amount <= 0 {
		http.Error(w, "Missing or invalid input parameters", http.StatusBadRequest)
		return
	}

	// Fetch the user's private key from the accounts table
	var privateKeyHex string
	query := `SELECT private_key FROM accounts WHERE user_id = $1`
	err := db.Get(&privateKeyHex, query, req.UserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Decode the private key from hex
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		http.Error(w, "Invalid private key stored in database", http.StatusInternalServerError)
		return
	}

	// Create the account from the private key
	keypair, err := types.AccountFromBytes(privateKeyBytes)
	if err != nil {
		http.Error(w, "Failed to create account from private key", http.StatusInternalServerError)
		return
	}

	// Create a context
	ctx := context.Background()

	// Initialize the Jupiter client
	jupClient, err := jupiter.NewClientWithResponses(jupiter.DefaultAPIURL)
	if err != nil {
		http.Error(w, "Failed to create Jupiter client", http.StatusInternalServerError)
		return
	}

	// Prepare the slippage Bps
	slippageBps := req.SlippageBps

	// Get the quote from Jupiter
	quoteResponse, err := jupClient.GetQuoteWithResponse(ctx, &jupiter.GetQuoteParams{
		InputMint:   req.InputMint,
		OutputMint:  req.OutputMint,
		Amount:      req.Amount,
		SlippageBps: &slippageBps,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get quote from Jupiter: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare the swap request
	prioritizationFeeLamports := jupiter.SwapRequest_PrioritizationFeeLamports{}
	if err = prioritizationFeeLamports.UnmarshalJSON([]byte(`"auto"`)); err != nil {
		http.Error(w, "Failed to set prioritization fee", http.StatusInternalServerError)
		return
	}

	dynamicComputeUnitLimit := true

	swapResponse, err := jupClient.PostSwapWithResponse(ctx, jupiter.PostSwapJSONRequestBody{
		UserPublicKey:             keypair.PublicKey.ToBase58(),
		QuoteResponse:             *quoteResponse.JSON200,
		PrioritizationFeeLamports: &prioritizationFeeLamports,
		DynamicComputeUnitLimit:   &dynamicComputeUnitLimit,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get swap transaction from Jupiter: %v", err), http.StatusInternalServerError)
		return
	}
	if swapResponse.JSON200 == nil {
		http.Error(w, "No swap transaction received from Jupiter", http.StatusInternalServerError)
		return
	}

	swapTransaction := swapResponse.JSON200.SwapTransaction

	// Create a wallet from the private key

	wallet, err := solana.NewWalletFromPrivateKeyBase58(base58.Encode(privateKeyBytes))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create wallet: %v", err), http.StatusInternalServerError)
		return
	}

	// Create a Solana client using the wallet and RPC URL
	solanaClient, err := solana.NewClient(wallet, os.Getenv("SOLANA_RPC_URL"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create Solana client: %v", err), http.StatusInternalServerError)
		return
	}

	// Send the transaction
	signedTxSignature, err := solanaClient.SendTransactionOnChain(ctx, swapTransaction)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send transaction: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the transaction signature as response
	response := SwapTokenResponse{
		TransactionSignature: string(signedTxSignature),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
