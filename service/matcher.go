package service

import (
    "order-matching-system/db"
    "order-matching-system/models"
    "errors"
)

var orderBooks = make(map[string][]models.Order)

func PlaceOrder(order models.Order) (models.Order, error) {
    order.Remaining = order.Quantity
    result, err := db.DB.Exec(`INSERT INTO orders (symbol, side, type, price, quantity, remaining, status) VALUES (?, ?, ?, ?, ?, ?, ?)`,
        order.Symbol, order.Side, order.Type, order.Price, order.Quantity, order.Remaining, "open")
    if err != nil {
        return order, err
    }
    order.ID, _ = result.LastInsertId()
    matchOrder(order)
    return order, nil
}

func CancelOrder(id int64) error {
    res, err := db.DB.Exec(`UPDATE orders SET status='canceled' WHERE id=? AND status='open'`, id)
    rows, _ := res.RowsAffected()
    if rows == 0 {
        return errors.New("not found or already filled")
    }
    return err
}

func GetOrderBook(symbol string) []models.Order {
    rows, _ := db.DB.Query(`SELECT id, symbol, side, type, price, quantity, remaining, status FROM orders WHERE symbol=? AND status='open' ORDER BY created_at ASC`, symbol)
    defer rows.Close()

    var orders []models.Order
    for rows.Next() {
        var o models.Order
        rows.Scan(&o.ID, &o.Symbol, &o.Side, &o.Type, &o.Price, &o.Quantity, &o.Remaining, &o.Status)
        orders = append(orders, o)
    }
    return orders
}

func ListTrades(symbol string) []models.Trade {
    rows, _ := db.DB.Query(`SELECT id, buy_order_id, sell_order_id, symbol, price, quantity FROM trades WHERE symbol=? ORDER BY executed_at DESC`, symbol)
    defer rows.Close()

    var trades []models.Trade
    for rows.Next() {
        var t models.Trade
        rows.Scan(&t.ID, &t.BuyOrderID, &t.SellOrderID, &t.Symbol, &t.Price, &t.Quantity)
        trades = append(trades, t)
    }
    return trades
}

func matchOrder(order models.Order) {
    opposite := "buy"
    if order.Side == "buy" {
        opposite = "sell"
    }

    rows, err := db.DB.Query(`
        SELECT id, price, remaining FROM orders
        WHERE symbol = ? AND side = ? AND status = 'open'
        AND (type = 'limit' OR type = 'market')
        ORDER BY price ASC, created_at ASC
    `, order.Symbol, opposite)
    if err != nil {
        return
    }
    defer rows.Close()

    for rows.Next() {
        var existingID int64
        var existingPrice, existingRemaining float64
        rows.Scan(&existingID, &existingPrice, &existingRemaining)

        match := false
        if order.Type == "market" || (order.Type == "limit" && (
            (order.Side == "buy" && order.Price >= existingPrice) ||
            (order.Side == "sell" && order.Price <= existingPrice))) {
            match = true
        }

        if !match {
            break
        }

        tradeQty := min(order.Remaining, existingRemaining)
        tradePrice := existingPrice

        db.DB.Exec(`INSERT INTO trades (buy_order_id, sell_order_id, symbol, price, quantity) VALUES (?, ?, ?, ?, ?)`,
            ifThen(order.Side == "buy", order.ID, existingID),
            ifThen(order.Side == "buy", existingID, order.ID),
            order.Symbol, tradePrice, tradeQty)

        db.DB.Exec(`UPDATE orders SET remaining = remaining - ?, status = IF(remaining - ? <= 0, 'filled', 'open') WHERE id = ?`,
            tradeQty, tradeQty, existingID)
        order.Remaining -= tradeQty

        if order.Remaining <= 0 {
            db.DB.Exec(`UPDATE orders SET remaining = 0, status = 'filled' WHERE id = ?`, order.ID)
            return
        }
    }

    // If limit order remains unmatched
    if order.Type == "limit" && order.Remaining > 0 {
        db.DB.Exec(`UPDATE orders SET remaining = ?, status = 'open' WHERE id = ?`, order.Remaining, order.ID)
    }

    // If market order has no more matches
    if order.Type == "market" && order.Remaining > 0 {
        db.DB.Exec(`UPDATE orders SET remaining = ?, status = 'filled' WHERE id = ?`, order.Remaining, order.ID)
    }
}

func min(a, b float64) float64 {
    if a < b {
        return a
    }
    return b
}

func ifThen(cond bool, a, b int64) int64 {
    if cond {
        return a
    }
    return b
}

