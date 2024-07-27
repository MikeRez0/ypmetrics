package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/MikeRez0/ypmetrics/internal/config"
	"github.com/MikeRez0/ypmetrics/internal/model"
)

type Config struct {
	HostString     string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func Run() error {
	conf, err := config.NewConfigAgent()
	if err != nil {
		return fmt.Errorf("error while load config: %w", err)
	}

	var metricStore = NewMetricStore()

	tickerPoll := time.NewTicker(time.Duration(conf.PollInterval) * time.Second)
	tickerReport := time.NewTicker(time.Duration(conf.ReportInterval) * time.Second)

	for {
		select {
		case <-tickerPoll.C:
			poll(metricStore)
		case <-tickerReport.C:
			reportBatch(metricStore, conf.HostString)
			clear(metricStore.MetricsGauge)
			clear(metricStore.MetricsCounter)
		}
	}
}

func poll(metricStore *MetricStore) {
	ReadRuntimeMetrics(metricStore)

	metricStore.PushCounterMetric("PollCount", model.CounterValue(1))
	metricStore.PushGaugeMetric("RandomValue", model.GaugeValue(rand.Float64()*1_000))
}

func report(metricStore *MetricStore, serverURL string) {
	serverURL = "http://" + serverURL

	metricType := model.MetricType(model.CounterType)
	for metricName, val := range metricStore.MetricsCounter {
		metric := model.Metrics{ID: metricName, MType: metricType, Delta: (*int64)(&val)}

		err := sendMetricJSON(serverURL, metric)
		if err != nil {
			log.Println(err)
		}
	}

	metricType = model.MetricType(model.GaugeType)
	for metricName, val := range metricStore.MetricsGauge {
		metric := model.Metrics{ID: metricName, MType: metricType, Value: (*float64)(&val)}
		err := sendMetricJSON(serverURL, metric)
		if err != nil {
			log.Println(err)
		}
	}
}

func reportBatch(metricStore *MetricStore, serverURL string) {
	serverURL = "http://" + serverURL

	metrics := make([]model.Metrics, 0)

	metricType := model.MetricType(model.CounterType)
	for metricName, val := range metricStore.MetricsCounter {
		metric := model.Metrics{ID: metricName, MType: metricType, Delta: (*int64)(&val)}
		metrics = append(metrics, metric)
	}
	metricType = model.MetricType(model.GaugeType)
	for metricName, val := range metricStore.MetricsGauge {
		metric := model.Metrics{ID: metricName, MType: metricType, Value: (*float64)(&val)}
		metrics = append(metrics, metric)
	}

	err := sendMetricBatchJSON(serverURL, metrics)
	if err != nil {
		log.Println(err)
	}
}

func sendJSON(requestStr string, jsonStr []byte) error {
	req, err := http.NewRequest(http.MethodPost, requestStr, bytes.NewBuffer(jsonStr))
	if err != nil {
		return fmt.Errorf("error on %s : %w", requestStr, err)
	}
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error on %s : %w", requestStr, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response %v for request %s", resp.StatusCode, requestStr)
	}
	return nil
}

func sendMetricJSON(serverURL string, metric model.Metrics) error {
	requestStr := serverURL + "/update/"

	jsonStr, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("erron while json encode: %w", err)
	}

	return sendJSON(requestStr, jsonStr)
}

func sendMetricBatchJSON(serverURL string, metrics []model.Metrics) error {
	requestStr := serverURL + "/updates/"

	jsonStr, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("erron while json encode: %w", err)
	}

	return sendJSON(requestStr, jsonStr)
}
