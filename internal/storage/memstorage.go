package storage

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/MikeRez0/ypmetrics/internal/model"
)

type MemStorage struct {
	MetricsGauge   sync.Map //map[string]model.GaugeValue
	MetricsCounter sync.Map //map[string]model.CounterValue
}

func NewMemStorage() *MemStorage {
	// mg := make(map[string]model.GaugeValue)
	// mc := make(map[string]model.CounterValue)
	return &MemStorage{sync.Map{}, sync.Map{}}
}

func (ms *MemStorage) Metrics() (res []model.Metrics) {
	ms.MetricsCounter.Range(func(name, value any) bool {
		if val, ok := value.(model.CounterValue); ok {
			// val := value.(model.CounterValue).(int64)
			res = append(res, model.Metrics{
				ID:    name.(string),
				MType: model.CounterType,
				Delta: (*int64)(&val),
			})
		}
		return true
	})
	// for name, value := range ms.MetricsCounter {
	// 	res = append(res, model.Metrics{
	// 		ID:    name,
	// 		MType: model.CounterType,
	// 		Delta: (*int64)(&value),
	// 	})
	// }

	ms.MetricsGauge.Range(func(name, value any) bool {
		if val, ok := value.(model.GaugeValue); ok {
			// val := value.(float64)
			res = append(res, model.Metrics{
				ID:    name.(string),
				MType: model.GaugeType,
				Value: (*float64)(&val),
			})
		}

		return true
	})

	// for name, value := range ms.MetricsGauge {
	// 	res = append(res, model.Metrics{
	// 		ID:    name,
	// 		MType: model.GaugeType,
	// 		Value: (*float64)(&value),
	// 	})
	// }

	return res
}

func (ms *MemStorage) StoreMetric(ctx context.Context, metric model.Metrics) error {
	switch metric.MType {
	case model.CounterType:
		ms.MetricsCounter.Store(metric.ID, model.CounterValue(*metric.Delta))
		// ms.MetricsCounter[metric.ID] = model.CounterValue(*metric.Delta)
	case model.GaugeType:
		ms.MetricsGauge.Store(metric.ID, model.GaugeValue(*metric.Value))
		// ms.MetricsGauge[metric.ID] = model.GaugeValue(*metric.Value)
	}

	return nil
}

func (ms *MemStorage) MetricStrings() (res []struct{ Name, Value string }) {
	ms.MetricsCounter.Range(func(name, value any) bool {
		if val, ok := value.(model.CounterValue); ok {
			res = append(res, struct {
				Name  string
				Value string
			}{name.(string), fmt.Sprint(val)})
		}
		return true
	})

	// for name, value := range ms.MetricsCounter {
	// 	res = append(res, struct {
	// 		Name  string
	// 		Value string
	// 	}{name, strconv.Itoa(int(value))})
	// }

	ms.MetricsGauge.Range(func(name, value any) bool {
		if val, ok := value.(model.GaugeValue); ok {
			res = append(res, struct {
				Name  string
				Value string
			}{name.(string), strconv.FormatFloat(float64(val), 'f', 5, 64)})
		}
		return true
	})

	// for name, value := range ms.MetricsGauge {
	// 	res = append(res, struct {
	// 		Name  string
	// 		Value string
	// 	}{name, strconv.FormatFloat(float64(value), 'f', 5, 64)})
	// }
	return res
}

func (ms *MemStorage) UpdateGauge(ctx context.Context,
	metric string, value model.GaugeValue) (model.GaugeValue, error) {
	ms.MetricsGauge.Store(metric, value)
	// ms.MetricsGauge[metric] = value

	v, _ := ms.MetricsGauge.Load(metric)

	return v.(model.GaugeValue), nil
}

func (ms *MemStorage) GetGauge(ctx context.Context, metric string) (model.GaugeValue, error) {
	if val, ok := ms.MetricsGauge.Load(metric); ok {
		return val.(model.GaugeValue), nil
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
	val := m.(model.CounterValue)
	ms.MetricsCounter.Store(metric, val+value)
	// ms.MetricsCounter[metric] = m + value

	v, _ := ms.MetricsCounter.Load(metric)
	return v.(model.CounterValue), nil
}

func (ms *MemStorage) GetCounter(ctx context.Context, metric string) (model.CounterValue, error) {
	if val, ok := ms.MetricsCounter.Load(metric); ok {
		return val.(model.CounterValue), nil
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
