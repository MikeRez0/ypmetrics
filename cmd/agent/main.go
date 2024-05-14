package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/MikeRez0/ypmetrics/internal/agent"
	"github.com/MikeRez0/ypmetrics/internal/storage"
)

func main() {

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	hostString := flag.String("a", `localhost:8080`, "HTTP server endpoint")
	pollInterval := flag.Int("p", 2, "Poll interval")
	reportInterval := flag.Int("r", 10, "Report interval")
	flag.Parse()

	var metricStore = agent.NewMetricStore()

	for i := 1; ; i++ {
		poll(metricStore)
		if i*(*pollInterval) >= *reportInterval {
			report(metricStore, *hostString)
			clear(metricStore.Metrics)
			i = 0
		}
		time.Sleep(time.Second * time.Duration(*pollInterval))
	}

	// return nil
}

func poll(metricStore *agent.MetricStore) {
	agent.ReadRuntimeMetrics(metricStore)

	metricStore.PushCounterMetric("PollCount", storage.CounterValue(1))
	metricStore.PushGaugeMetric("RandomValue", storage.GaugeValue(rand.Float64()*1_000))
}

func report(metricStore *agent.MetricStore, serverURL string) {
	serverURL = "http://" + serverURL

	for name, metric := range metricStore.Metrics {

		metricType := metric.MetricType
		metricName := name
		value := metric.Value
		requestStr := serverURL + "/update/" + metricType + "/" + metricName + "/"
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
