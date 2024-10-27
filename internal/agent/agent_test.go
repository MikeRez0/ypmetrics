package agent

import (
	"testing"

	"github.com/MikeRez0/ypmetrics/internal/logger"
	"github.com/MikeRez0/ypmetrics/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestReadRuntimeMetrics(t *testing.T) {
	ms := NewMetricStore()
	ReadRuntimeMetrics(ms)
	for _, v := range runtimeMetricNames {
		assert.Contains(t, ms.MetricsGauge, v)
	}
}

func Test_poll(t *testing.T) {
	ms := NewMetricStore()
	poll(ms)
	assert.Contains(t, ms.MetricsCounter, "PollCount")
	assert.Contains(t, ms.MetricsGauge, "RandomValue")
}
func Test_report(t *testing.T) {
	ms := NewMetricStore()
	ms.PushCounterMetric("TestCounter", model.CounterValue(10))
	ms.PushGaugeMetric("TestGauge", model.GaugeValue(15))

	log, _ := logger.Initialize("info")

	tests := []struct {
		name string
	}{
		{name: "test1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report(ms, `localhost:8080`, log, "")
		})
	}
}
