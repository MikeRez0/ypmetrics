package handlers

import (
	"log"
	"net/http"

	"github.com/MikeRez0/ypmetrics/internal/storage"
)

type Repository interface {
	Metrics() []struct{ Name, Value string }
	UpdateGauge(metric string, value storage.GaugeValue)
	GetGauge(metric string) (storage.GaugeValue, error)
	UpdateCounter(metric string, value storage.CounterValue)
	GetCounter(metric string) (storage.CounterValue, error)
}

type MetricsHandler struct {
	Store Repository
}

func NewMetricsHandler(s Repository) *MetricsHandler {
	return &MetricsHandler{Store: s}
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Status", http.StatusInternalServerError, r.URL)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}
func BadRequestErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("400 Bad request"))
}
func NotFoundErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not found"))
}
func MethodNotAllowedErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("405 Method not allowed"))
}
