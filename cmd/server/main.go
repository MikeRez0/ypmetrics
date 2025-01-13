package main

import (
	"fmt"
	"log"

	"github.com/MikeRez0/ypmetrics/internal/server"
)

var buildVersion string
var buildDate string
var buildCommit string

const cBuildInfoTemplate = `Build version: %s
Build date: %s
Build commit: %s
`

func main() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Printf(cBuildInfoTemplate, buildVersion, buildDate, buildCommit)

	if err := server.Run(); err != nil {
		// no custom logger at this line
		log.Fatalf("Fatal error: %v", err)
	}
}
