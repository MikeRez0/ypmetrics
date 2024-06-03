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

	// не получилось использовать свой мидлвар, потому что в ответ
	// встраивалось application/x-gzip, игнорируя "мои" заголовки
	// обсудить на 1-1
	r.GET("/", gzip.Gzip(gzip.DefaultCompression), h.MetricListView)
	r.POST("/update/:metricType/:metric/:value", h.UpdateMetricPlain)
	r.GET("/value/:metricType/:metric", h.GetMetricPlain)
	r.POST("/update/", handlers.GinCompress(mylog), h.UpdateMetricJSON)
	r.POST("/value/", handlers.GinCompress(mylog), h.GetMetricJSON)

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

	if conf.FileStoragePath != "" {
		repo, err = storage.NewFileStorage(
			conf.FileStoragePath,
			conf.StoreInterval,
			conf.Restore,
			mylog)
		if err != nil {
			return fmt.Errorf("error creating file repo: %w", err)
		}
	} else {
		repo = storage.NewMemStorage()
	}
	h, err := handlers.NewMetricsHandler(repo, mylog)
	if err != nil {
		return fmt.Errorf("error creating handler: %w", err)
	}
	r := setupRouter(h, mylog)

	err = r.Run(conf.HostString)
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error while run server: %w", err)
	}
	return nil
}
