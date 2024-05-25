package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/MikeRez0/ypmetrics/internal/config"
	"github.com/MikeRez0/ypmetrics/internal/handlers"
	"github.com/MikeRez0/ypmetrics/internal/logger"
	"github.com/MikeRez0/ypmetrics/internal/storage"
)

func main() {
	if err := run(); err != nil {
		logger.Log.Fatal("Fatal error", zap.Error(err))
	}
}

func setupRouter(h *handlers.MetricsHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logger.GinLogger())
	r.HandleMethodNotAllowed = true

	r.GET("/", h.MetricListView)
	r.POST("/update/:metricType/:metric/:value", h.UpdateMetricGin)
	r.GET("/value/:metricType/:metric", h.GetMetricGin)

	return r
}

func run() error {
	conf, err := config.NewConfigServer()
	if err != nil {
		return fmt.Errorf("error while config load: %w", err)
	}

	err = logger.Initialize(conf.LogLevel)
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	var ms = storage.NewMemStorage()
	h, err := handlers.NewMetricsHandler(ms)
	if err != nil {
		return fmt.Errorf("error creating handler: %w", err)
	}
	r := setupRouter(h)

	err = r.Run(conf.HostString)
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error while run server: %w", err)
	}
	return nil
}
