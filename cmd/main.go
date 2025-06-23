package main

import (
    "log"
    "net/http"
    "order-matching-system/api"
    "order-matching-system/db"
)

func main() {
    db.InitDB("root:Lavi.572353@tcp(localhost:3306)/orderbook")

    r := api.SetupRouter()
    log.Println("Listening on :8080")
    http.ListenAndServe(":8080", r)
}
