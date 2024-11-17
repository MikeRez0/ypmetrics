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
	"github.com/MikeRez0/ypmetrics/internal/utils/retrier"
	"go.uber.org/zap"
)

type DBStorage struct {
	log  *zap.Logger
	pool *pgxpool.Pool
}

//go:embed migrations/*.sql
var migrationsDir embed.FS

func NewDBStorage(dsn string, log *zap.Logger) (*DBStorage, error) {
	if err := runMigrations(dsn); err != nil {
		return nil, fmt.Errorf("failed to run DB migrations: %w", err)
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create a connection pool: %w", err)
	}

	dbs := DBStorage{
		pool: pool,
		log:  log,
	}

	log.Debug("Success connected to db")

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
	err := retrier.Retry(ctx, func() error {
		tx, err := ds.pool.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			return fmt.Errorf("error starting transaction: %w", err)
		}

		defer func() {
			err := tx.Rollback(ctx)
			if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
				ds.log.Error("error while rollback", zap.Error(err))
			}
		}()

		mt := model.MetricType(model.GaugeType)

		_, err = tx.Exec(ctx,
			`INSERT INTO "metric" ("id", "mtype", "value", "updts")
			VALUES ($1, $2, $3, $4);`,
			metric, mt, value, time.Now())
		if err != nil {
			return fmt.Errorf("error inserting metric: %w", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			return fmt.Errorf("error commiting transaction: %w", err)
		}
		return nil
	}, 3, ds.log)

	if err != nil {
		return 0, err //nolint:wrapcheck // callback error
	}

	return value, nil
}

func (ds *DBStorage) UpdateCounter(ctx context.Context,
	metric string, value model.CounterValue) (model.CounterValue, error) {
	newVal := value

	err := retrier.Retry(ctx, func() error {
		tx, err := ds.pool.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			return fmt.Errorf("error starting transaction: %w", err)
		}

		defer func() {
			err := tx.Rollback(ctx)
			if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
				ds.log.Error("error while rollback", zap.Error(err))
			}
		}()

		mt := model.MetricType(model.GaugeType)

		row := tx.QueryRow(ctx,
			`INSERT INTO metric
				(id, mtype, delta, updts)
				values ($1, $2, $3, $4)
				ON CONFLICT (id) DO UPDATE 
				SET delta= metric.delta + EXCLUDED.delta
				RETURNING delta;`,
			metric, mt, value, time.Now())
		err = row.Scan(&newVal)
		if err != nil {
			return fmt.Errorf("error inserting metric: %w", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			return fmt.Errorf("error commiting transaction: %w", err)
		}
		return nil
	}, 3, ds.log)

	if err != nil {
		return 0, err //nolint:wrapcheck // callback error
	}

	return newVal, nil
}

func (ds *DBStorage) readMetrics(ctx context.Context, id string) ([]model.Metrics, error) {
	ds.log.Debug("Start reading metrics from database")

	metricsList := make([]model.Metrics, 0)

	var (
		rows pgx.Rows
		err  error
	)

	if id == "" {
		rows, err = ds.pool.Query(ctx,
			`SELECT "id", "mtype", "delta", "value"
			FROM "metric"`)
	} else {
		rows, err = ds.pool.Query(ctx,
			`SELECT "id", "mtype", "delta", "value"
			FROM "metric" where "id" = $1`, id)
	}
	if err != nil {
		return nil, fmt.Errorf("error selecting metric: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var metric model.Metrics

		err = rows.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)
		if err != nil {
			return nil, fmt.Errorf("error reading metric: %w", err)
		}
		metricsList = append(metricsList, metric)
	}

	ds.log.Debug("End reading metrics from database")
	return metricsList, nil
}

func (ds *DBStorage) Ping() error {
	err := retrier.Retry(context.Background(),
		func() error {
			return ds.pool.Ping(context.Background())
		}, 3, ds.log)
	if err != nil {
		return fmt.Errorf("error connecting DB: %w", err)
	}
	return nil
}

func (ds *DBStorage) GetCounter(ctx context.Context, metric string) (model.CounterValue, error) {
	ms, err := ds.readMetrics(ctx, metric)
	if err != nil {
		return 0, err
	}
	if len(ms) != 1 {
		return 0, fmt.Errorf("metric %s not found", metric)
	}
	return model.CounterValue(*ms[0].Delta), nil
}

func (ds *DBStorage) GetGauge(ctx context.Context, metric string) (model.GaugeValue, error) {
	ms, err := ds.readMetrics(ctx, metric)
	if err != nil {
		return 0, err
	}
	if len(ms) != 1 {
		return 0, fmt.Errorf("metric %s not found", metric)
	}
	return model.GaugeValue(*ms[0].Value), nil
}

func (ds *DBStorage) Metrics() (res []model.Metrics) {
	res, err := ds.readMetrics(context.Background(), "")

	if err != nil {
		return nil
	}

	return res
}

func (ds *DBStorage) BatchUpdate(ctx context.Context, metrics []model.Metrics) error {
	ds.log.Info("Start writing Batch metrics to database")

	err := retrier.Retry(ctx, func() error {
		tx, err := ds.pool.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			return fmt.Errorf("error starting transaction: %w", err)
		}
		defer func() {
			err := tx.Rollback(ctx)
			if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
				ds.log.Error("error while rollback", zap.Error(err))
			}
		}()

		for _, m := range metrics {
			mt, _ := m.MType.Value()

			statement := `INSERT INTO "metric" ("id", "mtype", "delta", "value", "updts")
			VALUES ($1, $2, $3, $4, $5)`

			switch m.MType {
			case model.GaugeType:
				statement += `ON CONFLICT ("id") DO UPDATE
				SET "mtype" = $2, "delta" = $3, "value" = $4, "updts" = $5;`
			case model.CounterType:
				statement += `ON CONFLICT ("id") DO UPDATE
				SET "mtype" = $2, "delta" = metric.delta + EXCLUDED.delta, "value" = $4, "updts" = $5;`
			}

			_, err := tx.Exec(ctx, statement,
				m.ID, mt, m.Delta, m.Value, time.Now())
			if err != nil {
				return fmt.Errorf("error upserting metric: %w", err)
			}
		}

		err = tx.Commit(ctx)
		if err != nil {
			return fmt.Errorf("error commiting transaction: %w", err)
		}
		return nil
	}, 3, ds.log)
	if err != nil {
		return err //nolint:wrapcheck //error from callback
	}

	ds.log.Info("End writing Batch metrics to database")
	return nil
}
