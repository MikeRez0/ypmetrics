package utils

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type RetryFunc func() error

func Retry(ctx context.Context, f RetryFunc, retryAttempts int, logger *zap.Logger) error {
	const retryIntervalStep = 2

	var err error
	for i := 0; i < retryAttempts; i++ {
		select {
		case <-ctx.Done():
			logger.Error("All retries failed :( ")
			return ctx.Err() //nolint:wrapcheck //nothing to wrap
		default:
		}
		err := f()
		if err == nil {
			return nil
		}
		logger.Info(fmt.Sprintf("Going to retry # %v with error", i+1), zap.Error(err))
		time.Sleep(time.Duration(1+retryIntervalStep*i) * time.Second)
	}
	return err
}
