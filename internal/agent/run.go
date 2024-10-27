package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/MikeRez0/ypmetrics/internal/config"
	"github.com/MikeRez0/ypmetrics/internal/logger"
	"github.com/MikeRez0/ypmetrics/internal/model"
	"github.com/MikeRez0/ypmetrics/internal/utils/retrier"
	"github.com/MikeRez0/ypmetrics/internal/utils/signer"
	"go.uber.org/zap"
)

func Run() error {
	conf, err := config.NewConfigAgent()
	if err != nil {
		return fmt.Errorf("error while load config: %w", err)
	}

	log := logger.GetLogger(conf.LogLevel)

	var metricStore = NewMetricStore()

	tickerPoll := time.NewTicker(time.Duration(conf.PollInterval) * time.Second)
	tickerReport := time.NewTicker(time.Duration(conf.ReportInterval) * time.Second)

	for {
		select {
		case <-tickerPoll.C:
			poll(metricStore)
		case <-tickerReport.C:
			reportBatch(metricStore, conf.HostString, log, conf.SignKey)
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

func report(metricStore *MetricStore, serverURL string, log *zap.Logger, keyHash string) {
	serverURL = "http://" + serverURL

	metricType := model.MetricType(model.CounterType)
	for metricName, val := range metricStore.MetricsCounter {
		metric := model.Metrics{ID: metricName, MType: metricType, Delta: (*int64)(&val)}

		err := sendMetricJSON(serverURL, metric, log, keyHash)
		if err != nil {
			log.Error("error sending counter metric json", zap.Error(err))
		}
	}

	metricType = model.MetricType(model.GaugeType)
	for metricName, val := range metricStore.MetricsGauge {
		metric := model.Metrics{ID: metricName, MType: metricType, Value: (*float64)(&val)}
		err := sendMetricJSON(serverURL, metric, log, keyHash)
		if err != nil {
			log.Error("error sending guage metric json", zap.Error(err))
		}
	}
}

func reportBatch(metricStore *MetricStore, serverURL string, log *zap.Logger, keyHash string) {
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

	err := sendMetricBatchJSON(serverURL, metrics, log, keyHash)
	if err != nil {
		log.Error("error sending guage metric json", zap.Error(err))
	}
}

func sendJSON(requestStr string, jsonStr []byte, log *zap.Logger, keyHash string) error {
	req, err := http.NewRequest(http.MethodPost, requestStr, bytes.NewBuffer(jsonStr))
	if err != nil {
		return fmt.Errorf("error on %s : %w", requestStr, err)
	}
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Content-Type", "application/json")

	if keyHash != "" {
		sgn := signer.NewSigner(keyHash)
		h, err := sgn.GetHashBA(jsonStr)
		if err != nil {
			return fmt.Errorf("signer error: %w", err)
		}

		log.Debug("Hash value", zap.String("Hash", h))
		req.Header.Add(model.HeaderSignerHash, h)
	}

	return retrier.Retry(context.Background(), func() error { //nolint:wrapcheck //error from callback
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("error on %s : %w", requestStr, err)
		}
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("bad response %v for request %s", resp.StatusCode, requestStr)
		}
		return nil
	}, 3, log)
}

func sendMetricJSON(serverURL string, metric model.Metrics, log *zap.Logger, keyHash string) error {
	requestStr := serverURL + "/update/"

	jsonStr, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("erron while json encode: %w", err)
	}

	return sendJSON(requestStr, jsonStr, log, keyHash)
}

func sendMetricBatchJSON(serverURL string, metrics []model.Metrics, log *zap.Logger, keyHash string) error {
	requestStr := serverURL + "/updates/"

	jsonStr, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("erron while json encode: %w", err)
	}

	return sendJSON(requestStr, jsonStr, log, keyHash)
}
