package logger

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	once   sync.Once
)

func GetLogger(level string) *zap.Logger {
	once.Do(func() {
		var err error
		logger, err = Initialize(level)
		if err != nil {
			panic(fmt.Sprintf("error init logger: %v", err))
		}
	})

	return logger
}

func Initialize(level string) (*zap.Logger, error) {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the atomic level %s: %w", level, err)
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	// устанавливаем уровень
	cfg.Level = lvl
	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("bad logger config: %w", err)
	}
	// устанавливаем синглтон
	return zl, nil
}

func GinLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		log.Info("Incoming HTTP request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.RequestURI),
			zap.Int("status", c.Writer.Status()),
			zap.Int("size", c.Writer.Size()),
			zap.String("duration", time.Since(t).String()))
	}
}
func LoggerWithComponent(logger *zap.Logger, name string) *zap.Logger {
	return logger.With(zap.String("component", name))
}
