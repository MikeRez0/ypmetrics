package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MikeRez0/ypmetrics/internal/config"
	"github.com/MikeRez0/ypmetrics/internal/logger"
	"github.com/MikeRez0/ypmetrics/internal/model"
)

func TestReadRuntimeMetrics(t *testing.T) {
	app := NewAgentApp(config.ConfigAgent{}, logger.GetLogger("info"))
	app.ReadRuntimeMetrics()
	for _, v := range runtimeMetricNames {
		assert.Contains(t, app.metrics.GetGaugeMetrics(), v)
	}
}

func Test_poll(t *testing.T) {
	app := NewAgentApp(config.ConfigAgent{}, logger.GetLogger("info"))
	app.Poll()
	assert.Contains(t, app.metrics.GetCounterMetrics(), "PollCount")
	assert.Contains(t, app.metrics.GetGaugeMetrics(), "RandomValue")
}
func Test_report(t *testing.T) {
	ms := NewMetricStore()
	ms.PushCounterMetric("TestCounter", model.CounterValue(10))
	ms.PushGaugeMetric("TestGauge", model.GaugeValue(15))

	tests := []struct {
		name string
	}{
		{name: "test1"},
	}
	app := NewAgentApp(config.ConfigAgent{}, logger.GetLogger("info"))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.Report()
		})
	}
}
