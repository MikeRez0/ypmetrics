package agent

import (
	"fmt"
	"runtime"

	"github.com/MikeRez0/ypmetrics/internal/model"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
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

	metrics.PushGaugeMetric(`Alloc`, model.GaugeValue(memStats.Alloc))
	metrics.PushGaugeMetric(`BuckHashSys`, model.GaugeValue(memStats.BuckHashSys))
	metrics.PushGaugeMetric(`Frees`, model.GaugeValue(memStats.Frees))
	metrics.PushGaugeMetric(`GCCPUFraction`, model.GaugeValue(memStats.GCCPUFraction))
	metrics.PushGaugeMetric(`GCSys`, model.GaugeValue(memStats.GCSys))
	metrics.PushGaugeMetric(`HeapAlloc`, model.GaugeValue(memStats.HeapAlloc))
	metrics.PushGaugeMetric(`HeapIdle`, model.GaugeValue(memStats.HeapIdle))
	metrics.PushGaugeMetric(`HeapInuse`, model.GaugeValue(memStats.HeapInuse))
	metrics.PushGaugeMetric(`HeapObjects`, model.GaugeValue(memStats.HeapObjects))
	metrics.PushGaugeMetric(`HeapReleased`, model.GaugeValue(memStats.HeapReleased))
	metrics.PushGaugeMetric(`HeapSys`, model.GaugeValue(memStats.HeapSys))
	metrics.PushGaugeMetric(`LastGC`, model.GaugeValue(memStats.LastGC))
	metrics.PushGaugeMetric(`Lookups`, model.GaugeValue(memStats.Lookups))
	metrics.PushGaugeMetric(`MCacheInuse`, model.GaugeValue(memStats.MCacheInuse))
	metrics.PushGaugeMetric(`MCacheSys`, model.GaugeValue(memStats.MCacheSys))
	metrics.PushGaugeMetric(`MSpanInuse`, model.GaugeValue(memStats.MSpanInuse))
	metrics.PushGaugeMetric(`MSpanSys`, model.GaugeValue(memStats.MSpanSys))
	metrics.PushGaugeMetric(`Mallocs`, model.GaugeValue(memStats.Mallocs))
	metrics.PushGaugeMetric(`NextGC`, model.GaugeValue(memStats.NextGC))
	metrics.PushGaugeMetric(`NumForcedGC`, model.GaugeValue(memStats.NumForcedGC))
	metrics.PushGaugeMetric(`NumGC`, model.GaugeValue(memStats.NumGC))
	metrics.PushGaugeMetric(`OtherSys`, model.GaugeValue(memStats.OtherSys))
	metrics.PushGaugeMetric(`PauseTotalNs`, model.GaugeValue(memStats.PauseTotalNs))
	metrics.PushGaugeMetric(`StackInuse`, model.GaugeValue(memStats.StackInuse))
	metrics.PushGaugeMetric(`StackSys`, model.GaugeValue(memStats.StackSys))
	metrics.PushGaugeMetric(`Sys`, model.GaugeValue(memStats.Sys))
	metrics.PushGaugeMetric(`TotalAlloc`, model.GaugeValue(memStats.TotalAlloc))

	return metrics
}

func ReadGopsutilMetrics(metrics *MetricStore) *MetricStore {
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Percent(0, true)

	metrics.PushGaugeMetric("TotalMemory", model.GaugeValue(v.Total))
	metrics.PushGaugeMetric("FreeMemory", model.GaugeValue(v.Free))
	for i, u := range c {
		metrics.PushGaugeMetric(fmt.Sprintf("CPUutilization%d", i), model.GaugeValue(u))
	}

	return metrics
}
