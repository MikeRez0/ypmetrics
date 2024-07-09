package model

const GaugeType = "gauge"
const CounterType = "counter"

type GaugeValue float64
type CounterValue int64

type Metrics struct { //nolint:govet //external rule from course author
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
