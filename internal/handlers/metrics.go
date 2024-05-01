package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/MikeRez0/ypmetrics/internal/storage"
)

func (mh *MetricsHandler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	matches := UpdateMetricRe.FindStringSubmatch(r.URL.Path)
	fmt.Println(matches)
	//Expected fullstring + groups: [1]:metricType, [2]:metric, [3]:value
	if len(matches) < 4 {
		BadRequestErrorHandler(w, r)
		return
	}

	var (
		metricType = matches[1]
		metric     = matches[2]
		valueRaw   = matches[3]
	)

	switch metricType {
	case storage.GaugeType:
		value, err := strconv.ParseFloat(valueRaw, 64)
		if err != nil {
			BadRequestErrorHandler(w, r)
			return
		}
		mh.store.UpdateGauge(metric, storage.GaugeValue(value))
	case storage.CounterType:
		value, err := strconv.ParseInt(valueRaw, 10, 64)
		if err != nil {
			BadRequestErrorHandler(w, r)
			return
		}
		mh.store.UpdateCounter(metric, storage.CounterValue(value))
	default:
		NotFoundErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (mh *MetricsHandler) GetMetric(w http.ResponseWriter, r *http.Request) {
	matches := GetMetricRe.FindStringSubmatch(r.URL.Path)
	//Expected fullstring + groups: [1]:metricType, [2]:metric
	if len(matches) < 3 {
		BadRequestErrorHandler(w, r)
		return
	}

	var (
		metricType = matches[1]
		metric     = matches[2]
	)

	switch metricType {
	case storage.GaugeType:
		value, err := mh.store.GetGauge(metric)

		if err != nil {
			BadRequestErrorHandler(w, r)
			return
		}
		io.WriteString(w, strconv.FormatFloat(float64(value), 'f', 5, 64))
	case storage.CounterType:
		value, err := mh.store.GetCounter(metric)
		if err != nil {
			BadRequestErrorHandler(w, r)
			return
		}
		io.WriteString(w, strconv.FormatInt(int64(value), 10))
	default:
		NotFoundErrorHandler(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	// w.WriteHeader(http.StatusOK)

}
