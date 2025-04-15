package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type PriceResponse struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

type Candle struct {
	OpenTime  int64   `json:"open_time"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
	CloseTime int64   `json:"close_time"`
}

// Enable CORS for local dev
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func priceHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	resp, err := http.Get("https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	var data struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}
	json.NewDecoder(resp.Body).Decode(&data)
	price, _ := strconv.ParseFloat(data.Price, 64)
	json.NewEncoder(w).Encode(PriceResponse{Symbol: data.Symbol, Price: price})
}

func candlesHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	// Get last 50 1m candles
	resp, err := http.Get("https://api.binance.com/api/v3/klines?symbol=BTCUSDT&interval=1m&limit=50")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	var arr [][]interface{}
	json.NewDecoder(resp.Body).Decode(&arr)
	var candles []Candle
	for _, v := range arr {
		candles = append(candles, Candle{
			OpenTime:  int64(v[0].(float64)),
			Open:      parseStr(v[1]),
			High:      parseStr(v[2]),
			Low:       parseStr(v[3]),
			Close:     parseStr(v[4]),
			Volume:    parseStr(v[5]),
			CloseTime: int64(v[6].(float64)),
		})
	}
	json.NewEncoder(w).Encode(candles)
}

func parseStr(val interface{}) float64 {
	s, ok := val.(string)
	if !ok {
		return 0
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func main() {
	http.HandleFunc("/price", priceHandler)
	http.HandleFunc("/candles", candlesHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"
	}
	fmt.Printf("BTC API backend running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
