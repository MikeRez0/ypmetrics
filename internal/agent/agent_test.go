package agent

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MikeRez0/ypmetrics/internal/config"
	"github.com/MikeRez0/ypmetrics/internal/logger"
	"github.com/MikeRez0/ypmetrics/internal/model"
)

func TestReadRuntimeMetrics(t *testing.T) {
	app, err := NewAgentApp(&config.ConfigAgent{}, logger.GetLogger("info"))
	assert.NoError(t, err)
	app.ReadRuntimeMetrics()
	for _, v := range runtimeMetricNames {
		assert.Contains(t, app.metrics.GetGaugeMetrics(), v)
	}
}

func Test_poll(t *testing.T) {
	app, err := NewAgentApp(&config.ConfigAgent{}, logger.GetLogger("info"))
	assert.NoError(t, err)
	app.Poll()
	assert.Contains(t, app.metrics.GetCounterMetrics(), "PollCount")
	assert.Contains(t, app.metrics.GetGaugeMetrics(), "RandomValue")
}

func Test_ReadGopsutil(t *testing.T) {
	app, err := NewAgentApp(&config.ConfigAgent{}, logger.GetLogger("info"))
	assert.NoError(t, err)
	ms := app.ReadGopsutilMetrics()
	assert.Contains(t, ms.GetGaugeMetrics(), "TotalMemory")
	assert.Contains(t, ms.GetGaugeMetrics(), "FreeMemory")
	assert.Contains(t, ms.GetGaugeMetrics(), "CPUutilization0")
}

func Test_report(t *testing.T) {
	tests := []struct {
		name        string
		fillMetrics func(ms *MetricStore)
		batch       bool
		wantBody    string
		wantStatus  int
	}{
		{name: "Counter report",
			fillMetrics: func(ms *MetricStore) {
				ms.PushCounterMetric("TestCounter", model.CounterValue(10))
			},
			wantStatus: 200,
			wantBody:   `{"type":"counter","delta":10,"id":"TestCounter"}`,
		},
		{name: "Gauge report",
			fillMetrics: func(ms *MetricStore) {
				ms.PushGaugeMetric("TestGauge", model.GaugeValue(15))
			},
			wantStatus: 200,
			wantBody:   `{"type":"gauge","value":15,"id":"TestGauge"}`,
		},
		{name: "Batch report",
			fillMetrics: func(ms *MetricStore) {
				ms.PushCounterMetric("TestCounter", model.CounterValue(10))
				ms.PushGaugeMetric("TestGauge", model.GaugeValue(15))
			},
			batch:      true,
			wantStatus: 200,
			wantBody:   `[{"type":"counter","delta":10,"id":"TestCounter"},{"type":"gauge","value":15,"id":"TestGauge"}]`,
		},
	}

	testID := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := "/update/"
		if tests[testID].batch {
			url = "/updates/"
		}
		assert.Equal(t, url, r.URL.Path)

		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		defer r.Body.Close() //nolint:errcheck // that's why

		assert.Equal(t, tests[testID].wantBody, string(body))

		w.WriteHeader(http.StatusOK)
	}))

	app, err := NewAgentApp(&config.ConfigAgent{HostString: srv.URL[7:], SignKey: "test"}, logger.GetLogger("info"))
	assert.NoError(t, err)
	for i, tt := range tests {
		testID = i
		t.Run(tt.name, func(t *testing.T) {
			tt.fillMetrics(app.metrics)
			if !tt.batch {
				app.Report()
				app.metrics.Clear()
			} else {
				app.ReportBatch()
			}
		})
	}
}
