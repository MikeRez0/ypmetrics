package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/MikeRez0/ypmetrics/internal/model"
)

type MemStorage struct {
	MetricsGauge   sync.Map
	MetricsCounter sync.Map
}

func NewMemStorage() *MemStorage {
	return &MemStorage{sync.Map{}, sync.Map{}}
}

func (ms *MemStorage) Metrics() (res []model.Metrics) {
	ms.MetricsCounter.Range(func(key, value any) bool {
		name, ok := key.(string)
		if !ok {
			return false
		}
		if val, ok := value.(model.CounterValue); ok {
			res = append(res, model.Metrics{
				ID:    name,
				MType: model.CounterType,
				Delta: (*int64)(&val),
			})
		}
		return true
	})

	ms.MetricsGauge.Range(func(key, value any) bool {
		name, ok := key.(string)
		if !ok {
			return false
		}
		if val, ok := value.(model.GaugeValue); ok {
			// val := value.(float64)
			res = append(res, model.Metrics{
				ID:    name,
				MType: model.GaugeType,
				Value: (*float64)(&val),
			})
		}

		return true
	})

	return res
}

func (ms *MemStorage) StoreMetric(ctx context.Context, metric model.Metrics) error {
	switch metric.MType {
	case model.CounterType:
		ms.MetricsCounter.Store(metric.ID, model.CounterValue(*metric.Delta))
	case model.GaugeType:
		ms.MetricsGauge.Store(metric.ID, model.GaugeValue(*metric.Value))
	}

	return nil
}

func (ms *MemStorage) UpdateGauge(ctx context.Context,
	metric string, value model.GaugeValue) (model.GaugeValue, error) {
	ms.MetricsGauge.Store(metric, value)

	v, _ := ms.MetricsGauge.Load(metric)

	return v.(model.GaugeValue), nil //nolint:forcetypeassert //this is why
}

func (ms *MemStorage) GetGauge(ctx context.Context, metric string) (model.GaugeValue, error) {
	if val, ok := ms.MetricsGauge.Load(metric); ok {
		return val.(model.GaugeValue), nil //nolint:forcetypeassert //this is why
	} else {
		return 0, fmt.Errorf("not found %s", metric)
	}
}

func (ms *MemStorage) UpdateCounter(ctx context.Context,
	metric string, value model.CounterValue) (model.CounterValue, error) {
	var m, ok = ms.MetricsCounter.Load(metric)
	if !ok {
		m = model.CounterValue(0)
	}
	val, _ := m.(model.CounterValue)
	ms.MetricsCounter.Store(metric, val+value)

	v, _ := ms.MetricsCounter.Load(metric)
	return v.(model.CounterValue), nil //nolint:forcetypeassert //this is why
}

func (ms *MemStorage) GetCounter(ctx context.Context, metric string) (model.CounterValue, error) {
	if val, ok := ms.MetricsCounter.Load(metric); ok {
		return val.(model.CounterValue), nil //nolint:forcetypeassert //this is why
	} else {
		return 0, fmt.Errorf("not found %s", metric)
	}
}

func (ms *MemStorage) Ping() error {
	return errors.New("Ping not supported")
}

func (ms *MemStorage) BatchUpdate(ctx context.Context, metrics []model.Metrics) error {
	for _, metric := range metrics {
		var err error
		switch metric.MType {
		case model.GaugeType:
			_, err = ms.UpdateGauge(ctx, metric.ID, model.GaugeValue(*metric.Value))
		case model.CounterType:
			_, err = ms.UpdateCounter(ctx, metric.ID, model.CounterValue(*metric.Delta))
		default:
			err = fmt.Errorf("unrecognized metric type %s", metric.MType)
		}

		if err != nil {
			return err
		}
	}
	return nil
}
