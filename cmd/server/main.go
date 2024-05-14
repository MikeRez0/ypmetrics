package main

import (
	"flag"
	"fmt"

	"os"

	"github.com/gin-gonic/gin"

	"github.com/MikeRez0/ypmetrics/internal/handlers"
	"github.com/MikeRez0/ypmetrics/internal/storage"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		panic(err)
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
	hostString := flag.String("a", `localhost:8080`, "HTTP server endpoint")
	flag.Parse()

	if envHostString := os.Getenv("ADDRESS"); envHostString != "" {
		*hostString = envHostString
	}

	var ms = storage.NewMemStorage()
	var h = handlers.NewMetricsHandler(ms)
	r := setupRouter(h)

	return r.Run(*hostString)
}
