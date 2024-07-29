package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type ConfigServer struct { //nolint:govet //no need for opimization
	HostString      string `env:"ADDRESS"`
	LogLevel        string `env:"LOG_LEVEL"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DSN             string `env:"DATABASE_DSN"`
}

func NewConfigServer() (*ConfigServer, error) {
	// null config
	config := ConfigServer{}

	// cmd string params
	flag.StringVar(&config.HostString, "a", `localhost:8080`, "HTTP server endpoint")
	flag.StringVar(&config.LogLevel, "l", `error`, "Log level")
	flag.IntVar(&config.StoreInterval, "i", 300, "File store interval, 0 - synchrose")
	flag.StringVar(&config.FileStoragePath, "f", "/tmp/metrics-db.json", "File store path, empty - without store")
	flag.BoolVar(&config.Restore, "r", true, "Needs restore on start")
	flag.StringVar(&config.DSN, "d", "", "Database string")
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
