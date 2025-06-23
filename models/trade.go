package models

type Trade struct {
    ID          int64   `json:"id"`
    BuyOrderID  int64   `json:"buy_order_id"`
    SellOrderID int64   `json:"sell_order_id"`
    Symbol      string  `json:"symbol"`
    Price       float64 `json:"price"`
    Quantity    float64 `json:"quantity"`
}
