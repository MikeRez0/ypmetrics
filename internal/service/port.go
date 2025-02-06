package service

import (
	"context"

	"github.com/MikeRez0/ypmetrics/internal/model"
)

// Repository - Interface for metrics repository. Access/update metric value by name.
// Batch update multiple metrics.
//
//go:generate mockgen -source=./service.go -package mock -destination ./mock/service.go
type Repository interface {
	// List all metrics with values
	Metrics() []model.Metrics
	// Update gauge metric
	UpdateGauge(context context.Context, metric string, value model.GaugeValue) (model.GaugeValue, error)
	// Get gauge metric
	GetGauge(context context.Context, metric string) (model.GaugeValue, error)
	// Update counter metric
	UpdateCounter(context context.Context, metric string, value model.CounterValue) (model.CounterValue, error)
	// Get counter metric
	GetCounter(context context.Context, metric string) (model.CounterValue, error)
	// Update multiple metrics
	BatchUpdate(ctx context.Context, metrics []model.Metrics) error
	// Ping storage
	Ping() error
}

type IMetricService interface {
	GetMetric(ctx context.Context, metric *model.Metrics) error
	UpdateMetric(ctx context.Context, metric *model.Metrics) error
	BatchUpdateMetrics(ctx context.Context, metrics *[]model.Metrics) error
	Metrics() []model.Metrics
	Ping() error
}
