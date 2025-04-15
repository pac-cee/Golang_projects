# BTC Dashboard Backend (Go)

This Go backend provides real-time BTC/USDT price and candlestick data by proxying the Binance API. It is designed for use with a modern frontend dashboard.

## Endpoints
- `GET /price` — Current BTC/USDT price
- `GET /candles` — Recent 1-minute candlesticks (last 50)

## How it Works
- Fetches data from Binance public API
- Enables CORS for local frontend development
- Designed for easy integration with React/TS frontend

## How to Run
```sh
cd backend
# Install dependencies if needed
 go mod tidy
# Run the server
 go run main.go
# Server runs at http://localhost:8088
```

## Environment
- Go 1.20+
- No authentication required (public data)

---

See the frontend folder for the React dashboard.
