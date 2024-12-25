package retrier

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type Retrier struct {
	logger       *zap.Logger
	attempts     int
	intervalStep int
}

type RetryFunc func() error

func NewRetrier(log *zap.Logger, attempts int, intervalStep int) *Retrier {
	return &Retrier{
		logger:       log,
		attempts:     attempts,
		intervalStep: intervalStep,
	}
}
func (r *Retrier) Retry(ctx context.Context,
	f RetryFunc,
	clauseCanRetry func(error) bool) error {
	var err error
	for i := range r.attempts {
		select {
		case <-ctx.Done():
			r.logger.Error("All retries failed :( ")
			return ctx.Err() //nolint:wrapcheck //nothing to wrap
		default:
		}
		err = f()
		if err == nil {
			return nil
		}
		if !clauseCanRetry(err) {
			return err
		}
		r.logger.Info(fmt.Sprintf("Going to retry # %v with error", i+1), zap.Error(err))

		<-time.After(time.Duration(r.intervalStep*(i+1)) * time.Second)
	}
	return err
}
