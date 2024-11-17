package main

import (
	"log"
	"net/http"

	handlers "smolshot_api/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	handlers.InitDB()

	r.HandleFunc("/api/v3/coins/{pool}/ohlc", handlers.Handler).Methods("GET")
	r.HandleFunc("/api/v3/set_private_key", handlers.SetPrivateKeyHandler).Methods("POST")
	r.HandleFunc("/api/v3/get_sol_balance", handlers.GetSolBalanceHandler).Methods("GET")
	r.HandleFunc("/api/v3/get_token_balance", handlers.GetTokenBalanceHandler).Methods("GET")

	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
