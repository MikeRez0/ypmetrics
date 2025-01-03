package main

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/MikeRez0/ypmetrics/internal/handlers"
	"github.com/MikeRez0/ypmetrics/internal/model"
	"github.com/MikeRez0/ypmetrics/internal/storage"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestMetricsHandler_Server(t *testing.T) {
	testMemStorage := func() *storage.MemStorage {
		store := storage.NewMemStorage()

		cval, err := store.UpdateCounter(context.Background(), "MetricCounter", 5)
		assert.NoError(t, err)
		assert.Equal(t, model.CounterValue(5), cval)

		gval, err := store.UpdateGauge(context.Background(), "MetricGauge", 10)
		assert.NoError(t, err)
		assert.Equal(t, model.GaugeValue(10), gval)

		return store
	}

	mh := &handlers.MetricsHandler{
		Store: testMemStorage(),
	}

	router := setupRouter(mh, zap.L())
	srv := httptest.NewServer(router)

	tests := getTestData()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.method
			req.URL = srv.URL + tt.request
			if tt.requestBody != "" {
				req.Body = tt.requestBody
			}
			if tt.contentType != "" {
				req.SetHeader("Content-Type", tt.contentType)
			}

			res, err := req.Send()
			assert.NoError(t, err)
			assert.Equal(t, tt.want.code, res.StatusCode())

			if tt.want.body != "" {
				if tt.contentType != "application/json" {
					assert.Equal(t, tt.want.body, string(res.Body()))
				} else {
					assert.JSONEq(t, tt.want.body, string(res.Body()))
				}
			}
			if tt.want.contentType != "" {
				assert.Contains(t, res.Header().Get("Content-Type"), tt.want.contentType)
			}
		})
	}
}
