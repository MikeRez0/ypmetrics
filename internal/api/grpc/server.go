package grpc

import (
	"context"
	"errors"

	pb "github.com/MikeRez0/ypmetrics/internal/api/grpc/proto"
	"github.com/MikeRez0/ypmetrics/internal/model"
	"github.com/MikeRez0/ypmetrics/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/metrics.proto

type MetricService struct {
	pb.UnimplementedMetricServiceServer
	service service.IMetricService
	log     *zap.Logger
}

func (m *MetricService) GetMetric(ctx context.Context, in *pb.RequestMetric) (*pb.Metric, error) {
	metric := model.Metrics{
		ID:    in.GetID(),
		MType: model.MetricType(in.GetType()),
	}

	err := m.service.GetMetric(ctx, &metric)
	switch {
	case errors.Is(err, model.ErrDataNotFound):
		return nil, status.Errorf(codes.NotFound, "Metric not found")
	case errors.Is(err, model.ErrBadRequest):
		return nil, status.Errorf(codes.InvalidArgument, "Metric type not found")
	}

	result := pb.Metric{
		ID:    metric.ID,
		Type:  string(metric.MType),
		Delta: *metric.Delta,
		Value: *metric.Value,
	}

	return &result, nil
}
func (m *MetricService) UpdateMetric(ctx context.Context, in *pb.Metric) (*pb.Metric, error) {
	value := in.GetValue()
	delta := in.GetDelta()
	metric := model.Metrics{
		ID:    in.GetID(),
		MType: model.MetricType(in.GetType()),
		Value: &value,
		Delta: &delta,
	}

	err := m.service.UpdateMetric(ctx, &metric)
	switch {
	case errors.Is(err, model.ErrDataNotFound):
		return nil, status.Errorf(codes.NotFound, "Metric not found")
	case errors.Is(err, model.ErrBadRequest):
		return nil, status.Errorf(codes.InvalidArgument, "Metric type not found")
	}

	result := pb.Metric{
		ID:    metric.ID,
		Type:  string(metric.MType),
		Delta: *metric.Delta,
		Value: *metric.Value,
	}

	return &result, nil
}
func (m *MetricService) UpdateMetricBatch(ctx context.Context, in *pb.RequestMetricList) (*pb.Empty, error) {
	metricList := make([]model.Metrics, 0, len(in.GetMetrics()))
	for _, m := range in.GetMetrics() {
		metricList = append(metricList, readMetric(m))
	}

	err := m.service.BatchUpdateMetrics(ctx, &metricList)
	switch {
	case errors.Is(err, model.ErrDataNotFound):
		return nil, status.Errorf(codes.NotFound, "Metric not found")
	case errors.Is(err, model.ErrBadRequest):
		return nil, status.Errorf(codes.InvalidArgument, "Metric type not found")
	}

	return &pb.Empty{}, nil
}

func readMetric(m *pb.Metric) model.Metrics {
	value := m.GetValue()
	delta := m.GetDelta()
	return model.Metrics{
		ID:    m.GetID(),
		MType: model.MetricType(m.GetType()),
		Value: &value,
		Delta: &delta,
	}
}

func CreateServer(serv service.IMetricService, log *zap.Logger) (*grpc.Server, error) {
	gs := grpc.NewServer()
	pb.RegisterMetricServiceServer(gs, &MetricService{
		service: serv,
		log:     log,
	})

	return gs, nil
}
