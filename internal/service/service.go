package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/MikeRez0/ypmetrics/internal/model"
	"go.uber.org/zap"
)

type MetricService struct {
	Store Repository
	log   *zap.Logger
}

func NewMetricService(repo Repository, log *zap.Logger) (*MetricService, error) {
	return &MetricService{
		Store: repo,
		log:   log,
	}, nil
}

func (s *MetricService) GetMetric(c context.Context, metric *model.Metrics) error {
	if metric.ID == "" {
		return model.ErrDataNotFound
	}

	switch metric.MType {
	case model.GaugeType:
		value, err := s.Store.GetGauge(c, metric.ID)
		if err != nil {
			return model.ErrDataNotFound
		}
		metric.Value = (*float64)(&value)
	case model.CounterType:
		value, err := s.Store.GetCounter(c, metric.ID)
		if err != nil {
			return model.ErrDataNotFound
		}
		metric.Delta = (*int64)(&value)
	default:
		return model.ErrBadRequest
	}
	return nil
}
func (s *MetricService) UpdateMetric(c context.Context, metric *model.Metrics) error {
	if metric.ID == "" {
		return model.ErrDataNotFound
	}

	switch metric.MType {
	case model.GaugeType:
		v, err := s.Store.UpdateGauge(c, metric.ID, model.GaugeValue(*metric.Value))
		if err != nil {
			return model.ErrInternal
		}
		var newVal = float64(v)
		metric.Value = &newVal
	case model.CounterType:
		v, err := s.Store.UpdateCounter(c, metric.ID, model.CounterValue(*metric.Delta))
		if err != nil {
			return model.ErrInternal
		}
		var newVal = int64(v)
		metric.Delta = &newVal
	default:
		return model.ErrBadRequest
	}
	return nil
}
func (s *MetricService) BatchUpdateMetrics(c context.Context, metrics *[]model.Metrics) error {
	err := s.Store.BatchUpdate(c, *metrics)
	if err != nil {
		if errors.As(err, &model.BadValueError{}) {
			return model.ErrBadRequest
		}
		return model.ErrInternal
	}
	return nil
}

func (s *MetricService) Metrics() []model.Metrics {
	return s.Store.Metrics()
}

func (s *MetricService) Ping() error {
	err := s.Store.Ping()
	if err != nil {
		return fmt.Errorf("error on ping repository:%w", err)
	}
	return nil
}
