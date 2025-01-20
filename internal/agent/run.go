package agent

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/MikeRez0/ypmetrics/internal/config"
	"github.com/MikeRez0/ypmetrics/internal/logger"
)

// Run - run agent app with config.
func Run() error {
	conf, err := config.NewConfigAgent()
	if err != nil {
		return fmt.Errorf("error while load config: %w", err)
	}

	log := logger.GetLogger(conf.LogLevel)
	log.Info(fmt.Sprintf("cmd args: %v", os.Args[1:]))
	log.Info(fmt.Sprintf("start agent with config: %v", conf))

	app, err := NewAgentApp(conf, log)
	if err != nil {
		return fmt.Errorf("error creating app: %w", err)
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	var wg sync.WaitGroup

	jobStart(ctx, &wg, func() error {
		app.Poll()
		return nil
	}, conf.PollInterval.Duration, 1, log.Named("Poll go metrics job"))

	jobStart(ctx, &wg, func() error {
		app.ReadGopsutilMetrics()
		return nil
	}, conf.PollInterval.Duration, 1, log.Named("Poll gopsutil metrics job"))

	jobStart(ctx, &wg, func() error {
		app.ReportBatch()
		return nil
	}, conf.ReportInterval.Duration, conf.RateLimit, log.Named("Report metrics job"))

	wg.Wait()
	log.Info("Application graceful shutdowned")
	return nil
}

// ticker with worker pool.
func jobStart(ctx context.Context, wg *sync.WaitGroup,
	job func() error, interval time.Duration, workers int,
	log *zap.Logger) {
	jobFire := make(chan struct{})

	for i := 0; i < workers; i++ {
		i := i
		wg.Add(1)
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
					wg.Done()
					log.Debug(fmt.Sprintf("Worker %d stopped", i))
					return
				}
			}
		}(jobFire)
	}

	ticker := time.NewTicker(interval)
	wg.Add(1)
	go func() {
		for {
			select {
			case <-ticker.C:
				jobFire <- struct{}{}
			case <-ctx.Done():
				wg.Done()
				ticker.Stop()
				log.Debug("Ticker stopped")
				return
			}
		}
	}()
}
