package agent

import (
	"sync"

	"github.com/MikeRez0/ypmetrics/internal/model"
)

// MetricStore - metric store in agent app.
type MetricStore struct {
	metricsGauge   map[string]model.GaugeValue
	metricsCounter map[string]model.CounterValue
	l              sync.RWMutex
}

// NewMetricStore - create metric store for agent.
func NewMetricStore() *MetricStore {
	var ms MetricStore
	ms.metricsGauge = make(map[string]model.GaugeValue)
	ms.metricsCounter = make(map[string]model.CounterValue)
	return &ms
}

// PushGaugeMetric - save gauge metric.
func (ms *MetricStore) PushGaugeMetric(name string, value model.GaugeValue) {
	ms.l.Lock()
	ms.metricsGauge[name] = value
	ms.l.Unlock()
}

// PushCounterMetric - save counter metric.
func (ms *MetricStore) PushCounterMetric(name string, value model.CounterValue) {
	ms.l.Lock()
	newValue := model.CounterValue(0)
	if val, ok := ms.metricsCounter[name]; ok {
		newValue = val
	}
	ms.metricsCounter[name] = newValue + value
	ms.l.Unlock()
}

// GetGaugeMetrics - get gauge metric.
func (ms *MetricStore) GetGaugeMetrics() map[string]model.GaugeValue {
	ms.l.RLock()
	res := make(map[string]model.GaugeValue, len(ms.metricsGauge))
	for k, v := range ms.metricsGauge {
		res[k] = v
	}
	ms.l.RUnlock()

	return res
}

// GetCounterMetrics - get counter metric.
func (ms *MetricStore) GetCounterMetrics() map[string]model.CounterValue {
	ms.l.RLock()
	res := make(map[string]model.CounterValue, len(ms.metricsCounter))
	for k, v := range ms.metricsCounter {
		res[k] = v
	}
	ms.l.RUnlock()

	return res
}

// Clear - clear metrics in agent store.
func (ms *MetricStore) Clear() {
	ms.l.Lock()

	clear(ms.metricsCounter)
	clear(ms.metricsGauge)

	ms.l.Unlock()
}
