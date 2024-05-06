package agent

import (
	"reflect"
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

	readStat := func(name string) storage.GaugeValue {
		s := reflect.ValueOf(&memStats).Elem()
		v := s.FieldByName(name)
		// fmt.Println(name, v.Type(), v.Interface())
		if v.CanFloat() {
			metric := s.FieldByName(name).Float()
			return storage.GaugeValue(metric)
		} else if v.CanInt() {
			metric := s.FieldByName(name).Int()
			return storage.GaugeValue(metric)
		} else if v.CanUint() {
			metric := s.FieldByName(name).Uint()
			return storage.GaugeValue(metric)
		}

		return storage.GaugeValue(0)
	}

	for _, mname := range runtimeMetricNames {
		metrics.PushGaugeMetric(mname, storage.GaugeValue(readStat(mname)))
	}

	return metrics
}
