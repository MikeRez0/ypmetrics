package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/MikeRez0/ypmetrics/internal/model"
	"go.uber.org/zap"
)

type DBStorage struct {
	MemStorage
	log      *zap.Logger
	db       *sql.DB
	syncSave bool
}

func NewDBStorage(dsn string, log *zap.Logger) (*DBStorage, error) {
	fs := DBStorage{
		MemStorage: *NewMemStorage(),
		db:         nil,
		log:        log,
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening DB connection: %w", err)
	}
	fs.db = db
	// defer func() {
	// 	err := db.Close()
	// 	if err != nil {
	// 		log.Error("db close error", zap.Error(err))
	// 	}
	// }()

	log.Debug("Success connected to db")

	return &fs, nil
}

func (ds *DBStorage) UpdateGauge(metric string, value model.GaugeValue) (model.GaugeValue, error) {
	val, err := ds.MemStorage.UpdateGauge(metric, value)
	if err != nil {
		return model.GaugeValue(0), err
	}

	if ds.syncSave {
		err = ds.WriteMetrics()
		if err != nil {
			return model.GaugeValue(0), err
		}
	}

	return val, nil
}

func (ds *DBStorage) UpdateCounter(metric string, value model.CounterValue) (model.CounterValue, error) {
	val, err := ds.MemStorage.UpdateCounter(metric, value)
	if err != nil {
		return model.CounterValue(0), err
	}

	if ds.syncSave {
		err = ds.WriteMetrics()
		if err != nil {
			return model.CounterValue(0), err
		}
	}

	return val, nil
}

func (ds *DBStorage) WriteMetrics() error {
	// ds.log.Info("Start writing metrics to file")
	// file, err := os.OpenFile(ds.filename, os.O_CREATE|os.O_WRONLY, 0o600)
	// if err != nil {
	// 	return fmt.Errorf("error opening file: %w", err)
	// }
	// defer func() {
	// 	err := file.Close()
	// 	if err != nil {
	// 		ds.log.Error("error while closing file", zap.Error(err))
	// 	}
	// }()

	// encoder := json.NewEncoder(file)

	// for _, m := range ds.Metrics() {
	// 	err = encoder.Encode(m)
	// 	if err != nil {
	// 		return fmt.Errorf("error encoding metric %s: %w", m.ID, err)
	// 	}
	// }

	// ds.log.Info("End writing metrics to file")
	return nil
}

func (ds *DBStorage) ReadMetrics() error {
	// ds.log.Info("Start reading metrics from file")
	// file, err := os.OpenFile(ds.filename, os.O_CREATE|os.O_RDWR, 0o600)
	// if err != nil {
	// 	return fmt.Errorf("error while open file %s: %w", ds.filename, err)
	// }
	// defer func() {
	// 	err := file.Close()
	// 	if err != nil {
	// 		ds.log.Error("Error while closing file", zap.Error(err))
	// 	}
	// }()

	// scan := bufio.NewScanner(file)

	// var metric model.Metrics
	// for scan.Scan() {
	// 	data := scan.Bytes()
	// 	err = json.Unmarshal(data, &metric)
	// 	if err != nil {
	// 		return fmt.Errorf("error while read file %s: %w", ds.filename, err)
	// 	}
	// 	err = ds.StoreMetric(metric)
	// 	if err != nil {
	// 		return fmt.Errorf("error while store metric: %w", err)
	// 	}
	// }

	// ds.log.Info("End reading metrics from file")
	return nil
}

func (ds *DBStorage) Ping() error {
	err := ds.db.Ping()
	if err != nil {
		return fmt.Errorf("error connecting DB: %w", err)
	}
	return nil
}
