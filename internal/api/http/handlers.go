// Package http - API request handlers for http-server.
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
package http

import (
	_ "embed"
	"fmt"
	"html/template"

	"go.uber.org/zap"

	"github.com/MikeRez0/ypmetrics/internal/service"
	"github.com/MikeRez0/ypmetrics/internal/utils/signer"
)

type MetricsHandler struct {
	// Store     service.Repository
	service   service.IMetricService
	Template  *template.Template
	Log       *zap.Logger
	Signer    *signer.Signer
	Decrypter *signer.Decrypter
}

//go:embed "templates/metrics.html"
var templateContent string

// NewMetricsHandler - create new metrics handler
//
// - Repository - metric storage implementation.
//
// - zap.Logger - logger.
func NewMetricsHandler(s service.IMetricService, log *zap.Logger) (*MetricsHandler, error) {
	tmpl := template.New("metrics")
	var err error
	tmpl, err = tmpl.Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("error while parsing template: %w", err)
	}

	return &MetricsHandler{service: s, Template: tmpl, Log: log}, nil
}
