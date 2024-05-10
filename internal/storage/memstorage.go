package storage

import (
	"fmt"
)

type MemStorage struct {
	MetricsGauge   map[string]GaugeValue
	MetricsCounter map[string]CounterValue
}

func NewMemStorage() *MemStorage {
	mg := make(map[string]GaugeValue)
	mc := make(map[string]CounterValue)
	return &MemStorage{mg, mc}
}

func (ms *MemStorage) Metrics() (res []struct{ Name, Value string }) {
	// var res []struct{ Name, Value string }
	for name, value := range ms.MetricsCounter {
		res = append(res, struct {
			Name  string
			Value string
		}{name, fmt.Sprint(value)})
	}
	for name, value := range ms.MetricsGauge {
		res = append(res, struct {
			Name  string
			Value string
		}{name, fmt.Sprint(value)})
	}
	return res
}

func (ms *MemStorage) UpdateGauge(metric string, value GaugeValue) {
	ms.MetricsGauge[metric] = value
}

func (ms *MemStorage) GetGauge(metric string) (GaugeValue, error) {
	if val, ok := ms.MetricsGauge[metric]; ok {
		return val, nil
	} else {
		return 0, fmt.Errorf("not found %s", metric)
	}
}

func (ms *MemStorage) UpdateCounter(metric string, value CounterValue) {
	var m, ok = ms.MetricsCounter[metric]
	if !ok {
		m = 0
	}
	ms.MetricsCounter[metric] = m + value
}

func (ms *MemStorage) GetCounter(metric string) (CounterValue, error) {
	if val, ok := ms.MetricsCounter[metric]; ok {
		return val, nil
	} else {
		return 0, fmt.Errorf("not found %s", metric)
	}
}
