package handlers

import (
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
