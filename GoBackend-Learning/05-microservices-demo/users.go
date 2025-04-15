package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Users microservice: list of users")
	})
	fmt.Println("Users service running at http://localhost:8085/users")
	log.Fatal(http.ListenAndServe(":8085", nil))
}
