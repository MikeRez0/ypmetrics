package main

import (
	"testing"

	"github.com/MikeRez0/ypmetrics/internal/agent"
	"github.com/MikeRez0/ypmetrics/internal/storage"
	"github.com/stretchr/testify/assert"
)

func Test_poll(t *testing.T) {
	ms := agent.NewMetricStore()
	poll(ms)
	assert.Contains(t, ms.Metrics, "PollCount")
	assert.Contains(t, ms.Metrics, "RandomValue")
}
func Test_report(t *testing.T) {
	ms := agent.NewMetricStore()
	ms.PushCounterMetric("TestCounter", storage.CounterValue(10))
	ms.PushGaugeMetric("TestGauge", storage.GaugeValue(15))

	tests := []struct {
		name string
	}{
		{name: "test1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report(ms, `localhost:8080`)
		})
	}
}
