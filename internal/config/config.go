package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

// ConfigServer - config params for server.
type ConfigServer struct { //nolint:govet //no need for opimization
	HostString      string `env:"ADDRESS"`
	LogLevel        string `env:"LOG_LEVEL"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DSN             string `env:"DATABASE_DSN"`
	SignKey         string `env:"KEY"`
}

// NewConfigServer - parse and create new server config.
func NewConfigServer() (*ConfigServer, error) {
	// null config
	config := ConfigServer{}

	// cmd string params
	flag.StringVar(&config.HostString, "a", `localhost:8080`, "HTTP server endpoint")
	flag.StringVar(&config.LogLevel, "l", `debug`, "Log level")
	flag.IntVar(&config.StoreInterval, "i", 300, "File store interval, 0 - synchrose")
	flag.StringVar(&config.FileStoragePath, "f", "/tmp/metrics-db.json", "File store path, empty - without store")
	flag.BoolVar(&config.Restore, "r", true, "Needs restore on start")
	flag.StringVar(&config.DSN, "d", "", "Database string")
	flag.StringVar(&config.SignKey, "k", "", "SighHash Key")
	flag.Parse()

	// environment override
	err := env.Parse(&config)
	if err != nil {
		return nil, fmt.Errorf("error parsing env config: %w", err)
	}

	return &config, nil
}

// ConfigAgent - config params for agent.
type ConfigAgent struct {
	HostString     string `env:"ADDRESS"`
	SignKey        string `env:"KEY"`
	LogLevel       string `env:"LOG_LEVEL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

// Parse and create new agent config
func NewConfigAgent() (*ConfigAgent, error) {
	// null config
	config := ConfigAgent{}

	// cmd string params
	flag.StringVar(&config.HostString, "a", `localhost:8080`, "HTTP server endpoint")
	flag.IntVar(&config.PollInterval, "p", 2, "Poll interval")
	flag.IntVar(&config.ReportInterval, "r", 10, "Report interval")
	flag.IntVar(&config.RateLimit, "l", 3, "Rate limit")
	flag.StringVar(&config.SignKey, "k", "", "SighHash Key")
	flag.StringVar(&config.LogLevel, "log", `error`, "Log level")
	flag.Parse()

	// environment override
	err := env.Parse(&config)
	if err != nil {
		return nil, fmt.Errorf("error parsing env config: %w", err)
	}

	return &config, nil
}
