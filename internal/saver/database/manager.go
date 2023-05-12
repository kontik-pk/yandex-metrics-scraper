package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/collector"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/flags"
	"time"
)

func (m *manager) Restore(ctx context.Context) ([]collector.StoredMetric, error) {
	const query = `select id, mtype, delta, mvalue from metrics`
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	var metrics []collector.StoredMetric
	for rows.Next() {
		var (
			id          string
			mtype       string
			deltaFromDB sql.NullInt64
			valueFromDB sql.NullFloat64
		)
		if err := rows.Scan(&id, &mtype, &deltaFromDB, &valueFromDB); err != nil {
			return nil, err
		}
		var delta *int64
		if deltaFromDB.Valid {
			delta = &deltaFromDB.Int64
		}
		var mvalue *float64
		if valueFromDB.Valid {
			mvalue = &valueFromDB.Float64
		}
		metric := collector.StoredMetric{
			ID:           id,
			MType:        mtype,
			CounterValue: delta,
			GaugeValue:   mvalue,
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func (m *manager) Save(ctx context.Context, metrics []collector.StoredMetric) error {
	retries := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
	for _, metric := range metrics {
		switch metric.MType {
		case collector.Gauge:
			query := `insert into metrics (id, mtype, mvalue) values ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET mvalue = EXCLUDED.mvalue;`
			var err error
			for _, t := range retries {
				if _, err = m.db.ExecContext(ctx, query, metric.ID, metric.MType, &metric.GaugeValue); err == nil {
					break
				} else {
					time.Sleep(t)
				}
			}
			if err != nil {
				return fmt.Errorf("error while executing insert counter query: %w", err)
			}
		case collector.Counter:
			query := `insert into metrics (id, mtype, delta) values ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET delta = EXCLUDED.delta;`
			//TODO: так себе ретраи, ретраить нужно не любые ошибки, ну и реализация фу
			var err error
			for _, t := range retries {
				if _, err = m.db.ExecContext(ctx, query, metric.ID, metric.MType, &metric.CounterValue); err == nil {
					break
				} else {
					time.Sleep(t)
				}
			}
			if err != nil {
				return fmt.Errorf("error while executing insert counter query: %w", err)
			}
		}
	}
	return nil
}

func (m *manager) init(ctx context.Context) error {
	const query = `create table if not exists metrics (id text primary key, mtype text, delta bigint, mvalue double precision)`
	if _, err := m.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("error while trying to create table: %w", err)
	}
	return nil
}

func New(params *flags.Params) (*manager, error) {
	ctx := context.Background()
	db, err := sql.Open("pgx", params.DatabaseAddress)
	if err != nil {
		return nil, err
	}

	m := manager{
		db: db,
	}
	if err := m.init(ctx); err != nil {
		return nil, err
	}
	return &m, nil
}

type manager struct {
	db *sql.DB
}
