package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type ConfigServer struct {
	HostString string `env:"ADDRESS"`
}

func NewConfigServer() (*ConfigServer, error) {
	// null config
	config := ConfigServer{}

	// cmd string params
	flag.StringVar(&config.HostString, "a", `localhost:8080`, "HTTP server endpoint")
	flag.Parse()

	// environment override
	err := env.Parse(&config)
	if err != nil {
		return nil, fmt.Errorf("error parsing env config: %w", err)
	}

	return &config, nil
}

type ConfigAgent struct {
	HostString     string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func NewConfigAgent() (*ConfigAgent, error) {
	// null config
	config := ConfigAgent{}

	// cmd string params
	flag.StringVar(&config.HostString, "a", `localhost:8080`, "HTTP server endpoint")
	flag.IntVar(&config.PollInterval, "p", 2, "Poll interval")
	flag.IntVar(&config.ReportInterval, "r", 10, "Report interval")
	flag.Parse()

	// environment override
	err := env.Parse(&config)
	if err != nil {
		return nil, fmt.Errorf("error parsing env config: %w", err)
	}

	return &config, nil
}
