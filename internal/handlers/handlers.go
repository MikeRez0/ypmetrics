package handlers

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"

	"github.com/MikeRez0/ypmetrics/internal/model"
	"github.com/MikeRez0/ypmetrics/internal/utils/signer"
	"go.uber.org/zap"
)

type Repository interface {
	Metrics() []model.Metrics
	MetricStrings() []struct{ Name, Value string }
	StoreMetric(context context.Context, metric model.Metrics) error
	UpdateGauge(context context.Context, metric string, value model.GaugeValue) (model.GaugeValue, error)
	GetGauge(context context.Context, metric string) (model.GaugeValue, error)
	UpdateCounter(context context.Context, metric string, value model.CounterValue) (model.CounterValue, error)
	GetCounter(context context.Context, metric string) (model.CounterValue, error)
	BatchUpdate(ctx context.Context, metrics []model.Metrics) error
	Ping() error
}

type MetricsHandler struct {
	Store    Repository
	Template *template.Template
	Log      *zap.Logger
	Signer   *signer.Signer
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
