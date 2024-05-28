package main

import "net/http"

type want struct {
	code        int
	body        string
	contentType string
}

type testData struct {
	name        string
	request     string
	requestBody string
	contentType string
	method      string
	want        want
}

func getTestData() []testData {
	return []testData{
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
		{
			name:    "Pos udpate Counter JSON",
			request: "/update/",
			requestBody: `
			{"id":"testCounterJSON",
			"type":"counter",
			"delta":5
			}
			`,
			contentType: "application/json",
			method:      http.MethodPost,
			want: want{code: 200, contentType: "application/json", body: `
			{"id":"testCounterJSON",
			"type":"counter",
			"delta":5
			}
			`},
		},
		{
			name:    "Pos Get Counter JSON",
			request: "/value/",
			requestBody: `
			{"id":"MetricCounter",
			"type":"counter"
			}
			`,
			contentType: "application/json",
			method:      http.MethodPost,
			want: want{code: 200, contentType: "application/json", body: `
			{"id":"MetricCounter",
			"type":"counter",
			"delta":5
			}
			`},
		},
		{
			name:    "Pos Get Gauge JSON",
			request: "/value/",
			requestBody: `
			{"id":"MetricGauge",
			"type":"gauge"
			}
			`,
			contentType: "application/json",
			method:      http.MethodPost,
			want: want{code: 200, contentType: "application/json", body: `
			{"id":"MetricGauge",
			"type":"gauge",
			"value":10
			}
			`},
		},
	}
}
