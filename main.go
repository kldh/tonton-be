package main

import (
	"log"
	"os"
)

func main() {
	cp := os.Getenv("CONFIG_PATH")
	if cp == "" {
		log.Fatalf("Config path is empty")
	}
}
