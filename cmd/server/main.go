package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/MikeRez0/ypmetrics/internal/config"
	"github.com/MikeRez0/ypmetrics/internal/handlers"
	"github.com/MikeRez0/ypmetrics/internal/logger"
	"github.com/MikeRez0/ypmetrics/internal/storage"
	"github.com/MikeRez0/ypmetrics/internal/utils/signer"
)

func main() {
	if err := run(); err != nil {
		// no custom logger at this line
		log.Fatalf("Fatal error: %v", err)
	}
}

func run() error {
	conf, err := config.NewConfigServer()
	if err != nil {
		return fmt.Errorf("error while config load: %w", err)
	}

	mylog, err := logger.Initialize(conf.LogLevel)
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	mylog.Info(fmt.Sprintf("start server with config: %v", conf))

	var repo handlers.Repository

	switch {
	case conf.DSN != "":
		repo, err = storage.NewDBStorage(
			conf.DSN,
			logger.LoggerWithComponent(mylog, "dbstorage"))
		if err != nil {
			return fmt.Errorf("error creating db repo: %w", err)
		}
	case conf.FileStoragePath != "":
		repo, err = storage.NewFileStorage(
			conf.FileStoragePath,
			conf.StoreInterval,
			conf.Restore,
			logger.LoggerWithComponent(mylog, "filestorage"))
		if err != nil {
			return fmt.Errorf("error creating file repo: %w", err)
		}
	default:
		repo = storage.NewMemStorage()
	}

	h, err := handlers.NewMetricsHandler(repo, logger.LoggerWithComponent(mylog, "handlers"))
	if err != nil {
		return fmt.Errorf("error creating handler: %w", err)
	}
	r := handlers.SetupRouter(h, logger.LoggerWithComponent(mylog, "handlers"))

	if conf.SignKey != "" {
		h.Signer = signer.NewSigner(conf.SignKey)
	}

	err = r.Run(conf.HostString)
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error while run server: %w", err)
	}
	return nil
}
