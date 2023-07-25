package database

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/collector"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestManager_Restore(t *testing.T) {
	testCases := []struct {
		name     string
		rows     *sqlmock.Rows
		expected []collector.StoredMetric
	}{
		{
			name:     "positive: no saved metrics",
			rows:     sqlmock.NewRows([]string{"id", "mtype", "delta", "mvalue"}),
			expected: []collector.StoredMetric(nil),
		},
		{
			name: "positive: one saved metric",
			rows: sqlmock.NewRows([]string{"id", "mtype", "delta", "mvalue"}).AddRow("metricName", "counter", 5, nil),
			expected: []collector.StoredMetric{
				{
					ID:           "metricName",
					MType:        "counter",
					CounterValue: collector.PtrInt64(5),
				},
			},
		},
		{
			name: "positive: some saved metrics",
			rows: sqlmock.NewRows([]string{"id", "mtype", "delta", "mvalue"}).
				AddRow("metricName", "counter", 5, nil).
				AddRow("otherMetricName", "gauge", nil, 10.502),
			expected: []collector.StoredMetric{
				{
					ID:           "metricName",
					MType:        "counter",
					CounterValue: collector.PtrInt64(5),
				},
				{
					ID:         "otherMetricName",
					MType:      "gauge",
					GaugeValue: collector.PtrFloat64(10.502),
				},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			mock.ExpectExec("create table if not exists metrics").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectQuery("select id, mtype, delta, mvalue from metrics").WillReturnRows(tt.rows)
			manager, err := New(db)
			assert.NoError(t, err)

			ctx := context.Background()
			metrics, err := manager.Restore(ctx)
			assert.NoError(t, err)
			assert.Equal(t, metrics, tt.expected)
		})
	}
}

func TestManager_Save(t *testing.T) {
	testCases := []struct {
		name    string
		metrics []collector.StoredMetric
	}{
		{
			name: "positive: store gauge",
			metrics: []collector.StoredMetric{
				{
					ID:         "metricName",
					MType:      "gauge",
					GaugeValue: collector.PtrFloat64(10.502),
				},
				{
					ID:           "otherMetricName",
					MType:        "counter",
					CounterValue: collector.PtrInt64(10),
				},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			mock.ExpectExec("create table if not exists metrics").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("insert into metrics").WithArgs(tt.metrics[0].ID, tt.metrics[0].MType, &tt.metrics[0].GaugeValue).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("insert into metrics").WithArgs(tt.metrics[1].ID, tt.metrics[1].MType, &tt.metrics[1].CounterValue).WillReturnResult(sqlmock.NewResult(1, 1))
			manager, err := New(db)
			assert.NoError(t, err)

			ctx := context.Background()
			err = manager.Save(ctx, tt.metrics)
			assert.NoError(t, err)
		})
	}
}
