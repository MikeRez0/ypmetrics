package handlers

import (
	_ "embed"
	"fmt"
	"html/template"

	"github.com/MikeRez0/ypmetrics/internal/storage"
)

type Repository interface {
	Metrics() []struct{ Name, Value string }
	UpdateGauge(metric string, value storage.GaugeValue) storage.GaugeValue
	GetGauge(metric string) (storage.GaugeValue, error)
	UpdateCounter(metric string, value storage.CounterValue) storage.CounterValue
	GetCounter(metric string) (storage.CounterValue, error)
}

type MetricsHandler struct {
	Store    Repository
	Template *template.Template
}

//go:embed "templates/metrics.html"
var templateContent string

func NewMetricsHandler(s Repository) (*MetricsHandler, error) {
	tmpl := template.New("metrics")
	var err error
	tmpl, err = tmpl.Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("error while parsing template: %w", err)
	}

	return &MetricsHandler{Store: s, Template: tmpl}, nil
}
