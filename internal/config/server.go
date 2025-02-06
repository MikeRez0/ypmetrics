// Package config contains configuration structures and parsers.
package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
)

// ConfigServer - config params for server.
//
// Config file example:
//
//	{
//	    "address": "localhost:8080", // аналог переменной окружения ADDRESS или флага -a
//	    "restore": true, // аналог переменной окружения RESTORE или флага -r
//	    "store_interval": "1s", // аналог переменной окружения STORE_INTERVAL или флага -i
//	    "store_file": "/path/to/file.db", // аналог переменной окружения STORE_FILE или -f
//	    "database_dsn": "", // аналог переменной окружения DATABASE_DSN или флага -d
//	    "crypto_key": "/path/to/key.pem" // аналог переменной окружения CRYPTO_KEY или флага -crypto-key
//	}
type ConfigServer struct {
	HostString      string   `env:"ADDRESS" json:"address"`
	GRPCHost        string   `env:"GRPC_ADDRESS" json:"grpc_address"`
	LogLevel        string   `env:"LOG_LEVEL"`
	FileStoragePath string   `env:"FILE_STORAGE_PATH" json:"store_file"`
	DSN             string   `env:"DATABASE_DSN" json:"database_dsn"`
	SignKey         string   `env:"KEY"`
	CryptoKey       string   `env:"CRYPTO_KEY" json:"crypto_key"`
	TrustedSubnet   string   `json:"trusted_subnet" env:"TRUSTED_SUBNET"`
	StoreInterval   Duration `json:"store_interval"` //env:"STORE_INTERVAL"
	Restore         bool     `env:"RESTORE" json:"restore"`
}

// NewConfigServer - parse and create new server config.
func NewConfigServer() (*ConfigServer, error) {
	// null config
	config := ConfigServer{
		HostString:      `localhost:8080`,
		GRPCHost:        "",
		LogLevel:        `debug`,
		StoreInterval:   Duration{300 * time.Second},
		FileStoragePath: "",
		Restore:         true,
		DSN:             "",
		SignKey:         "",
		CryptoKey:       "",
		TrustedSubnet:   "",
	}

	err := loadConfigFile(&config)
	if err != nil {
		return nil, fmt.Errorf("error loading config file:%w", err)
	}

	var storeInterval int
	// cmd string params
	flag.String("c", "", cConfigFilenameUsage)
	flag.String("config", "", cConfigFilenameUsage)
	flag.StringVar(&config.HostString, "a", config.HostString, "HTTP server endpoint")
	flag.StringVar(&config.GRPCHost, "g", config.GRPCHost, "GRPC server endpoint")
	flag.StringVar(&config.LogLevel, "l", config.LogLevel, "Log level")
	flag.IntVar(&storeInterval, "i", -1, "File store interval, 0 - synchrose")
	flag.StringVar(&config.FileStoragePath, "f", config.FileStoragePath, "File store path, empty - without store")
	flag.BoolVar(&config.Restore, "r", config.Restore, "Needs restore on start")
	flag.StringVar(&config.DSN, "d", config.DSN, "Database string")
	flag.StringVar(&config.SignKey, "k", config.SignKey, "SighHash Key")
	flag.StringVar(&config.CryptoKey, "crypto-key", config.CryptoKey, "Crypto Key")
	flag.StringVar(&config.TrustedSubnet, "t", config.TrustedSubnet, "Trusted subnet (CIDR)")
	flag.Parse()

	if storeInterval != -1 {
		config.StoreInterval = Duration{time.Duration(storeInterval) * time.Second}
	}

	// environment override
	err = env.Parse(&config)
	if err != nil {
		return nil, fmt.Errorf("error parsing env config: %w", err)
	}

	err = lookupEnvDuration("STORE_INTERVAL", &config.StoreInterval)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
