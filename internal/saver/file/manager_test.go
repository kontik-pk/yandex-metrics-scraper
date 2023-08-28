package file

import (
	"context"
	"os"
	"testing"

	"github.com/kontik-pk/yandex-metrics-scraper/internal/collector"
	"github.com/stretchr/testify/assert"
)

func TestManager_Restore(t *testing.T) {
	testCases := []struct {
		name            string
		fileContent     string
		expectedMetrics []collector.StoredMetric
	}{
		{
			name:        "positive: some metrics",
			fileContent: `[{"id":"FreeMemory","type":"gauge","gauge_value":341491712,"text_value":"341491712.00000000000"},{"id":"TotalMemory","type":"gauge","gauge_value":34359738368,"text_value":"34359738368.00000000000"},{"id":"PollCount","type":"counter","counter_value":100500,"text_value":"100500"}]`,
			expectedMetrics: []collector.StoredMetric{
				{
					ID:         "FreeMemory",
					MType:      "gauge",
					GaugeValue: collector.PtrFloat64(341491712),
					TextValue:  collector.PtrString("341491712.00000000000"),
				},
				{
					ID:         "TotalMemory",
					MType:      "gauge",
					GaugeValue: collector.PtrFloat64(34359738368),
					TextValue:  collector.PtrString("34359738368.00000000000"),
				},
				{
					ID:           "PollCount",
					MType:        "counter",
					CounterValue: collector.PtrInt64(100500),
					TextValue:    collector.PtrString("100500"),
				},
			},
		},
		{
			name:            "positive: some metrics",
			fileContent:     ``,
			expectedMetrics: []collector.StoredMetric(nil),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tdir := t.TempDir()
			file, err := os.CreateTemp(tdir, "")
			assert.NoError(t, err)

			_, err = file.Write([]byte(tt.fileContent))
			assert.NoError(t, err)

			ctx := context.Background()
			manager := New(file.Name())
			metrics, err := manager.Restore(ctx)
			assert.NoError(t, err)
			assert.Equal(t, metrics, tt.expectedMetrics)
		})
	}
}

func TestManager_Save(t *testing.T) {
	testCases := []struct {
		name                string
		metrics             []collector.StoredMetric
		expectedFileContent string
	}{
		{
			name: "positive",
			metrics: []collector.StoredMetric{
				{
					ID:         "FreeMemory",
					MType:      "gauge",
					GaugeValue: collector.PtrFloat64(341491712),
					TextValue:  collector.PtrString("341491712.00000000000"),
				},
				{
					ID:         "TotalMemory",
					MType:      "gauge",
					GaugeValue: collector.PtrFloat64(34359738368),
					TextValue:  collector.PtrString("34359738368.00000000000"),
				},
				{
					ID:           "PollCount",
					MType:        "counter",
					CounterValue: collector.PtrInt64(100500),
					TextValue:    collector.PtrString("100500"),
				},
			},
			expectedFileContent: `[{"id":"FreeMemory","type":"gauge","gauge_value":341491712,"text_value":"341491712.00000000000"},{"id":"TotalMemory","type":"gauge","gauge_value":34359738368,"text_value":"34359738368.00000000000"},{"id":"PollCount","type":"counter","counter_value":100500,"text_value":"100500"}]
`,
		},
		{
			name:                "positive: no metrics",
			metrics:             []collector.StoredMetric{},
			expectedFileContent: "[]\n",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tdir := t.TempDir()
			file, err := os.CreateTemp(tdir, "")
			assert.NoError(t, err)

			ctx := context.Background()
			manager := New(file.Name())
			err = manager.Save(ctx, tt.metrics)
			assert.NoError(t, err)

			b, err := os.ReadFile(file.Name())
			assert.NoError(t, err)
			assert.Equal(t, string(b), tt.expectedFileContent)
		})
	}
}
