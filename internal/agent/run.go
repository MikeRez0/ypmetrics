package agent

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/MikeRez0/ypmetrics/internal/config"
	"github.com/MikeRez0/ypmetrics/internal/storage"
)

type Config struct {
	HostString     string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func Run() error {

	conf, err := config.NewConfigAgent()
	if err != nil {
		return fmt.Errorf("error while load config: %w", err)
	}

	var metricStore = NewMetricStore()

	for i := 1; ; i++ {
		poll(metricStore)
		if i*(conf.PollInterval) >= conf.ReportInterval {
			report(metricStore, conf.HostString)
			clear(metricStore.MetricsGauge)
			clear(metricStore.MetricsCounter)
			i = 0
		}
		time.Sleep(time.Second * time.Duration(conf.PollInterval))
	}
}

func poll(metricStore *MetricStore) {
	ReadRuntimeMetrics(metricStore)

	metricStore.PushCounterMetric("PollCount", storage.CounterValue(1))
	metricStore.PushGaugeMetric("RandomValue", storage.GaugeValue(rand.Float64()*1_000))
}

func report(metricStore *MetricStore, serverURL string) {
	serverURL = "http://" + serverURL

	metricType := storage.CounterType
	for metricName, metric := range metricStore.MetricsCounter {
		value := strconv.FormatInt(int64(metric), 10)

		err := sendMetric(serverURL, metricType, metricName, value)
		if err != nil {
			log.Println(err)
		}
	}

	metricType = storage.GaugeType
	for metricName, metric := range metricStore.MetricsGauge {
		value := strconv.FormatFloat(float64(metric), 'f', 5, 64)

		err := sendMetric(serverURL, metricType, metricName, value)
		if err != nil {
			log.Println(err)
		}
	}
}

func sendMetric(serverURL, metricType, metricName, value string) error {
	requestStr := serverURL + "/update/" + metricType + "/" + metricName + "/" + value

	resp, err := http.Post(requestStr, "text/plain", nil)
	if err != nil {
		// log.Println(requestStr)
		// log.Println(err)
		return fmt.Errorf("error on %s : %w", requestStr, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		// log.Println(requestStr)
		// log.Println("Status from server:", resp.StatusCode)
		return fmt.Errorf("bad response %v for request %s", resp.StatusCode, requestStr)
	}
	return nil
}
