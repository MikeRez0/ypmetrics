package storage

import (
	"fmt"
	"strconv"
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

func (ms *MemStorage) UpdateGauge(metric string, value GaugeValue) GaugeValue {
	ms.MetricsGauge[metric] = value
	return ms.MetricsGauge[metric]
}

func (ms *MemStorage) GetGauge(metric string) (GaugeValue, error) {
	if val, ok := ms.MetricsGauge[metric]; ok {
		return val, nil
	} else {
		return 0, fmt.Errorf("not found %s", metric)
	}
}

func (ms *MemStorage) UpdateCounter(metric string, value CounterValue) CounterValue {
	var m, ok = ms.MetricsCounter[metric]
	if !ok {
		m = 0
	}
	ms.MetricsCounter[metric] = m + value
	return ms.MetricsCounter[metric]
}

func (ms *MemStorage) GetCounter(metric string) (CounterValue, error) {
	if val, ok := ms.MetricsCounter[metric]; ok {
		return val, nil
	} else {
		return 0, fmt.Errorf("not found %s", metric)
	}
}
