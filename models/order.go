package models

type Order struct {
    ID        int64   `json:"id"`
    Symbol    string  `json:"symbol"`
    Side      string  `json:"side"`
    Type      string  `json:"type"`
    Price     float64 `json:"price"`
    Quantity  float64 `json:"quantity"`
    Remaining float64 `json:"remaining"`
    Status    string  `json:"status"`
}
