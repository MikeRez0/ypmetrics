// Package server - server app for collecting and storing metrics.
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/MikeRez0/ypmetrics/internal/config"
	"github.com/MikeRez0/ypmetrics/internal/handlers"
	"github.com/MikeRez0/ypmetrics/internal/logger"
	"github.com/MikeRez0/ypmetrics/internal/storage"
	"github.com/MikeRez0/ypmetrics/internal/utils/netctrl"
	"github.com/MikeRez0/ypmetrics/internal/utils/signer"
	"go.uber.org/zap"
)

// Run - runs server on config params.
func Run() error {
	conf, err := config.NewConfigServer()
	if err != nil {
		return fmt.Errorf("error while config load: %w", err)
	}

	mylog := logger.GetLogger(conf.LogLevel)
	mylog.Info(fmt.Sprintf("cmd args: %v", os.Args[1:]))
	mylog.Info(fmt.Sprintf("start server with config: %v", conf))

	var repo handlers.Repository

	ctxBackround, cancelBackround := context.WithCancel(context.Background())
	defer cancelBackround()

	wg := &sync.WaitGroup{}

	switch {
	case conf.DSN != "":
		repo, err = storage.NewDBStorage(
			conf.DSN,
			logger.LoggerWithComponent(mylog, "dbstorage"))
		if err != nil {
			return fmt.Errorf("error creating db repo: %w", err)
		}
	case conf.FileStoragePath != "":
		repo, err = storage.NewFileStorage(ctxBackround,
			conf,
			wg,
			logger.LoggerWithComponent(mylog, "filestorage"))
		if err != nil {
			return fmt.Errorf("error creating file repo: %w", err)
		}
	default:
		repo = storage.NewMemStorage()
	}

	h, err := handlers.NewMetricsHandler(repo, logger.LoggerWithComponent(mylog, "handlers"))
	if err != nil {
		return fmt.Errorf("error creating handler: %w", err)
	}

	var netc *netctrl.IPControl
	if conf.TrustedSubnet != "" {
		netc, err = netctrl.NewIPControl(conf.TrustedSubnet, mylog.Named("netcontrol"))
		if err != nil {
			return fmt.Errorf("error creating net control: %w", err)
		}
	}

	r := handlers.SetupRouter(h, logger.LoggerWithComponent(mylog, "handlers"), netc)

	if conf.SignKey != "" {
		h.Signer = signer.NewSigner(conf.SignKey)
	}

	if conf.CryptoKey != "" {
		decrypter, err := signer.NewDecrypter(conf.CryptoKey, mylog.Named("decrypt"))
		if err != nil {
			return fmt.Errorf("error creating decryptor: %w", err)
		}
		h.Decrypter = decrypter
	}

	server := &http.Server{
		Addr:    conf.HostString,
		Handler: r.Handler(),
	}

	shutdown := make(chan os.Signal, 1)
	waitForShutdown := make(chan struct{})
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-shutdown
		mylog.Info("Start graceful shutdown...")

		cancelBackround()

		err := server.Shutdown(context.Background())
		if err != nil {
			mylog.Error("error while shutdown", zap.Error(err))
		}
		wg.Wait()
		waitForShutdown <- struct{}{}
	}()

	err = server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error while run server: %w", err)
	}
	<-waitForShutdown
	fmt.Println("Server was shut down gracefully")
	return nil
}
