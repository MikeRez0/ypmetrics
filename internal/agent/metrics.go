package agent

import (
	"github.com/MikeRez0/ypmetrics/internal/storage"
)

type MetricValue struct {
	MetricType string
	Name       string
	Value      any
}

type MetricStore struct {
	Metrics map[string]MetricValue
}

func NewMetricStore() *MetricStore {
	var ms MetricStore
	ms.Metrics = make(map[string]MetricValue)
	return &ms
}

func (ms *MetricStore) PushGaugeMetric(name string, value storage.GaugeValue) {
	ms.Metrics[name] = MetricValue{
		MetricType: storage.GaugeType,
		Name:       name,
		Value:      value,
	}
}
func (ms *MetricStore) PushCounterMetric(name string, value storage.CounterValue) {
	newValue := storage.CounterValue(0)
	if val, ok := ms.Metrics[name]; ok {
		newValue = val.Value.(storage.CounterValue)
	}
	ms.Metrics[name] = MetricValue{
		MetricType: storage.CounterType,
		Name:       name,
		Value:      newValue + value,
	}
}
