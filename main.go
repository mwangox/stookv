package main

import (
	"log"
	"stoo-kv/cmd"
)

func main() {
	if err := cmd.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
