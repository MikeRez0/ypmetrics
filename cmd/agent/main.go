package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/MikeRez0/ypmetrics/internal/agent"
	"github.com/MikeRez0/ypmetrics/internal/storage"
)

func main() {

	run()
}

func run() error {
	pollInterval := 2
	reportInterval := 10

	var metricStore = agent.NewMetricStore()

	for i := 1; ; i++ {
		// fmt.Println("Poll", i)
		poll(metricStore)
		// fmt.Println(metricStore.Metrics)
		if i*pollInterval == reportInterval {
			// fmt.Println("Report...")
			report(metricStore)
			clear(metricStore.Metrics)
			i = 0
		}
		time.Sleep(time.Second * time.Duration(pollInterval))
	}

	// return nil
}

func poll(metricStore *agent.MetricStore) {
	agent.ReadRuntimeMetrics(metricStore)

	metricStore.PushCounterMetric("PollCount", storage.CounterValue(1))
	metricStore.PushGaugeMetric("RandomValue", storage.GaugeValue(rand.Float64()*1_000))
}

func report(metricStore *agent.MetricStore) {
	serverURL := "http://localhost:8080/"

	for name, metric := range metricStore.Metrics {

		metricType := metric.MetricType
		metricName := name
		value := metric.Value
		requestStr := serverURL + "update/" + metricType + "/" + metricName + "/"
		if metricType == storage.CounterType {
			requestStr += fmt.Sprint(value)
		} else if metricType == storage.GaugeType {
			requestStr += fmt.Sprintf("%.5f", value)
		}
		resp, err := http.Post(requestStr, "text/plain", nil)
		if err != nil {
			fmt.Println(requestStr)
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Println(requestStr)
			fmt.Println("Status from server:", resp.StatusCode)
		}
	}
}
