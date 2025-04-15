package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Orders microservice: list of orders")
	})
	fmt.Println("Orders service running at http://localhost:8084/orders")
	log.Fatal(http.ListenAndServe(":8084", nil))
}
