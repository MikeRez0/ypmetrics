package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

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

func setupRouter(h *handlers.MetricsHandler, mylog *zap.Logger) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logger.GinLogger(mylog))
	r.HandleMethodNotAllowed = true

	r.GET("/", gzip.Gzip(gzip.DefaultCompression), h.MetricListView)
	r.POST("/update/:metricType/:metric/:value", h.UpdateMetricPlain)
	r.GET("/value/:metricType/:metric", h.GetMetricPlain)

	jsonGroup := r.Group("/")
	jsonGroup.Use(handlers.GinCompress(logger.LoggerWithComponent(mylog, "compress")))
	jsonGroup.POST("/update/", h.UpdateMetricJSON)
	jsonGroup.POST("/value/", h.GetMetricJSON)
	jsonGroup.POST("/updates/", h.BatchUpdateMetricsJSON)

	r.GET("/ping", h.PingDB)

	return r
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

	var repo handlers.Repository

	switch {
	case conf.DSN != "":
		repo, err = storage.NewDBStorage(
			conf.DSN,
			conf.StoreInterval,
			conf.Restore,
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
	r := setupRouter(h, logger.LoggerWithComponent(mylog, "handlers"))

	if conf.SignKey != "" {
		h.Signer = signer.NewSigner(conf.SignKey)
	}

	err = r.Run(conf.HostString)
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error while run server: %w", err)
	}
	return nil
}
