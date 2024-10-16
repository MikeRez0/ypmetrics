package storage

import (
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

func (ms *MemStorage) StoreMetric(metric model.Metrics) error {
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

func (ms *MemStorage) UpdateGauge(metric string, value model.GaugeValue) (model.GaugeValue, error) {
	ms.MetricsGauge[metric] = value
	return ms.MetricsGauge[metric], nil
}

func (ms *MemStorage) GetGauge(metric string) (model.GaugeValue, error) {
	if val, ok := ms.MetricsGauge[metric]; ok {
		return val, nil
	} else {
		return 0, fmt.Errorf("not found %s", metric)
	}
}

func (ms *MemStorage) UpdateCounter(metric string, value model.CounterValue) (model.CounterValue, error) {
	var m, ok = ms.MetricsCounter[metric]
	if !ok {
		m = 0
	}
	ms.MetricsCounter[metric] = m + value
	return ms.MetricsCounter[metric], nil
}

func (ms *MemStorage) GetCounter(metric string) (model.CounterValue, error) {
	if val, ok := ms.MetricsCounter[metric]; ok {
		return val, nil
	} else {
		return 0, fmt.Errorf("not found %s", metric)
	}
}
