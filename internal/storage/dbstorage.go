package storage

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MikeRez0/ypmetrics/internal/model"
	"go.uber.org/zap"
)

type DBStorage struct {
	MemStorage
	log      *zap.Logger
	pool     *pgxpool.Pool
	syncSave bool
}

//go:embed migrations/*.sql
var migrationsDir embed.FS

func NewDBStorage(dsn string, saveInterval int, restore bool, log *zap.Logger) (*DBStorage, error) {
	if err := runMigrations(dsn); err != nil {
		return nil, fmt.Errorf("failed to run DB migrations: %w", err)
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create a connection pool: %w", err)
	}

	dbs := DBStorage{
		MemStorage: *NewMemStorage(),
		pool:       pool,
		syncSave:   saveInterval == 0,
		log:        log,
	}

	log.Debug("Success connected to db")

	if restore {
		err := dbs.ReadMetrics(context.Background())
		if err != nil {
			return nil, fmt.Errorf("error restoring from db : %w", err)
		}
	}

	if !dbs.syncSave {
		ticker := time.NewTicker(time.Duration(saveInterval) * time.Second)
		go func() {
			for range ticker.C {
				err := dbs.WriteMetrics(context.Background())
				if err != nil {
					log.Error("error writing async metrics", zap.Error(err))
				}
			}
		}()
	}

	return &dbs, nil
}

func runMigrations(dsn string) error {
	d, err := iofs.New(migrationsDir, "migrations")
	if err != nil {
		return fmt.Errorf("failed to return an iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return fmt.Errorf("failed to get a new migrate instance: %w", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations to the DB: %w", err)
		}
	}
	return nil
}

func (ds *DBStorage) UpdateGauge(ctx context.Context,
	metric string, value model.GaugeValue) (model.GaugeValue, error) {
	val, err := ds.MemStorage.UpdateGauge(ctx, metric, value)
	if err != nil {
		return model.GaugeValue(0), err
	}

	if ds.syncSave {
		err = ds.WriteMetrics(ctx)
		if err != nil {
			return model.GaugeValue(0), err
		}
	}

	return val, nil
}

func (ds *DBStorage) UpdateCounter(ctx context.Context,
	metric string, value model.CounterValue) (model.CounterValue, error) {
	val, err := ds.MemStorage.UpdateCounter(ctx, metric, value)
	if err != nil {
		return model.CounterValue(0), err
	}

	if ds.syncSave {
		err = ds.WriteMetrics(ctx)
		if err != nil {
			return model.CounterValue(0), err
		}
	}

	return val, nil
}

func (ds *DBStorage) WriteMetrics(ctx context.Context) error {
	ds.log.Info("Start writing metrics to database")

	tx, err := ds.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			ds.log.Error("error while rollback", zap.Error(err))
		}
	}()

	_, err = tx.Exec(ctx, `TRUNCATE "metric"`)
	if err != nil {
		return fmt.Errorf("error truncating metrics: %w", err)
	}

	for _, m := range ds.Metrics() {
		mt, _ := m.MType.Value()
		_, err := tx.Exec(ctx,
			`INSERT INTO "metric" ("id", "mtype", "delta", "value")
			VALUES ($1, $2, $3, $4);`,
			m.ID, mt, m.Delta, m.Value)
		if err != nil {
			return fmt.Errorf("error inserting metric: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error commiting transaction: %w", err)
	}

	ds.log.Info("End writing metrics to database")
	return nil
}

func (ds *DBStorage) ReadMetrics(ctx context.Context) error {
	ds.log.Info("Start reading metrics from database")

	rows, err := ds.pool.Query(ctx,
		`SELECT "id", "mtype", "delta", "value"
		FROM "metric"`)
	if err != nil {
		return fmt.Errorf("error selecting metric: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var metric model.Metrics

		err = rows.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)
		if err != nil {
			return fmt.Errorf("error reading metric: %w", err)
		}
		err = ds.StoreMetric(ctx, metric)
		if err != nil {
			return fmt.Errorf("error while store metric: %w", err)
		}
	}

	ds.log.Info("End reading metrics from database")
	return nil
}

func (ds *DBStorage) Ping() error {
	err := ds.pool.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("error connecting DB: %w", err)
	}
	return nil
}
