// handlers/submit.go
package handlers

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/mr-tron/base58"
	"github.com/weeaa/jito-go/clients/searcher_client"
)

// AddInstructionToTransaction adds a new instruction to an existing transaction.
// It preserves the existing instructions and appends the new one.
func AddInstructionToTransaction(existingTx *solana.Transaction, newInstr solana.Instruction) (*solana.Transaction, error) {
	// Initialize the TransactionBuilder
	builder := solana.NewTransactionBuilder()

	// Add existing instructions to the builder
	for _, compiledInstr := range existingTx.Message.Instructions {
		// Convert CompiledInstruction to Instruction interface
		instr, err := DecodeInstruction(compiledInstr, &existingTx.Message)
		if err != nil {
			return nil, fmt.Errorf("failed to decode existing instruction: %w", err)
		}
		builder.AddInstruction(instr)
	}

	// Add the new instruction
	builder.AddInstruction(newInstr)

	// Set the recent blockhash from the existing transaction
	builder.SetRecentBlockHash(existingTx.Message.RecentBlockhash)

	// Optionally, set the fee payer if different from the existing one
	// builder.SetFeePayer(newFeePayerPublicKey)

	// Build the new transaction
	newTx, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build new transaction: %w", err)
	}

	return newTx, nil
}

// DecodeInstruction decodes a CompiledInstruction into an Instruction.
func DecodeInstruction(compiledInstr solana.CompiledInstruction, message *solana.Message) (solana.Instruction, error) {
	// Decode the instruction data
	data, err := base58.Decode(string(compiledInstr.Data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode instruction data: %w", err)
	}

	// Create a new Instruction
	instr := &CustomInstruction{
		programID: message.AccountKeys[compiledInstr.ProgramIDIndex],
		accounts:  convertAccountMetas(message, compiledInstr.Accounts),
		data:      data,
	}

	return instr, nil
}

// CustomInstruction implements the solana.Instruction interface.
type CustomInstruction struct {
	programID solana.PublicKey
	accounts  []*solana.AccountMeta
	data      []byte
}

func (ci *CustomInstruction) ProgramID() solana.PublicKey {
	return ci.programID
}

func (ci *CustomInstruction) Accounts() []*solana.AccountMeta {
	return ci.accounts
}

func (ci *CustomInstruction) Data() ([]byte, error) {
	return ci.data, nil
}

func convertAccountMetas(message *solana.Message, accountIndices []uint16) []*solana.AccountMeta {
	var accounts []*solana.AccountMeta
	for _, index := range accountIndices {
		accounts = append(accounts, &solana.AccountMeta{
			PublicKey:  message.AccountKeys[index],
			IsSigner:   index < uint16(message.Header.NumRequiredSignatures),
			IsWritable: index >= uint16(message.Header.NumRequiredSignatures)-uint16(message.Header.NumReadonlySignedAccounts),
		})
	}
	return accounts
}

// SwapTokenHandler handles the token swap request
func SwapTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Define request and response structures
	type SwapRequest struct {
		UserID      string `json:"user_id"`
		InputMint   string `json:"input_mint"`
		OutputMint  string `json:"output_mint"`
		Amount      uint64 `json:"amount"`       // in smallest unit (e.g., lamports)
		SlippageBps int    `json:"slippage_bps"` // Basis points (e.g., 100 = 1%)
	}

	type SwapResponse struct {
		TxHash string `json:"txHash"`
	}

	// Decode the incoming JSON request
	var swapReq SwapRequest
	if err := json.NewDecoder(r.Body).Decode(&swapReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required parameters
	if swapReq.UserID == "" || swapReq.InputMint == "" || swapReq.OutputMint == "" || swapReq.Amount == 0 {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Initialize environment variables
	if err := godotenv.Load(); err != nil {
		// Proceed if .env is not present; environment variables might be set otherwise
	}

	jitoRpcUrl := os.Getenv("JITO_RPC")
	if jitoRpcUrl == "" {
		http.Error(w, "JITO_RPC not set in environment", http.StatusInternalServerError)
		return
	}

	jitoBEUrl := os.Getenv("JITO_BLOCK_ENGINE_URL")
	if jitoBEUrl == "" {
		http.Error(w, "JITO_BLOCK_ENGINE_URL not set in environment", http.StatusInternalServerError)
		return
	}

	// Fetch the user's private key from the database
	var privateKeyHex string
	query := `SELECT private_key FROM accounts WHERE user_id = $1`
	err := db.Get(&privateKeyHex, query, swapReq.UserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Decode the private key from hex string
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		http.Error(w, "Invalid private key format", http.StatusInternalServerError)
		return
	}

	if len(privateKeyBytes) != 64 {
		http.Error(w, "Invalid private key length", http.StatusInternalServerError)
		return
	}

	// Convert hex string to base58 private key
	privateKeyBase58 := base58.Encode(privateKeyBytes)
	privateKey, err := solana.PrivateKeyFromBase58(privateKeyBase58)
	if err != nil {
		http.Error(w, "Failed to create private key", http.StatusInternalServerError)
		return
	}
	publicKey := privateKey.PublicKey()

	keyMap := make(map[solana.PublicKey]solana.PrivateKey)
	keyMap[publicKey] = privateKey

	// Initialize Jito client
	ctx := context.TODO()
	solanaRpcURL := os.Getenv("SOLANA_RPC_URL")
	if solanaRpcURL == "" {
		http.Error(w, "SOLANA_RPC_URL not set in environment", http.StatusInternalServerError)
		return
	}
	jitoClient, err := searcher_client.NewNoAuth(
		ctx,
		jitoBEUrl,
		rpc.New(solanaRpcURL),
		rpc.New(jitoRpcUrl),
		nil,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to initialize Jito client: %v", err), http.StatusInternalServerError)
		return
	}
	defer jitoClient.Close()

	// Step 1: Get a quote from Jupiter
	baseURL := "https://quote-api.jup.ag/v4/quote"
	quoteURL := fmt.Sprintf("%s?inputMint=%s&outputMint=%s&amount=%d&slippageBps=%d",
		baseURL, swapReq.InputMint, swapReq.OutputMint, swapReq.Amount, swapReq.SlippageBps)

	resp, err := http.Get(quoteURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get quote: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Failed to get quote: status code %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	var quoteResponse struct {
		Data []json.RawMessage `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&quoteResponse)
	if err != nil {
		http.Error(w, "Failed to parse quote response", http.StatusInternalServerError)
		return
	}

	if len(quoteResponse.Data) == 0 {
		http.Error(w, "No quote data returned", http.StatusInternalServerError)
		return
	}

	// Get the first route
	routeData := quoteResponse.Data[0]

	// Step 2: Get the swap transaction from Jupiter
	swapRequestBody := struct {
		Route         json.RawMessage `json:"route"`
		UserPublicKey string          `json:"userPublicKey"`
		WrapUnwrapSOL bool            `json:"wrapUnwrapSOL"`
	}{
		Route:         routeData,
		UserPublicKey: publicKey.String(),
		WrapUnwrapSOL: false,
	}

	swapReqBytes, err := json.Marshal(swapRequestBody)
	if err != nil {
		http.Error(w, "Failed to marshal swap request", http.StatusInternalServerError)
		return
	}

	swapURL := "https://quote-api.jup.ag/v4/swap"
	reqSwap, err := http.NewRequest("POST", swapURL, bytes.NewBuffer(swapReqBytes))
	if err != nil {
		http.Error(w, "Failed to create swap request", http.StatusInternalServerError)
		return
	}
	reqSwap.Header.Set("Content-Type", "application/json")

	swapResp, err := http.DefaultClient.Do(reqSwap)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get swap transaction: %v", err), http.StatusInternalServerError)
		return
	}
	defer swapResp.Body.Close()

	if swapResp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(swapResp.Body)
		http.Error(w, fmt.Sprintf("Failed to get swap transaction: status code %d, body: %s", swapResp.StatusCode, string(bodyBytes)), http.StatusInternalServerError)
		return
	}

	var swapResponse struct {
		SwapTransaction string `json:"swapTransaction"`
	}

	err = json.NewDecoder(swapResp.Body).Decode(&swapResponse)
	if err != nil {
		http.Error(w, "Failed to parse swap response", http.StatusInternalServerError)
		return
	}

	if swapResponse.SwapTransaction == "" {
		http.Error(w, "Swap transaction not provided", http.StatusInternalServerError)
		return
	}

	// Step 3: Decode, deserialize, and modify the transaction
	tx, err := solana.TransactionFromBase64(swapResponse.SwapTransaction)
	if err != nil {
		http.Error(w, "Failed to deserialize transaction", http.StatusInternalServerError)
		return
	}

	// Step 4: Add Jito tip instruction
	tipAmount := uint64(10000) // Tip amount in lamports
	tipInst, err := jitoClient.GenerateTipRandomAccountInstruction(tipAmount, publicKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate tip instruction: %v", err), http.StatusInternalServerError)
		return
	}

	// Append the tip instruction to the transaction
	newTx, err := AddInstructionToTransaction(tx, tipInst)

	// Step 5: Update recent blockhash
	blockhashResp, err := jitoClient.RpcConn.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get recent blockhash: %v", err), http.StatusInternalServerError)
		return
	}
	newTx.Message.RecentBlockhash = blockhashResp.Value.Blockhash

	// Step 6: Sign the transaction
	// Clear existing signatures
	newTx.Signatures = nil

	// Sign the transaction with the user's private key
	privateKeyGetter := func(key solana.PublicKey) *solana.PrivateKey {
		if privKey, exists := keyMap[key]; exists {
			return &privKey
		}
		return nil
	}
	_, err = newTx.Sign(privateKeyGetter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to sign transaction: %v", err), http.StatusInternalServerError)
		return
	}

	// Step 7: Submit the transaction via Jito client
	txns := []*solana.Transaction{newTx}

	// Broadcast the bundle with confirmation
	respBundle, err := jitoClient.BroadcastBundle(txns)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to broadcast transaction via Jito: %v", err), http.StatusInternalServerError)
		return
	}

	// Ensure respBundle has at least one transaction
	// Adjust the following based on the actual structure of SendBundleResponse

	// print
	fmt.Println(respBundle)

	// Prepare and send the response
	swapRespPayload := SwapResponse{
		TxHash: "meow",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(swapRespPayload)
}
