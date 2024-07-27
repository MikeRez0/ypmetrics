package storage

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/MikeRez0/ypmetrics/internal/model"
)

type MemStorage struct {
	MetricsGauge   map[string]model.GaugeValue
	MetricsCounter map[string]model.CounterValue
}

func NewMemStorage() *MemStorage {
	mg := make(map[string]model.GaugeValue)
	mc := make(map[string]model.CounterValue)
	return &MemStorage{mg, mc}
}

func (ms *MemStorage) Metrics() (res []model.Metrics) {
	for name, value := range ms.MetricsCounter {
		res = append(res, model.Metrics{
			ID:    name,
			MType: model.CounterType,
			Delta: (*int64)(&value),
		})
	}

	for name, value := range ms.MetricsGauge {
		res = append(res, model.Metrics{
			ID:    name,
			MType: model.GaugeType,
			Value: (*float64)(&value),
		})
	}

	return res
}

func (ms *MemStorage) StoreMetric(ctx context.Context, metric model.Metrics) error {
	switch metric.MType {
	case model.CounterType:
		ms.MetricsCounter[metric.ID] = model.CounterValue(*metric.Delta)
	case model.GaugeType:
		ms.MetricsGauge[metric.ID] = model.GaugeValue(*metric.Value)
	}

	return nil
}

func (ms *MemStorage) MetricStrings() (res []struct{ Name, Value string }) {
	for name, value := range ms.MetricsCounter {
		res = append(res, struct {
			Name  string
			Value string
		}{name, strconv.Itoa(int(value))})
	}
	for name, value := range ms.MetricsGauge {
		res = append(res, struct {
			Name  string
			Value string
		}{name, strconv.FormatFloat(float64(value), 'f', 5, 64)})
	}
	return res
}

func (ms *MemStorage) UpdateGauge(ctx context.Context,
	metric string, value model.GaugeValue) (model.GaugeValue, error) {
	ms.MetricsGauge[metric] = value
	return ms.MetricsGauge[metric], nil
}

func (ms *MemStorage) GetGauge(ctx context.Context, metric string) (model.GaugeValue, error) {
	if val, ok := ms.MetricsGauge[metric]; ok {
		return val, nil
	} else {
		return 0, fmt.Errorf("not found %s", metric)
	}
}

func (ms *MemStorage) UpdateCounter(ctx context.Context,
	metric string, value model.CounterValue) (model.CounterValue, error) {
	var m, ok = ms.MetricsCounter[metric]
	if !ok {
		m = 0
	}
	ms.MetricsCounter[metric] = m + value
	return ms.MetricsCounter[metric], nil
}

func (ms *MemStorage) GetCounter(ctx context.Context, metric string) (model.CounterValue, error) {
	if val, ok := ms.MetricsCounter[metric]; ok {
		return val, nil
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
