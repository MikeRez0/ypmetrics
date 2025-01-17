package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/MikeRez0/ypmetrics/internal/model"
)

type FileStorage struct {
	MemStorage
	log      *zap.Logger
	filename string
	syncSave bool
}

func NewFileStorage(filename string, saveInterval int, restore bool, log *zap.Logger) (*FileStorage, error) {
	fs := FileStorage{
		MemStorage: *NewMemStorage(),
		filename:   filename,
		syncSave:   saveInterval == 0,
		log:        log,
	}

	if restore {
		err := fs.ReadMetrics(context.Background())
		if err != nil {
			return nil, fmt.Errorf("error restoring from file %s : %w", filename, err)
		}
	}

	if !fs.syncSave {
		ticker := time.NewTicker(time.Duration(saveInterval) * time.Second)
		go func() {
			for range ticker.C {
				err := fs.WriteMetrics()
				if err != nil {
					log.Error("error writing async metrics", zap.Error(err))
				}
			}
		}()
	}

	return &fs, nil
}

func (fs *FileStorage) BatchUpdate(ctx context.Context, metrics []model.Metrics) error {
	err := fs.MemStorage.BatchUpdate(ctx, metrics)
	if err != nil {
		return err
	}

	if fs.syncSave {
		err := fs.WriteMetrics()
		if err != nil {
			return err
		}
	}

	return nil
}

func (fs *FileStorage) UpdateGauge(ctx context.Context,
	metric string, value model.GaugeValue) (model.GaugeValue, error) {
	val, err := fs.MemStorage.UpdateGauge(ctx, metric, value)
	if err != nil {
		return model.GaugeValue(0), err
	}

	if fs.syncSave {
		err = fs.WriteMetrics()
		if err != nil {
			return model.GaugeValue(0), err
		}
	}

	return val, nil
}

func (fs *FileStorage) UpdateCounter(ctx context.Context,
	metric string, value model.CounterValue) (model.CounterValue, error) {
	val, err := fs.MemStorage.UpdateCounter(ctx, metric, value)
	if err != nil {
		return model.CounterValue(0), err
	}
	if fs.syncSave {
		err = fs.WriteMetrics()
		if err != nil {
			return model.CounterValue(0), err
		}
	}

	return val, nil
}

func (fs *FileStorage) WriteMetrics() error {
	fs.log.Info("Start writing metrics to file")
	file, err := os.OpenFile(fs.filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			fs.log.Error("error while closing file", zap.Error(err))
		}
	}()

	encoder := json.NewEncoder(file)

	for _, m := range fs.Metrics() {
		err = encoder.Encode(m)
		if err != nil {
			return fmt.Errorf("error encoding metric %s: %w", m.ID, err)
		}
	}

	fs.log.Info("End writing metrics to file")
	return nil
}

func (fs *FileStorage) ReadMetrics(ctx context.Context) error {
	fs.log.Info("Start reading metrics from file")
	file, err := os.OpenFile(fs.filename, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return fmt.Errorf("error while open file %s: %w", fs.filename, err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			fs.log.Error("Error while closing file", zap.Error(err))
		}
	}()

	scan := bufio.NewScanner(file)

	var metric model.Metrics
	for scan.Scan() {
		data := scan.Bytes()
		err = json.Unmarshal(data, &metric)
		if err != nil {
			return fmt.Errorf("error while read file %s: %w", fs.filename, err)
		}
		err = fs.StoreMetric(ctx, metric)
		if err != nil {
			return fmt.Errorf("error while store metric: %w", err)
		}
	}

	fs.log.Info("End reading metrics from file")
	return nil
}

func (fs *FileStorage) Ping() error {
	return errors.New("Ping not supported")
}
