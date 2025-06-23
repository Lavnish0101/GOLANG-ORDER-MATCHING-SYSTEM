package api

import (
    "github.com/gorilla/mux"
    // "net/http"
)

func SetupRouter() *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/orders", PlaceOrderHandler).Methods("POST")
    r.HandleFunc("/orders/{id}", CancelOrderHandler).Methods("DELETE")
    r.HandleFunc("/orderbook", OrderBookHandler).Methods("GET")
    r.HandleFunc("/trades", TradesHandler).Methods("GET")
    return r
}
