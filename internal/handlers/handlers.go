package handlers

import (
	_ "embed"
	"fmt"
	"html/template"

	"github.com/MikeRez0/ypmetrics/internal/model"
	"go.uber.org/zap"
)

type Repository interface {
	Metrics() []model.Metrics
	MetricStrings() []struct{ Name, Value string }
	StoreMetric(metric model.Metrics) error
	UpdateGauge(metric string, value model.GaugeValue) (model.GaugeValue, error)
	GetGauge(metric string) (model.GaugeValue, error)
	UpdateCounter(metric string, value model.CounterValue) (model.CounterValue, error)
	GetCounter(metric string) (model.CounterValue, error)
}

type MetricsHandler struct {
	Store    Repository
	Template *template.Template
	Log      *zap.Logger
}

//go:embed "templates/metrics.html"
var templateContent string

func NewMetricsHandler(s Repository, log *zap.Logger) (*MetricsHandler, error) {
	tmpl := template.New("metrics")
	var err error
	tmpl, err = tmpl.Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("error while parsing template: %w", err)
	}

	return &MetricsHandler{Store: s, Template: tmpl, Log: log}, nil
}
