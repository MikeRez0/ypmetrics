package main

import (
	"log"

	"github.com/MikeRez0/ypmetrics/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		// no custom logger at this line
		log.Fatalf("Fatal error: %v", err)
	}
}
