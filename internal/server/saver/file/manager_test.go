package file

import (
	"context"
	collector2 "github.com/kontik-pk/yandex-metrics-scraper/internal/agent/collector"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManager_Restore(t *testing.T) {
	testCases := []struct {
		name            string
		fileContent     string
		expectedMetrics []collector2.StoredMetric
	}{
		{
			name:        "positive: some metrics",
			fileContent: `[{"id":"FreeMemory","type":"gauge","gauge_value":341491712,"text_value":"341491712.00000000000"},{"id":"TotalMemory","type":"gauge","gauge_value":34359738368,"text_value":"34359738368.00000000000"},{"id":"PollCount","type":"counter","counter_value":100500,"text_value":"100500"}]`,
			expectedMetrics: []collector2.StoredMetric{
				{
					ID:         "FreeMemory",
					MType:      "gauge",
					GaugeValue: collector2.PtrFloat64(341491712),
					TextValue:  collector2.PtrString("341491712.00000000000"),
				},
				{
					ID:         "TotalMemory",
					MType:      "gauge",
					GaugeValue: collector2.PtrFloat64(34359738368),
					TextValue:  collector2.PtrString("34359738368.00000000000"),
				},
				{
					ID:           "PollCount",
					MType:        "counter",
					CounterValue: collector2.PtrInt64(100500),
					TextValue:    collector2.PtrString("100500"),
				},
			},
		},
		{
			name:            "positive: some metrics",
			fileContent:     ``,
			expectedMetrics: []collector2.StoredMetric(nil),
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
		metrics             []collector2.StoredMetric
		expectedFileContent string
	}{
		{
			name: "positive",
			metrics: []collector2.StoredMetric{
				{
					ID:         "FreeMemory",
					MType:      "gauge",
					GaugeValue: collector2.PtrFloat64(341491712),
					TextValue:  collector2.PtrString("341491712.00000000000"),
				},
				{
					ID:         "TotalMemory",
					MType:      "gauge",
					GaugeValue: collector2.PtrFloat64(34359738368),
					TextValue:  collector2.PtrString("34359738368.00000000000"),
				},
				{
					ID:           "PollCount",
					MType:        "counter",
					CounterValue: collector2.PtrInt64(100500),
					TextValue:    collector2.PtrString("100500"),
				},
			},
			expectedFileContent: `[{"id":"FreeMemory","type":"gauge","gauge_value":341491712,"text_value":"341491712.00000000000"},{"id":"TotalMemory","type":"gauge","gauge_value":34359738368,"text_value":"34359738368.00000000000"},{"id":"PollCount","type":"counter","counter_value":100500,"text_value":"100500"}]
`,
		},
		{
			name:                "positive: no metrics",
			metrics:             []collector2.StoredMetric{},
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
