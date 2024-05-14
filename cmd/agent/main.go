package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"strconv"
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
	HostString     string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func run() error {
	hostString := flag.String("a", `localhost:8080`, "HTTP server endpoint")
	pollInterval := flag.Int("p", 2, "Poll interval")
	reportInterval := flag.Int("r", 10, "Report interval")
	flag.Parse()

	config := Config{HostString: *hostString, PollInterval: *pollInterval, ReportInterval: *reportInterval}
	err := env.Parse(&config)
	if err != nil {
		return err
	}

	var metricStore = agent.NewMetricStore()

	for i := 1; ; i++ {
		poll(metricStore)
		if i*(config.PollInterval) >= config.ReportInterval {
			report(metricStore, config.HostString)
			clear(metricStore.Metrics)
			i = 0
		}
		time.Sleep(time.Second * time.Duration(config.PollInterval))
	}
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
			if v, ok := value.(storage.CounterValue); ok {
				requestStr += strconv.FormatInt(int64(v), 10)
			}
		} else if metricType == storage.GaugeType {
			if v, ok := value.(storage.GaugeValue); ok {
				requestStr += strconv.FormatFloat(float64(v), 'f', 5, 64)
			}
		}
		resp, err := http.Post(requestStr, "text/plain", nil)
		if err != nil {
			log.Println(requestStr)
			log.Println(err)
			return
		}
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != http.StatusOK {
			log.Println(requestStr)
			log.Println("Status from server:", resp.StatusCode)
		}
	}
}
