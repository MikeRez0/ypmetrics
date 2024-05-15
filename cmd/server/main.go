package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MikeRez0/ypmetrics/internal/config"
	"github.com/MikeRez0/ypmetrics/internal/handlers"
	"github.com/MikeRez0/ypmetrics/internal/storage"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func setupRouter(h *handlers.MetricsHandler) *gin.Engine {
	r := gin.Default()
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

	var ms = storage.NewMemStorage()
	var h = handlers.NewMetricsHandler(ms)
	r := setupRouter(h)

	err = r.Run(conf.HostString)
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error while run server: %w", err)
	}
	return nil
}
