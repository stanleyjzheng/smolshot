package main

import (
	"log"
	"net/http"

	lib "smolshot_api/lib"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	lib.InitDB()

	r.HandleFunc("/api/v3/coins/{pool}/ohlc", lib.Handler).Methods("GET")
	r.HandleFunc("/api/v3/set_private_key", lib.SetPrivateKeyHandler).Methods("POST")
	r.HandleFunc("/api/v3/get_balance", lib.GetBalanceHandler).Methods("GET")

	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
