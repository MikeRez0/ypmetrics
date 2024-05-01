package main

import (
	"net/http"

	"github.com/MikeRez0/ypmetrics/internal/handlers"
	"github.com/MikeRez0/ypmetrics/internal/storage"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {

	var ms = storage.NewMemStorage()
	var h = handlers.NewMetricsHandler(ms)

	mux := http.NewServeMux()
	mux.Handle("/", h)

	return http.ListenAndServe(`:8080`, mux)
}
