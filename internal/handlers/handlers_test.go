package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MikeRez0/ypmetrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricsHandler_ServeHTTP(t *testing.T) {
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
			mh := &MetricsHandler{
				store: testMemStorage(),
			}

			request := httptest.NewRequest(tt.method, tt.request, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()

			mh.ServeHTTP(w, request)

			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)

			defer res.Body.Close()
			if tt.want.body != "" {
				resBody, err := io.ReadAll(res.Body)

				require.NoError(t, err)

				assert.Equal(t, tt.want.body, string(resBody))
			}
			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
