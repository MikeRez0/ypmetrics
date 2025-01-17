package storage_test

import "net/http"

type want struct {
	body        string
	contentType string
	code        int
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
			name:    "Scenario Test Counter - 1.Pos udpate Counter",
			request: "/update/counter/MetricCounter/5",
			method:  http.MethodPost,
			want:    want{code: 200},
		},
		{
			name:    "Scenario Test Counter - 2.Pos Get Counter",
			request: "/value/counter/MetricCounter",
			method:  http.MethodGet,
			want:    want{code: 200, body: "5"},
		},
		{
			name:    "Scenario Test Counter - 3.Pos udpate Counter",
			request: "/update/counter/MetricCounter/5",
			method:  http.MethodPost,
			want:    want{code: 200},
		},
		{
			name:    "Scenario Test Counter - 4.Pos Get Counter",
			request: "/value/counter/MetricCounter",
			method:  http.MethodGet,
			want:    want{code: 200, body: "10"},
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
			name:    "Pos udpate Batch JSON",
			request: "/updates/",
			requestBody: `
			[{"id":"testCounterJSON",
			"type":"counter",
			"delta":5
			},{"id":"testCounterJSON2",
			"type":"counter",
			"delta":52
			},{"id":"testGaugeJSON",
			"type":"gauge",
			"value":10.2}]
			`,
			contentType: "application/json",
			method:      http.MethodPost,
			want: want{code: 200, contentType: "application/json", body: `
			[{"id":"testCounterJSON",
			"type":"counter",
			"delta":5
			},{"id":"testCounterJSON2",
			"type":"counter",
			"delta":52
			},{"id":"testGaugeJSON",
			"type":"gauge",
			"value":10.2}]
			`},
		},
		{
			name:    "Neg 1 udpate Batch JSON",
			request: "/updates/",
			requestBody: `
			[{"id":"testCounterJSON",
			"type":"counter",
			"delta":qwe
			},{"id":"testCounterJSON2",
			"type":"counter",
			"delta":52
			},{"id":"testGaugeJSON",
			"type":"gauge",
			"delta":10.2}]
			`,
			contentType: "application/json",
			method:      http.MethodPost,
			want:        want{code: 400},
		},
		{
			name:    "Neg 2 udpate Batch JSON",
			request: "/updates/",
			requestBody: `
			[{"id":"testCounterJSON",
			"type":"value",
			"delta":123
			},
			{"id":"testCounterJSON2",
			"type":"counter",
			"delta":52
			}]
			`,
			contentType: "application/json",
			method:      http.MethodPost,
			want:        want{code: 400},
		},
		{
			name:    "Neg Get Counter JSON - bad type",
			request: "/value/",
			requestBody: `
			{"id":"MetricCounter",
			"type":"value"
			}
			`,
			contentType: "application/json",
			method:      http.MethodPost,
			want:        want{code: 400},
		},
		{
			name:    "Neg Get Counter JSON - not found",
			request: "/value/",
			requestBody: `
			{"id":"MetricCounterNOTFOUND",
			"type":"counter"
			}
			`,
			contentType: "application/json",
			method:      http.MethodPost,
			want:        want{code: 404},
		},
		{
			name:    "Scenario test gauge - 1.Pos Update Gauge JSON",
			request: "/update/",
			requestBody: `
			{"id":"MetricGauge",
			"type":"gauge",
			"value":10
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
		{
			name:    "Scenario test gauge - 2.Pos Get Gauge JSON",
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
		// {
		// 	name:    "Pos Get HTML",
		// 	request: "/",
		// 	method:  http.MethodGet,
		// 	want:    want{code: 200},
		// },
	}
}
