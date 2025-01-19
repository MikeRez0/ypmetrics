// Package model contains app model of metrics.
package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

// HeaderSignerHash - Header key in request for hash-body.
const HeaderSignerHash = "HashSHA256"

// HeaderEncryptKey - Header key in request for decrypt body.
const HeaderEncryptKey = "-X-Encrypt"

// GaugeType - name for gauge.
const GaugeType = "gauge"

// CounterType - name for counter.
const CounterType = "counter"

// Database keys for metric types.
const (
	counterTypeDB = iota + 1
	gaugeTypeDB
)

type MetricType string

type GaugeValue float64
type CounterValue int64

// Metrics - structure for metric value.
type Metrics struct {
	MType MetricType `json:"type" binding:"required"` // параметр, принимающий значение gauge или counter
	Delta *int64     `json:"delta,omitempty"`         // значение метрики в случае передачи counter
	Value *float64   `json:"value,omitempty"`         // значение метрики в случае передачи gauge
	ID    string     `json:"id" binding:"required"`   // имя метрики
}

func (mt MetricType) Value() (driver.Value, error) {
	switch mt {
	case CounterType:
		return counterTypeDB, nil
	case GaugeType:
		return gaugeTypeDB, nil
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
		case counterTypeDB:
			*mt = CounterType
			return nil
		case gaugeTypeDB:
			*mt = GaugeType
			return nil
		default:
			return fmt.Errorf(`failed to recognise value %d`, v)
		}
	}

	return fmt.Errorf("failed to scan value %s", value)
}

type BadValueError struct {
	err string
}

func NewErrBadValue(err string) BadValueError {
	return BadValueError{err: err}
}
func (e BadValueError) Error() string {
	return e.err
}
