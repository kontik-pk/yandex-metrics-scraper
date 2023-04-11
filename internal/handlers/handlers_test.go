package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/collector"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSaveMetric(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", SaveMetric)
	srv := httptest.NewServer(r)
	defer srv.Close()

	testCases := []struct {
		name          string
		mType         string
		mName         string
		mValue        string
		expectedCode  int
		expectedError error
	}{
		{
			name:         "case0",
			mType:        "counter",
			mName:        "Counter1",
			mValue:       "15",
			expectedCode: http.StatusOK,
		},
		{
			name:         "case1",
			mType:        "gauge",
			mName:        "Gauge1",
			mValue:       "12.282",
			expectedCode: http.StatusOK,
		},
		{
			name:          "case2",
			mType:         "invalid",
			mName:         "Gauge1",
			mValue:        "12.282",
			expectedCode:  http.StatusNotImplemented,
			expectedError: collector.ErrNotImplemented,
		},
		{
			name:          "case3",
			mType:         "counter",
			mName:         "Counter1",
			mValue:        "15.2562",
			expectedCode:  http.StatusBadRequest,
			expectedError: collector.ErrNotFound,
		},
		{
			name:          "case4",
			mType:         "gauge",
			mName:         "Gauge1",
			mValue:        "12.282dgh",
			expectedCode:  http.StatusBadRequest,
			expectedError: collector.ErrNotFound,
		},
		{
			name:          "case5",
			mType:         "gauge",
			mName:         "Gauge1",
			mValue:        "",
			expectedCode:  http.StatusNotFound,
			expectedError: collector.ErrNotFound,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := resty.New().R().
				SetHeader("Content-Type", "text/plain").
				Post(fmt.Sprintf("%s/update/%s/%s/%s", srv.URL, tt.mType, tt.mName, tt.mValue))

			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, resp.StatusCode(), tt.expectedCode)

			value, err := collector.Collector.GetMetric(tt.mName, tt.mType)
			if err != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			if tt.expectedCode == http.StatusOK {
				assert.Equal(t, value, tt.mValue)
			}
		})
	}
}
