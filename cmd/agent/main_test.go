package main

import (
	"testing"

	"github.com/MikeRez0/ypmetrics/internal/agent"
	"github.com/stretchr/testify/assert"
)

func Test_poll(t *testing.T) {
	ms := agent.NewMetricStore()
	poll(ms)
	assert.Contains(t, ms.Metrics, "PollCount")
	assert.Contains(t, ms.Metrics, "RandomValue")

}

func Test_report(t *testing.T) {
	type args struct {
		metricStore *agent.MetricStore
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report(tt.args.metricStore, `localhost:8080`)
		})
	}
}
