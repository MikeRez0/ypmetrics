package storage

type MemStorage struct {
	MetricsGauge   map[string]GaugeValue
	MetricsCounter map[string]CounterValue
}

func NewMemStorage() *MemStorage {
	mg := make(map[string]GaugeValue)
	mc := make(map[string]CounterValue)
	return &MemStorage{mg, mc}
}

func (ms *MemStorage) UpdateGauge(metric string, value GaugeValue) {
	ms.MetricsGauge[metric] = value
}

func (ms *MemStorage) GetGauge(metric string) (GaugeValue, error) {
	return ms.MetricsGauge[metric], nil
}

func (ms *MemStorage) UpdateCounter(metric string, value CounterValue) {
	var m, ok = ms.MetricsCounter[metric]
	if !ok {
		m = 0
	}
	ms.MetricsCounter[metric] = m + value
}

func (ms *MemStorage) GetCounter(metric string) (CounterValue, error) {
	return ms.MetricsCounter[metric], nil
}
