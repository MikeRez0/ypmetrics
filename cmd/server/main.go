package main

import (
	// "net/http"
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

func run() error {

	hostString := flag.String("a", `localhost:8080`, "HTTP server endpoint")
	flag.Parse()

	if envHostString := os.Getenv("ADDRESS"); envHostString != "" {
		fmt.Printf("ADDRESS=%s", envHostString)
		*hostString = envHostString
	}

	var ms = storage.NewMemStorage()
	var h = handlers.NewMetricsHandler(ms)

	r := gin.Default()
	r.GET("/", h.MetricListView)
	r.POST("/update/:metricType/:metric/:value", h.UpdateMetricGin)
	r.GET("/value/:metricType/:metric", h.GetMetricGin)

	fmt.Printf("Starting server on %s", *hostString)
	return r.Run(*hostString)
	// mux := http.NewServeMux()
	// mux.Handle("/", h)

	// return http.ListenAndServe(`:8080`, mux)
}
