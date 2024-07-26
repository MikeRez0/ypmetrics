package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

const GaugeType = "gauge"
const CounterType = "counter"

const (
	CounterTypeDB = iota + 1
	GaugeTypeDB
)

type MetricType string

type GaugeValue float64
type CounterValue int64

type Metrics struct { //nolint:govet //external rule from course author
	ID    string     `json:"id"`              // имя метрики
	MType MetricType `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64     `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64   `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (mt MetricType) Value() (driver.Value, error) {
	switch mt {
	case CounterType:
		return CounterTypeDB, nil
	case GaugeType:
		return GaugeTypeDB, nil
	default:
		return nil, fmt.Errorf(`unexpected value %s`, mt)
	}
}

func (mt *MetricType) Scan(value any) error {
	if value != nil {
		sv, err := driver.Int32.ConvertValue(value)
		if err != nil {
			return fmt.Errorf("cannot scan value. %w", err)
		}

		v, ok := sv.(int64)
		if !ok {
			return errors.New("cannot scan value. cannot convert value to MetricType")
		}

		switch v {
		case CounterTypeDB:
			*mt = CounterType
			return nil
		case GaugeTypeDB:
			*mt = GaugeType
			return nil
		default:
			return fmt.Errorf(`failed to recognise value %d`, v)
		}
	}

	return fmt.Errorf("failed to scan value %s", value)
}
