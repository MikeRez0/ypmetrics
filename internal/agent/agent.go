// Package agent - agent application for collecting and reporting metrics.
package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"go.uber.org/zap"

	"github.com/MikeRez0/ypmetrics/internal/config"
	"github.com/MikeRez0/ypmetrics/internal/model"
	"github.com/MikeRez0/ypmetrics/internal/utils/retrier"
	"github.com/MikeRez0/ypmetrics/internal/utils/signer"
)

var runtimeMetricNames []string = []string{
	`Alloc`,
	`BuckHashSys`, `Frees`,
	`GCCPUFraction`, `GCSys`, `HeapAlloc`,
	`HeapIdle`, `HeapInuse`, `HeapObjects`, `HeapReleased`,
	`HeapSys`, `LastGC`, `Lookups`, `MCacheInuse`,
	`MCacheSys`, `MSpanInuse`, `MSpanSys`,
	`Mallocs`, `NextGC`, `NumForcedGC`, `NumGC`, `OtherSys`, `PauseTotalNs`,
	`StackInuse`, `StackSys`, `Sys`, `TotalAlloc`}

// AgentApp - Agent application.
type AgentApp struct {
	log     *zap.Logger
	metrics *MetricStore
	retrier *retrier.Retrier
	host    string
	keyHash string
}

// NewAgentApp - Create new agent application.
func NewAgentApp(conf config.ConfigAgent, log *zap.Logger) *AgentApp {
	r := retrier.NewRetrier(log.Named("Retrier"), 3, 3)
	return &AgentApp{
		log:     log,
		metrics: NewMetricStore(),
		retrier: r,
		host:    conf.HostString,
		keyHash: conf.SignKey,
	}
}

// ReadRuntimeMetrics - read runtime metrics.
//
// Metric list:
// `Alloc`, `BuckHashSys`, `Frees`,
// `GCCPUFraction`, `GCSys`, `HeapAlloc`,
// `HeapIdle`, `HeapInuse`, `HeapObjects`, `HeapReleased`,
// `HeapSys`, `LastGC`, `Lookups`, `MCacheInuse`,
// `MCacheSys`, `MSpanInuse`, `MSpanSys`,
// `Mallocs`, `NextGC`, `NumForcedGC`, `NumGC`, `OtherSys`, `PauseTotalNs`,
// `StackInuse`, `StackSys`, `Sys`, `TotalAlloc`.
func (a *AgentApp) ReadRuntimeMetrics() *MetricStore {
	var memStats runtime.MemStats

	runtime.ReadMemStats(&memStats)

	a.metrics.PushGaugeMetric(`Alloc`, model.GaugeValue(memStats.Alloc))
	a.metrics.PushGaugeMetric(`BuckHashSys`, model.GaugeValue(memStats.BuckHashSys))
	a.metrics.PushGaugeMetric(`Frees`, model.GaugeValue(memStats.Frees))
	a.metrics.PushGaugeMetric(`GCCPUFraction`, model.GaugeValue(memStats.GCCPUFraction))
	a.metrics.PushGaugeMetric(`GCSys`, model.GaugeValue(memStats.GCSys))
	a.metrics.PushGaugeMetric(`HeapAlloc`, model.GaugeValue(memStats.HeapAlloc))
	a.metrics.PushGaugeMetric(`HeapIdle`, model.GaugeValue(memStats.HeapIdle))
	a.metrics.PushGaugeMetric(`HeapInuse`, model.GaugeValue(memStats.HeapInuse))
	a.metrics.PushGaugeMetric(`HeapObjects`, model.GaugeValue(memStats.HeapObjects))
	a.metrics.PushGaugeMetric(`HeapReleased`, model.GaugeValue(memStats.HeapReleased))
	a.metrics.PushGaugeMetric(`HeapSys`, model.GaugeValue(memStats.HeapSys))
	a.metrics.PushGaugeMetric(`LastGC`, model.GaugeValue(memStats.LastGC))
	a.metrics.PushGaugeMetric(`Lookups`, model.GaugeValue(memStats.Lookups))
	a.metrics.PushGaugeMetric(`MCacheInuse`, model.GaugeValue(memStats.MCacheInuse))
	a.metrics.PushGaugeMetric(`MCacheSys`, model.GaugeValue(memStats.MCacheSys))
	a.metrics.PushGaugeMetric(`MSpanInuse`, model.GaugeValue(memStats.MSpanInuse))
	a.metrics.PushGaugeMetric(`MSpanSys`, model.GaugeValue(memStats.MSpanSys))
	a.metrics.PushGaugeMetric(`Mallocs`, model.GaugeValue(memStats.Mallocs))
	a.metrics.PushGaugeMetric(`NextGC`, model.GaugeValue(memStats.NextGC))
	a.metrics.PushGaugeMetric(`NumForcedGC`, model.GaugeValue(memStats.NumForcedGC))
	a.metrics.PushGaugeMetric(`NumGC`, model.GaugeValue(memStats.NumGC))
	a.metrics.PushGaugeMetric(`OtherSys`, model.GaugeValue(memStats.OtherSys))
	a.metrics.PushGaugeMetric(`PauseTotalNs`, model.GaugeValue(memStats.PauseTotalNs))
	a.metrics.PushGaugeMetric(`StackInuse`, model.GaugeValue(memStats.StackInuse))
	a.metrics.PushGaugeMetric(`StackSys`, model.GaugeValue(memStats.StackSys))
	a.metrics.PushGaugeMetric(`Sys`, model.GaugeValue(memStats.Sys))
	a.metrics.PushGaugeMetric(`TotalAlloc`, model.GaugeValue(memStats.TotalAlloc))

	return a.metrics
}

// ReadGopsutilMetrics - read metrics:
//
// TotalMemory, FreeMemory, CPUutilization(i) for every CPU-core.
func (a *AgentApp) ReadGopsutilMetrics() *MetricStore {
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Percent(0, true)

	a.metrics.PushGaugeMetric("TotalMemory", model.GaugeValue(v.Total))
	a.metrics.PushGaugeMetric("FreeMemory", model.GaugeValue(v.Free))
	for i, u := range c {
		a.metrics.PushGaugeMetric(fmt.Sprintf("CPUutilization%d", i), model.GaugeValue(u))
	}

	return a.metrics
}

// Poll - read metrics and add metrics:
//
// PollCount = 1
//
// RandomValue = random value 0..1000.
func (a *AgentApp) Poll() {
	a.ReadRuntimeMetrics()

	a.metrics.PushCounterMetric("PollCount", model.CounterValue(1))
	a.metrics.PushGaugeMetric("RandomValue", model.GaugeValue(rand.Float64()*1_000))
}

// Report - Send metrics to server (one-by-one-request).
func (a *AgentApp) Report() {
	serverURL := "http://" + a.host

	metricType := model.MetricType(model.CounterType)
	for metricName, val := range a.metrics.GetCounterMetrics() {
		metric := model.Metrics{ID: metricName, MType: metricType, Delta: (*int64)(&val)}

		err := a.sendMetricJSON(serverURL, metric)
		if err != nil {
			a.log.Error("error sending counter metric json", zap.Error(err))
		}
	}

	metricType = model.MetricType(model.GaugeType)
	for metricName, val := range a.metrics.GetGaugeMetrics() {
		metric := model.Metrics{ID: metricName, MType: metricType, Value: (*float64)(&val)}
		err := a.sendMetricJSON(serverURL, metric)
		if err != nil {
			a.log.Error("error sending guage metric json", zap.Error(err))
		}
	}
}

// ReportBatch - Send metrics to server (all-in-one-request).
func (a *AgentApp) ReportBatch() {
	serverURL := "http://" + a.host

	metrics := make([]model.Metrics, 0)

	metricType := model.MetricType(model.CounterType)
	for metricName, val := range a.metrics.GetCounterMetrics() {
		metric := model.Metrics{ID: metricName, MType: metricType, Delta: (*int64)(&val)}
		metrics = append(metrics, metric)
	}
	metricType = model.MetricType(model.GaugeType)
	for metricName, val := range a.metrics.GetGaugeMetrics() {
		metric := model.Metrics{ID: metricName, MType: metricType, Value: (*float64)(&val)}
		metrics = append(metrics, metric)
	}

	err := a.sendMetricBatchJSON(serverURL, metrics)
	if err != nil {
		a.log.Error("error sending guage metric json", zap.Error(err))
	}

	a.metrics.Clear()
}

func checkCanRetry(err error) bool {
	return true
}

func (a *AgentApp) sendJSON(requestStr string, jsonStr []byte) error {
	req, err := http.NewRequest(http.MethodPost, requestStr, bytes.NewBuffer(jsonStr))
	if err != nil {
		return fmt.Errorf("error on %s : %w", requestStr, err)
	}
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Content-Type", "application/json")

	if a.keyHash != "" {
		sgn := signer.NewSigner(a.keyHash)
		h, err := sgn.GetHashBA(jsonStr)
		if err != nil {
			return fmt.Errorf("signer error: %w", err)
		}

		a.log.Debug("Hash value", zap.String("Hash", h))
		req.Header.Add(model.HeaderSignerHash, h)
	}

	return a.retrier.Retry(context.Background(), func() error { //nolint:wrapcheck //error from callback
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("error on %s : %w", requestStr, err)
		}
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("bad response %v for request %s", resp.StatusCode, requestStr)
		}
		return nil
	}, checkCanRetry)
}

func (a *AgentApp) sendMetricJSON(serverURL string, metric model.Metrics) error {
	requestStr := serverURL + "/update/"

	jsonStr, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("erron while json encode: %w", err)
	}

	return a.sendJSON(requestStr, jsonStr)
}

func (a *AgentApp) sendMetricBatchJSON(serverURL string, metrics []model.Metrics) error {
	requestStr := serverURL + "/updates/"

	jsonStr, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("erron while json encode: %w", err)
	}

	return a.sendJSON(requestStr, jsonStr)
}
