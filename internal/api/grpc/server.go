package gapi

import (
	"context"

	gapi "github.com/MikeRez0/ypmetrics/internal/api/grpc/proto"
	"github.com/MikeRez0/ypmetrics/internal/service"
	"go.uber.org/zap"
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/metrics.proto

type MetricService struct {
	gapi.UnimplementedMetricServiceServer
	repo service.Repository
	log  *zap.Logger
}

func (m *MetricService) GetMetric(ctx context.Context, in *gapi.RequestMetric) (*gapi.Metric, error) {
	return nil, nil
}
func (m *MetricService) UpdateMetric(ctx context.Context, in *gapi.Metric) (*gapi.Metric, error) {
	return nil, nil
}
func (m *MetricService) UpdateMetricBatch(ctx context.Context, in *gapi.RequestMetricList) (*gapi.Empty, error) {
	return nil, nil
}

func CreateServer(repo service.Repository, log *zap.Logger) (*MetricService, error) {
	return &MetricService{
		repo: repo,
		log:  log,
	}, nil
}
