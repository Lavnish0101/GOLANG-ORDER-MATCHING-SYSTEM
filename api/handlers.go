package api

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "order-matching-system/models"
    "order-matching-system/service"
)

func PlaceOrderHandler(w http.ResponseWriter, r *http.Request) {
    var order models.Order
    if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    createdOrder, err := service.PlaceOrder(order)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(createdOrder)
}

func CancelOrderHandler(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
    if err := service.CancelOrder(id); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func OrderBookHandler(w http.ResponseWriter, r *http.Request) {
    symbol := r.URL.Query().Get("symbol")
    book := service.GetOrderBook(symbol)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(book)
}

func TradesHandler(w http.ResponseWriter, r *http.Request) {
    symbol := r.URL.Query().Get("symbol")
    trades := service.ListTrades(symbol)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(trades)
}
