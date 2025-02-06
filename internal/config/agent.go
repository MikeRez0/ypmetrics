package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
)

// ConfigAgent - config params for agent.
//
// Config file example:
//
//	{
//	    "address": "localhost:8080", // аналог переменной окружения ADDRESS или флага -a
//	    "report_interval": "1s", // аналог переменной окружения REPORT_INTERVAL или флага -r
//	    "poll_interval": "1s", // аналог переменной окружения POLL_INTERVAL или флага -p
//	    "crypto_key": "/path/to/key.pem" // аналог переменной окружения CRYPTO_KEY или флага -crypto-key
//	}
type ConfigAgent struct {
	HostString     string   `env:"ADDRESS" json:"address"`
	SignKey        string   `env:"KEY"`
	LogLevel       string   `env:"LOG_LEVEL"`
	CryptoKey      string   `env:"CRYPTO_KEY" json:"crypto_key"`
	ReportInterval Duration `json:"report_interval"` //env:"REPORT_INTERVAL"
	PollInterval   Duration `json:"poll_interval"`   //env:"POLL_INTERVAL"
	RateLimit      int      `env:"RATE_LIMIT"`
	GRPC           bool     `env:"GRPC_MODE" json:"grpc_mode"`
}

// NewConfigAgent - Parse and create new agent config.
func NewConfigAgent() (*ConfigAgent, error) {
	// null config
	config := ConfigAgent{
		HostString:     `localhost:8080`,
		PollInterval:   Duration{2 * time.Second},
		ReportInterval: Duration{10 * time.Second},
		RateLimit:      3,
		SignKey:        "",
		LogLevel:       "error",
		CryptoKey:      "",
		GRPC:           false,
	}

	err := loadConfigFile(&config)
	if err != nil {
		return nil, fmt.Errorf("error loading config file:%w", err)
	}

	var pollInterval int
	var reportInterval int
	// cmd string params
	flag.String("c", "", cConfigFilenameUsage)
	flag.String("config", "", cConfigFilenameUsage)
	flag.StringVar(&config.HostString, "a", config.HostString, "HTTP/gRPC server endpoint")
	flag.BoolVar(&config.GRPC, "g", config.GRPC, "Enable gRPC Mode")
	flag.IntVar(&pollInterval, "p", -1, "Poll interval")
	flag.IntVar(&reportInterval, "r", -1, "Report interval")
	flag.IntVar(&config.RateLimit, "l", config.RateLimit, "Rate limit")
	flag.StringVar(&config.SignKey, "k", config.SignKey, "SighHash Key")
	flag.StringVar(&config.LogLevel, "log", config.LogLevel, "Log level")
	flag.StringVar(&config.CryptoKey, "crypto-key", config.CryptoKey, "Crypto Key")
	flag.Parse()

	if pollInterval != -1 {
		config.PollInterval = Duration{time.Duration(pollInterval) * time.Second}
	}
	if reportInterval != -1 {
		config.ReportInterval = Duration{time.Duration(reportInterval) * time.Second}
	}

	// environment override
	err = env.Parse(&config)
	if err != nil {
		return nil, fmt.Errorf("error parsing env config: %w", err)
	}

	err = lookupEnvDuration("REPORT_INTERVAL", &config.ReportInterval)
	if err != nil {
		return nil, err
	}

	err = lookupEnvDuration("POLL_INTERVAL", &config.PollInterval)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
