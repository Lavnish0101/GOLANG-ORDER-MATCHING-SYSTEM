# 🏦 Golang Order Matching System

A simplified stock exchange-style backend in **Golang** with **MySQL** that supports:
- Limit & Market Orders
- Order Matching Engine (price-time priority)
- REST API for placing, canceling, viewing orders and trades

---

## 🧰 Tech Stack

- **Language**: Go (1.20+)
- **Database**: MySQL (or TiDB - MySQL compatible)
- **HTTP Router**: [gorilla/mux](https://github.com/gorilla/mux)
- **MySQL Driver**: [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)

---

## ⚙️ Dependencies & Setup

### 🔧 Install Golang

Install Go from: https://golang.org/dl/  
Ensure it’s added to your system PATH.

```bash
go version
```

🐬 Install MySQL
Install MySQL (or TiDB) and ensure it’s running.

Login to MySQL:

```bash
mysql -u root -p
```

Then run:
```sql
CREATE DATABASE orderbook;
```

🗃️ Database Initialization
From the project root:

```bash
mysql -u root -p orderbook < db/schema.sql
```

This will create two tables:

-> orders: Tracks all submitted orders (limit/market, buy/sell)

-> trades: Tracks executed trades between matched orders

🚀 Starting the Server
Step 1: Configure DB DSN
In cmd/main.go update:
```go
db.InitDB("root:<your-password>@tcp(localhost:3306)/orderbook")
```

Step 2: Install Modules
```bash
go mod tidy
```

Step 3: Run the Server
```bash
go run cmd/main.go
```

Server runs at: http://localhost:8080

🧪 API Endpoints

All APIs use JSON format.

🔹 POST /orders — Place Order
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "AAPL",
    "side": "buy",
    "type": "limit",
    "price": 150.00,
    "quantity": 10
}'
```
Fields:

symbol: Stock symbol (e.g., "AAPL")

side: "buy" or "sell"

type: "limit" or "market"

price: only for limit orders

quantity: total units

🔹 DELETE /orders/{id} — Cancel Order
```bash
curl -X DELETE http://localhost:8080/orders/1
```
🔹 GET /orderbook?symbol=AAPL — View Order Book
```bash
curl "http://localhost:8080/orderbook?symbol=AAPL"
```
Returns all open orders (not filled or canceled).

🔹 GET /trades?symbol=AAPL — Trade History
```bash
curl "http://localhost:8080/trades?symbol=AAPL"
```


💡 Design Decisions & Assumptions
Only open orders are shown in /orderbook endpoint.

Market orders are matched immediately and never sit in the order book.

Matching uses price-time priority:

Best price first (buy: highest bid, sell: lowest ask)

FIFO within same price level

Matching engine and DB operations are fully transactional

No ORMs used — only raw SQL

Matching is triggered only when a new order is placed

No concurrency control for simultaneous placements (can be extended with mutex or DB locks)

📦 Directory Structure
```bash
.
├── api/           # HTTP handlers & router
├── db/            # DB connection + schema
├── models/        # Data models for Order, Trade
├── service/       # Matching logic
├── cmd/           # Entry point
├── go.mod
└── README.md
```


✅ Author

Built with ❤️ by [Lavnish Kumar](https://github.com/Lavnish0101)
