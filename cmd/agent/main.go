package main

import (
	"log"

	"github.com/MikeRez0/ypmetrics/internal/agent"
)

func main() {
	if err := agent.Run(); err != nil {
		log.Fatal(err)
	}
}
