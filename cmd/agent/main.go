package main

import (
	"fmt"
	"log"

	"github.com/MikeRez0/ypmetrics/internal/agent"
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

	if err := agent.Run(); err != nil {
		log.Fatal(err)
	}
}
