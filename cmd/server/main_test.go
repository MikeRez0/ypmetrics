package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MikeRez0/ypmetrics/internal/handlers"
	"github.com/MikeRez0/ypmetrics/internal/storage"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestMetricsHandler_Server(t *testing.T) {
	type want struct {
		code        int
		body        string
		contentType string
	}

	testMemStorage := func() *storage.MemStorage {
		store := storage.NewMemStorage()

		store.UpdateCounter("MetricCounter", 5)
		store.UpdateGauge("MetricGauge", 10)

		return store
	}

	mh := &handlers.MetricsHandler{
		Store: testMemStorage(),
	}

	router := setupRouter(mh)
	srv := httptest.NewServer(router)

	tests := []struct {
		name    string
		request string
		method  string
		want    want
	}{
		{
			name:    "Pos update Gauge",
			request: "/update/gauge/test/1",
			method:  http.MethodPost,
			want:    want{code: 200},
		},
		{
			name:    "Pos udpate Counter",
			request: "/update/counter/test2/1",
			method:  http.MethodPost,
			want:    want{code: 200},
		},
		{
			name:    "Pos Get Counter",
			request: "/value/counter/MetricCounter",
			method:  http.MethodGet,
			want:    want{code: 200, body: "5"},
		},
		{
			name:    "Neg Get Counter",
			request: "/value/counter/XXXMetricCounter",
			method:  http.MethodGet,
			want:    want{code: http.StatusNotFound},
		},
		{
			name:    "Neg update Gauge",
			request: "/update/gauge/test/xxx",
			method:  http.MethodPost,
			want:    want{code: 400},
		},
		{
			name:    "Neg udpate Counter",
			request: "/update/counter/test2/xxx",
			method:  http.MethodPost,
			want:    want{code: 400},
		},
		{
			name:    "Neg udpate XXX",
			request: "/update/XXX/test2/5",
			method:  http.MethodPost,
			want:    want{code: 400},
		},
		{
			name:    "Neg udpate No metric",
			request: "/update/counter/",
			method:  http.MethodPost,
			want:    want{code: 404},
		},
		{
			name:    "Neg method not allowed",
			request: "/update/counter/test/123",
			method:  http.MethodPatch,
			want:    want{code: http.StatusMethodNotAllowed},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.method
			req.URL = srv.URL + tt.request

			res, err := req.Send()
			assert.NoError(t, err)
			assert.Equal(t, tt.want.code, res.StatusCode())

			if tt.want.body != "" {
				assert.Equal(t, tt.want.body, string(res.Body()))
			}
			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, res.Header().Get("Content-Type"))
			}
		})
	}
}
