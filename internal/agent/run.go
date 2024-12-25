package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
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

	log.Info(fmt.Sprintf("start agent with config: %v", conf))

	var metricStore = NewMetricStore()

	ctx := context.Background()

	var wg sync.WaitGroup

	wg.Add(1)
	jobStart(ctx, func() error {
		poll(metricStore)
		return nil
	}, time.Duration(conf.PollInterval)*time.Second, 1, log.Named("Poll go metrics job"))

	jobStart(ctx, func() error {
		ReadGopsutilMetrics(metricStore)
		return nil
	}, time.Duration(conf.PollInterval)*time.Second, 1, log.Named("Poll gopsutil metrics job"))

	jobStart(ctx, func() error {
		reportBatch(metricStore, conf.HostString, log, conf.SignKey)
		return nil
	}, time.Duration(conf.ReportInterval)*time.Second, conf.RateLimit, log.Named("Report metrics job"))

	wg.Wait()
	return nil
}

// ticker with worker pool.
func jobStart(ctx context.Context, job func() error, interval time.Duration, workers int, log *zap.Logger) {
	jobFire := make(chan struct{})

	for i := 0; i < workers; i++ {
		i := i
		go func(j chan struct{}) {
			log.Debug(fmt.Sprintf("Worker %d init", i))
			for {
				select {
				case <-j:
					log.Debug(fmt.Sprintf("Worker %d start job", i))
					err := job()
					if err != nil {
						log.Error("job finished with error", zap.Error(err))
					}
					log.Debug(fmt.Sprintf("Worker %d end job", i))
				case <-ctx.Done():
					log.Debug(fmt.Sprintf("Worker %d stopped", i))
				}
			}
		}(jobFire)
	}

	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				jobFire <- struct{}{}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func poll(metricStore *MetricStore) {
	ReadRuntimeMetrics(metricStore)

	metricStore.PushCounterMetric("PollCount", model.CounterValue(1))
	metricStore.PushGaugeMetric("RandomValue", model.GaugeValue(rand.Float64()*1_000))
}

func report(metricStore *MetricStore, serverURL string, log *zap.Logger, keyHash string) {
	serverURL = "http://" + serverURL

	metricType := model.MetricType(model.CounterType)
	for metricName, val := range metricStore.GetCounterMetrics() {
		metric := model.Metrics{ID: metricName, MType: metricType, Delta: (*int64)(&val)}

		err := sendMetricJSON(serverURL, metric, log, keyHash)
		if err != nil {
			log.Error("error sending counter metric json", zap.Error(err))
		}
	}

	metricType = model.MetricType(model.GaugeType)
	for metricName, val := range metricStore.GetGaugeMetrics() {
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
	for metricName, val := range metricStore.GetCounterMetrics() {
		metric := model.Metrics{ID: metricName, MType: metricType, Delta: (*int64)(&val)}
		metrics = append(metrics, metric)
	}
	metricType = model.MetricType(model.GaugeType)
	for metricName, val := range metricStore.GetGaugeMetrics() {
		metric := model.Metrics{ID: metricName, MType: metricType, Value: (*float64)(&val)}
		metrics = append(metrics, metric)
	}

	err := sendMetricBatchJSON(serverURL, metrics, log, keyHash)
	if err != nil {
		log.Error("error sending guage metric json", zap.Error(err))
	}

	metricStore.Clear()
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
	}, 3, log.Named("http request"))
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
