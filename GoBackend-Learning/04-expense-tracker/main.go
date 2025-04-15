package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Expense struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Amount float64 `json:"amount"`
}

var db *sql.DB

func getExpenses(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, amount FROM expenses")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var expenses []Expense
	for rows.Next() {
		var e Expense
		rows.Scan(&e.ID, &e.Title, &e.Amount)
		expenses = append(expenses, e)
	}
	json.NewEncoder(w).Encode(expenses)
}

func main() {
	var err error
	connStr := "user=postgres password=postgres dbname=expenses sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/expenses", getExpenses)
	fmt.Println("Expense Tracker API running at http://localhost:8083/")
	log.Fatal(http.ListenAndServe(":8083", nil))
}
