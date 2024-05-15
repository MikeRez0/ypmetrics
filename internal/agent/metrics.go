package agent

import (
	"github.com/MikeRez0/ypmetrics/internal/storage"
)

type MetricStore struct {
	MetricsGauge   map[string]storage.GaugeValue
	MetricsCounter map[string]storage.CounterValue
}

func NewMetricStore() *MetricStore {
	var ms MetricStore
	ms.MetricsGauge = make(map[string]storage.GaugeValue)
	ms.MetricsCounter = make(map[string]storage.CounterValue)
	return &ms
}

func (ms *MetricStore) PushGaugeMetric(name string, value storage.GaugeValue) {
	ms.MetricsGauge[name] = value
}
func (ms *MetricStore) PushCounterMetric(name string, value storage.CounterValue) {
	newValue := storage.CounterValue(0)
	if val, ok := ms.MetricsCounter[name]; ok {
		newValue = val
	}
	ms.MetricsCounter[name] = newValue + value
}
