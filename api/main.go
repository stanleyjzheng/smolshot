package main

import (
	"log"
	"net/http"

	ohlc "smolshot_api/lib"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/v3/coins/{symbol}/ohlc", ohlc.Handler).Methods("GET")

	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
