package agent

import (
	"github.com/MikeRez0/ypmetrics/internal/model"
)

type MetricStore struct {
	MetricsGauge   map[string]model.GaugeValue
	MetricsCounter map[string]model.CounterValue
}

func NewMetricStore() *MetricStore {
	var ms MetricStore
	ms.MetricsGauge = make(map[string]model.GaugeValue)
	ms.MetricsCounter = make(map[string]model.CounterValue)
	return &ms
}

func (ms *MetricStore) PushGaugeMetric(name string, value model.GaugeValue) {
	ms.MetricsGauge[name] = value
}
func (ms *MetricStore) PushCounterMetric(name string, value model.CounterValue) {
	newValue := model.CounterValue(0)
	if val, ok := ms.MetricsCounter[name]; ok {
		newValue = val
	}
	ms.MetricsCounter[name] = newValue + value
}
