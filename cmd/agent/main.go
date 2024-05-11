package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/caarlos0/env/v6"

	"github.com/MikeRez0/ypmetrics/internal/agent"
	"github.com/MikeRez0/ypmetrics/internal/storage"
)

func main() {

	if err := run(); err != nil {
		panic(err)
	}
}

type Config struct {
	hostString     string `env:"ADDRESS"`
	reportInterval int    `env:"REPORT_INTERVAL"`
	pollInterval   int    `env:"POLL_INTERVAL"`
}

func run() error {
	hostString := flag.String("a", `localhost:8080`, "HTTP server endpoint")
	pollInterval := flag.Int("p", 2, "Poll interval")
	reportInterval := flag.Int("r", 10, "Report interval")
	flag.Parse()

	config := Config{hostString: *hostString, pollInterval: *pollInterval, reportInterval: *reportInterval}
	env.Parse(&config)

	var metricStore = agent.NewMetricStore()

	for i := 1; ; i++ {
		poll(metricStore)
		if i*(config.pollInterval) >= config.reportInterval {
			report(metricStore, config.hostString)
			clear(metricStore.Metrics)
			i = 0
		}
		time.Sleep(time.Second * time.Duration(config.pollInterval))
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
