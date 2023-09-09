package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/agent/collector"
	"time"
)

// Restore - a method for restoring metrics state from DB.
func (m *Manager) Restore(ctx context.Context) ([]collector.StoredMetric, error) {
	const query = `select id, mtype, delta, mvalue from metrics`
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if err = rows.Err(); err != nil {
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
		if err = rows.Scan(&id, &mtype, &deltaFromDB, &valueFromDB); err != nil {
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

// Save - a method for saving metric state to the DB.
func (m *Manager) Save(ctx context.Context, metrics []collector.StoredMetric) error {
	for _, metric := range metrics {
		switch metric.MType {
		case collector.Gauge:
			query := `insert into metrics (id, mtype, mvalue) values ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET mvalue = EXCLUDED.mvalue;`
			if err := m.execWithRetries(ctx, query, metric.ID, metric.MType, &metric.GaugeValue); err != nil {
				return fmt.Errorf("error while executing insert counter query: %w", err)
			}
		case collector.Counter:
			query := `insert into metrics (id, mtype, delta) values ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET delta = EXCLUDED.delta;`
			if err := m.execWithRetries(ctx, query, metric.ID, metric.MType, &metric.CounterValue); err != nil {
				return fmt.Errorf("error while executing insert counter query: %w", err)
			}
		}
	}
	return nil
}

func (m *Manager) execWithRetries(ctx context.Context, query string, args ...any) error {
	var err error
	retries := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
	for _, t := range retries {
		if _, err = m.db.ExecContext(ctx, query, args...); err == nil {
			break
		} else {
			time.Sleep(t)
		}
	}
	return err
}

// init - a method for preparing new DB for storing metrics.
func (m *Manager) init(ctx context.Context) error {
	const query = `create table if not exists metrics (id text primary key, mtype text, delta bigint, mvalue double precision)`
	if _, err := m.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("error while trying to create table: %w", err)
	}
	return nil
}

func New(db *sql.DB) (*Manager, error) {
	ctx := context.Background()
	m := Manager{
		db: db,
	}
	if err := m.init(ctx); err != nil {
		return nil, err
	}
	return &m, nil
}

type Manager struct {
	db *sql.DB
}
