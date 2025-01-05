// API request handlers for http-server.
// Example:
//
//	mylog := zap.NewProduction()
//	repo = storage.NewMemStorage()
//
//	h, err := handlers.NewMetricsHandler(repo, logger.LoggerWithComponent(mylog, "handlers"))
//	if err != nil {
//		return fmt.Errorf("error creating handler: %w", err)
//	}
//	r := handlers.SetupRouter(h, logger.LoggerWithComponent(mylog, "handlers"))
//	err = r.Run("localhost:8080")
package handlers

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"

	"go.uber.org/zap"

	"github.com/MikeRez0/ypmetrics/internal/model"
	"github.com/MikeRez0/ypmetrics/internal/utils/signer"
)

// Repository - Interface for metrics repository. Access/update metric value by name.
// Batch update multiple metrics.

//go:generate mockgen -source=./handlers.go -package mock -destination ./mock/handlers.go
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

type MetricsHandler struct {
	Store    Repository
	Template *template.Template
	Log      *zap.Logger
	Signer   *signer.Signer
}

//go:embed "templates/metrics.html"
var templateContent string

// NewMetricsHandler - create new metrics handler
//
// - Repository - metric storage implementation.
//
// - zap.Logger - logger.
func NewMetricsHandler(s Repository, log *zap.Logger) (*MetricsHandler, error) {
	tmpl := template.New("metrics")
	var err error
	tmpl, err = tmpl.Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("error while parsing template: %w", err)
	}

	return &MetricsHandler{Store: s, Template: tmpl, Log: log}, nil
}
