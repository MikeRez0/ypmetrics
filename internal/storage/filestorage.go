package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/MikeRez0/ypmetrics/internal/logger"
	"github.com/MikeRez0/ypmetrics/internal/model"
	"go.uber.org/zap"
)

type FileStorage struct {
	MemStorage
	filename string
	syncSave bool
}

func NewFileStorage(filename string, saveInterval int, restore bool) (*FileStorage, error) {
	fs := FileStorage{
		MemStorage: *NewMemStorage(),
		filename:   filename,
		syncSave:   saveInterval == 0,
	}

	if restore {
		err := fs.ReadMetrics()
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
					logger.Log.Error("error writing async metrics", zap.Error(err))
				}
			}
		}()
	}

	return &fs, nil
}

func (fs *FileStorage) UpdateGauge(metric string, value model.GaugeValue) (model.GaugeValue, error) {
	val, err := fs.MemStorage.UpdateGauge(metric, value)
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

func (fs *FileStorage) UpdateCounter(metric string, value model.CounterValue) (model.CounterValue, error) {
	val, err := fs.MemStorage.UpdateCounter(metric, value)
	if err != nil {
		return model.CounterValue(0), err
	}

	return val, nil
}

func (fs *FileStorage) WriteMetrics() error {
	logger.Log.Info("Start writing metrics to file")
	file, err := os.OpenFile(fs.filename, os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			logger.Log.Error("error while closing file", zap.Error(err))
		}
	}()

	encoder := json.NewEncoder(file)

	for _, m := range fs.Metrics() {
		err = encoder.Encode(m)
		if err != nil {
			return fmt.Errorf("error encoding metric %s: %w", m.ID, err)
		}
	}

	logger.Log.Info("End writing metrics to file")
	return nil
}

func (fs *FileStorage) ReadMetrics() error {
	logger.Log.Info("Start reading metrics from file")
	file, err := os.OpenFile(fs.filename, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return fmt.Errorf("error while open file %s: %w", fs.filename, err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			logger.Log.Error("Error while closing file", zap.Error(err))
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
		err = fs.StoreMetric(metric)
		if err != nil {
			return fmt.Errorf("error while store metric: %w", err)
		}
	}

	logger.Log.Info("End reading metrics from file")
	return nil
}
