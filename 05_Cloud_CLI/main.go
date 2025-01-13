package main

import (
	"log"
	"os"

	"cloud-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
