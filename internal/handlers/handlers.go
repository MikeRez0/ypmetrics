package handlers

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/MikeRez0/ypmetrics/internal/storage"
)

type Repository interface {
	UpdateGauge(metric string, value storage.GaugeValue)
	GetGauge(metric string) (storage.GaugeValue, error)
	UpdateCounter(metric string, value storage.CounterValue)
	GetCounter(metric string) (storage.CounterValue, error)
}

type MetricsHandler struct {
	store Repository
}

func NewMetricsHandler(s Repository) *MetricsHandler {
	return &MetricsHandler{store: s}
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Status", http.StatusInternalServerError, r.URL)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}
func BadRequestErrorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Status", http.StatusBadRequest, r.URL)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("400 Bad request"))
}
func NotFoundErrorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Status", http.StatusNotFound, r.URL)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not found"))
}
func MethodNotAllowedErrorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Status", http.StatusMethodNotAllowed, r.URL)
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("405 Method not allowed"))
}

var (
	UpdateMetricRe = regexp.MustCompile(`^/update/(\w+)/(\w*)/?(\d*(?:\.\d+)?)$`) //`^/update/(\w+)/(\w+)/(\d+(?:\.\d+)?)$`)
	GetMetricRe    = regexp.MustCompile(`^/value/(\w+)/(\w+)$`)
)

func (mh *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if UpdateMetricRe.MatchString(r.URL.Path) {
			mh.UpdateMetric(w, r)
		} else {
			BadRequestErrorHandler(w, r)
		}
	} else if r.Method == http.MethodGet {
		if GetMetricRe.MatchString(r.URL.Path) {
			mh.GetMetric(w, r)
		} else {
			BadRequestErrorHandler(w, r)
		}
	} else {
		MethodNotAllowedErrorHandler(w, r)
		return
	}
}
