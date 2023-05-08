package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/collector"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/flags"
)

func (m *manager) Restore(ctx context.Context) ([]collector.MetricJSON, error) {
	const query = `select id, mtype, delta, mvalue from metrics`
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []collector.MetricJSON
	for rows.Next() {
		var (
			id     string
			mtype  string
			delta  int64
			mvalue float64
		)
		if err := rows.Scan(&id, &mtype, &delta, &mvalue); err != nil {
			return nil, err
		}
		metric := collector.MetricJSON{
			ID:    id,
			MType: mtype,
			Delta: &delta,
			Value: &mvalue,
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func (m *manager) Save(ctx context.Context, metrics []collector.MetricJSON) error {
	for _, metric := range metrics {
		switch metric.MType {
		case "gauge":
			query := `insert into metrics (id, mtype, mvalue) values ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET mvalue = EXCLUDED.mvalue;`
			if _, err := m.db.ExecContext(ctx, query, metric.ID, metric.MType, metric.Value); err != nil {
				return fmt.Errorf("error while trying to save gauge metric %q: %w", metric.ID, err)
			}
		case "counter":
			query := `insert into metrics (id, mtype, delta) values ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET delta = EXCLUDED.delta;`
			if _, err := m.db.ExecContext(ctx, query, metric.ID, metric.MType, metric.Delta); err != nil {
				return fmt.Errorf("error while trying to save counter metric %q: %w", metric.ID, err)
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
	defer db.Close()

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
