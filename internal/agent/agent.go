package agent

import (
	"runtime"

	"github.com/MikeRez0/ypmetrics/internal/storage"
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

func ReadRuntimeMetrics(metrics *MetricStore) *MetricStore {
	var memStats runtime.MemStats

	runtime.ReadMemStats(&memStats)

	metrics.PushGaugeMetric(`Alloc`, storage.GaugeValue(memStats.Alloc))
	metrics.PushGaugeMetric(`BuckHashSys`, storage.GaugeValue(memStats.BuckHashSys))
	metrics.PushGaugeMetric(`Frees`, storage.GaugeValue(memStats.Frees))
	metrics.PushGaugeMetric(`GCCPUFraction`, storage.GaugeValue(memStats.GCCPUFraction))
	metrics.PushGaugeMetric(`GCSys`, storage.GaugeValue(memStats.GCSys))
	metrics.PushGaugeMetric(`HeapAlloc`, storage.GaugeValue(memStats.HeapAlloc))
	metrics.PushGaugeMetric(`HeapIdle`, storage.GaugeValue(memStats.HeapIdle))
	metrics.PushGaugeMetric(`HeapInuse`, storage.GaugeValue(memStats.HeapInuse))
	metrics.PushGaugeMetric(`HeapObjects`, storage.GaugeValue(memStats.HeapObjects))
	metrics.PushGaugeMetric(`HeapReleased`, storage.GaugeValue(memStats.HeapReleased))
	metrics.PushGaugeMetric(`HeapSys`, storage.GaugeValue(memStats.HeapSys))
	metrics.PushGaugeMetric(`LastGC`, storage.GaugeValue(memStats.LastGC))
	metrics.PushGaugeMetric(`Lookups`, storage.GaugeValue(memStats.Lookups))
	metrics.PushGaugeMetric(`MCacheInuse`, storage.GaugeValue(memStats.MCacheInuse))
	metrics.PushGaugeMetric(`MCacheSys`, storage.GaugeValue(memStats.MCacheSys))
	metrics.PushGaugeMetric(`MSpanInuse`, storage.GaugeValue(memStats.MSpanInuse))
	metrics.PushGaugeMetric(`MSpanSys`, storage.GaugeValue(memStats.MSpanSys))
	metrics.PushGaugeMetric(`Mallocs`, storage.GaugeValue(memStats.Mallocs))
	metrics.PushGaugeMetric(`NextGC`, storage.GaugeValue(memStats.NextGC))
	metrics.PushGaugeMetric(`NumForcedGC`, storage.GaugeValue(memStats.NumForcedGC))
	metrics.PushGaugeMetric(`NumGC`, storage.GaugeValue(memStats.NumGC))
	metrics.PushGaugeMetric(`OtherSys`, storage.GaugeValue(memStats.OtherSys))
	metrics.PushGaugeMetric(`PauseTotalNs`, storage.GaugeValue(memStats.PauseTotalNs))
	metrics.PushGaugeMetric(`StackInuse`, storage.GaugeValue(memStats.StackInuse))
	metrics.PushGaugeMetric(`StackSys`, storage.GaugeValue(memStats.StackSys))
	metrics.PushGaugeMetric(`Sys`, storage.GaugeValue(memStats.Sys))
	metrics.PushGaugeMetric(`TotalAlloc`, storage.GaugeValue(memStats.TotalAlloc))

	return metrics
}
