package handlers

import (
	"encoding/json"
	"fmt"
	collector2 "github.com/kontik-pk/yandex-metrics-scraper/internal/agent/collector"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestHandler_SaveListMetricsFromJSON(t *testing.T) {
	r := chi.NewRouter()
	h := handler{}
	r.Post("/updates/", h.SaveListMetricsFromJSON)
	srv := httptest.NewServer(r)
	defer srv.Close()

	testCases := []struct {
		name           string
		request        []collector2.MetricRequest
		expectedMetric []collector2.StoredMetric
		expectedCode   int
		expectedError  error
	}{
		{
			name: "positive",
			request: []collector2.MetricRequest{
				{
					MType: "counter",
					ID:    "Counter20",
					Delta: collector2.PtrInt64(20),
				},
				{
					MType: "gauge",
					ID:    "Gauge13",
					Value: collector2.PtrFloat64(13.1),
				},
			},
			expectedMetric: []collector2.StoredMetric{
				{
					MType:        "counter",
					ID:           "Counter20",
					CounterValue: collector2.PtrInt64(20),
					TextValue:    collector2.PtrString("20"),
				},
				{
					MType:      "gauge",
					ID:         "Gauge13",
					GaugeValue: collector2.PtrFloat64(13.1),
					TextValue:  collector2.PtrString("13.10000000000"),
				},
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "negative: unsupported metric type",
			request: []collector2.MetricRequest{
				{
					MType: "counter",
					ID:    "Counter20",
					Delta: collector2.PtrInt64(20),
				},
				{
					MType: "undefined",
					ID:    "Gauge13",
					Value: collector2.PtrFloat64(13.1),
				},
			},
			expectedMetric: []collector2.StoredMetric{},
			expectedCode:   http.StatusNotImplemented,
		},
		{
			name: "negative: invalid value",
			request: []collector2.MetricRequest{
				{
					MType: "counter",
					ID:    "Counter20",
					Delta: collector2.PtrInt64(-20),
				},
				{
					MType: "undefined",
					ID:    "Gauge13",
					Value: collector2.PtrFloat64(13.1),
				},
			},
			expectedMetric: []collector2.StoredMetric{},
			expectedCode:   http.StatusBadRequest,
		},
	}
	for _, tt := range testCases {

		t.Run(tt.name, func(t *testing.T) {
			resBody, err := json.Marshal(tt.request)
			assert.NoError(t, err)
			resp, err := resty.New().R().
				SetHeader("Content-Type", "text/plain").
				SetBody(resBody).
				Post(fmt.Sprintf("%s/updates/", srv.URL))

			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, resp.StatusCode(), tt.expectedCode)
			if resp.StatusCode() != http.StatusOK {
				return
			}
			for i, m := range tt.request {
				value, err := collector2.Collector().GetMetricJSON(m.ID)
				if err != nil {
					assert.EqualError(t, err, tt.expectedError.Error())
				} else {
					assert.NoError(t, err)
				}
				actual := collector2.StoredMetric{}
				json.Unmarshal(value, &actual)

				if tt.expectedCode == http.StatusOK {
					assert.Equal(t, actual, tt.expectedMetric[i])
				}
			}
		})
	}
}

func TestSaveMetric(t *testing.T) {
	r := chi.NewRouter()
	h := handler{}
	r.Post("/update/{type}/{name}/{value}", h.SaveMetric)
	srv := httptest.NewServer(r)
	defer srv.Close()

	testCases := []struct {
		name           string
		mType          string
		mName          string
		mValue         string
		expectedCode   int
		expectedMetric collector2.StoredMetric
		expectedError  error
	}{
		{
			name:   "case0",
			mType:  "counter",
			mName:  "Counter1",
			mValue: "15",
			expectedMetric: collector2.StoredMetric{
				ID:           "Counter1",
				MType:        "counter",
				CounterValue: collector2.PtrInt64(15),
				TextValue:    collector2.PtrString("15"),
			},
			expectedCode: http.StatusOK,
		},
		{
			name:   "case1",
			mType:  "gauge",
			mName:  "Gauge1",
			mValue: "12.282",
			expectedMetric: collector2.StoredMetric{
				ID:         "Gauge1",
				MType:      "gauge",
				GaugeValue: collector2.PtrFloat64(12.282),
				TextValue:  collector2.PtrString("12.282"),
			},
			expectedCode: http.StatusOK,
		},
		{
			name:          "case2",
			mType:         "invalid",
			mName:         "Gauge1",
			mValue:        "12.282",
			expectedCode:  http.StatusNotImplemented,
			expectedError: collector2.ErrNotFound,
		},
		{
			name:          "case3",
			mType:         "counter",
			mName:         "Counter1",
			mValue:        "15.2562",
			expectedCode:  http.StatusBadRequest,
			expectedError: collector2.ErrNotFound,
		},
		{
			name:          "case4",
			mType:         "gauge",
			mName:         "Gauge1",
			mValue:        "12.282dgh",
			expectedCode:  http.StatusBadRequest,
			expectedError: collector2.ErrNotFound,
		},
		{
			name:          "case5",
			mType:         "gauge",
			mName:         "Gauge1",
			mValue:        "",
			expectedCode:  http.StatusNotFound,
			expectedError: collector2.ErrNotFound,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := resty.New().R().
				SetHeader("Content-Type", "text/plain").
				Post(fmt.Sprintf("%s/update/%s/%s/%s", srv.URL, tt.mType, tt.mName, tt.mValue))

			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, resp.StatusCode(), tt.expectedCode)

			value, err := collector2.Collector().GetMetric(tt.mName)
			if err != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			if tt.expectedCode == http.StatusOK {
				assert.Equal(t, value, tt.expectedMetric)
			}
		})
	}
}

func TestSaveMetricFromJSON(t *testing.T) {
	testCases := []struct {
		name           string
		request        collector2.MetricRequest
		expectedMetric collector2.StoredMetric
		header         map[string]string
		trustedSubnet  string
		expectedCode   int
		expectedError  error
	}{
		{
			name: "positive (counter)",
			request: collector2.MetricRequest{
				MType: "counter",
				ID:    "Counter15",
				Delta: collector2.PtrInt64(15),
			},
			expectedMetric: collector2.StoredMetric{
				MType:        "counter",
				ID:           "Counter15",
				CounterValue: collector2.PtrInt64(15),
				TextValue:    collector2.PtrString("15"),
			},
			expectedCode: http.StatusOK,
			header: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name: "positive (gauge)",
			request: collector2.MetricRequest{
				MType: "gauge",
				ID:    "Gauge1",
				Value: collector2.PtrFloat64(12.282),
			},
			expectedMetric: collector2.StoredMetric{
				MType:      "gauge",
				ID:         "Gauge1",
				GaugeValue: collector2.PtrFloat64(12.282),
				TextValue:  collector2.PtrString("12.28200000000"),
			},
			expectedCode: http.StatusOK,
			header: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name: "positive (with trusted subnet)",
			request: collector2.MetricRequest{
				MType: "gauge",
				ID:    "Gauge1",
				Value: collector2.PtrFloat64(12.282),
			},
			expectedMetric: collector2.StoredMetric{
				MType:      "gauge",
				ID:         "Gauge1",
				GaugeValue: collector2.PtrFloat64(12.282),
				TextValue:  collector2.PtrString("12.28200000000"),
			},
			expectedCode: http.StatusOK,
			header: map[string]string{
				"Content-Type": "application/json",
				"X-Real-IP":    "172.17.0.2",
			},
			trustedSubnet: "172.17.0.0/16",
		},
		{
			name: "negative (invalid type)",
			request: collector2.MetricRequest{
				MType: "invalid",
				ID:    "Gauge1",
				Value: collector2.PtrFloat64(12.282),
			},
			expectedMetric: collector2.StoredMetric{},
			expectedCode:   http.StatusNotImplemented,
			expectedError:  collector2.ErrNotImplemented,
			header: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name: "negative (invalid name)",
			request: collector2.MetricRequest{
				MType: "gauge",
				ID:    "",
				Value: collector2.PtrFloat64(1),
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: collector2.ErrNotFound,
			header: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name: "negative (invalid gauge value)",
			request: collector2.MetricRequest{
				MType: "gauge",
				ID:    "invalidGauge",
				Value: collector2.PtrFloat64(-1.9),
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: collector2.ErrNotFound,
			header: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name: "negative (with trusted subnet)",
			request: collector2.MetricRequest{
				MType: "gauge",
				ID:    "Gauge1",
				Value: collector2.PtrFloat64(12.282),
			},
			expectedMetric: collector2.StoredMetric{
				MType:      "gauge",
				ID:         "Gauge1",
				GaugeValue: collector2.PtrFloat64(12.282),
				TextValue:  collector2.PtrString("12.28200000000"),
			},
			expectedCode: http.StatusForbidden,
			header: map[string]string{
				"Content-Type": "application/json",
				"X-Real-IP":    "127.0.0.1",
			},
			trustedSubnet: "172.17.0.0/16",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			h := handler{
				trustedSubnet: tt.trustedSubnet,
			}
			r.Use(h.CheckSubnet)
			r.Post("/update/", h.SaveMetricFromJSON)
			srv := httptest.NewServer(r)
			defer srv.Close()

			resBody, err := json.Marshal(tt.request)
			assert.NoError(t, err)
			req := resty.New().R()
			for k, v := range tt.header {
				req.SetHeader(k, v)
			}

			resp, err := req.SetBody(resBody).Post(fmt.Sprintf("%s/update/", srv.URL))
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, resp.StatusCode(), tt.expectedCode)

			value, err := collector2.Collector().GetMetricJSON(tt.request.ID)
			if err != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			actual := collector2.StoredMetric{}
			json.Unmarshal(value, &actual)

			if tt.expectedCode == http.StatusOK {
				assert.Equal(t, actual, tt.expectedMetric)
			}
		})
	}
}

func TestGetMetric(t *testing.T) {
	r := chi.NewRouter()
	h := handler{}
	r.Use(h.CheckSubscription)
	r.Post("/update/{type}/{name}/{value}", h.SaveMetric)
	r.Get("/value/{type}/{name}", h.GetMetric)
	srv := httptest.NewServer(r)
	defer srv.Close()

	client := resty.New()
	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/counter/Counter3/15", srv.URL))
	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/counter/Counter2/0", srv.URL))

	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/gauge/Gauge1/100500.2780001", srv.URL))
	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/gauge/Gauge2/100500.278000100", srv.URL))
	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/gauge/Gauge3/100500", srv.URL))

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
			mName:        "Counter3",
			mValue:       "15",
			expectedCode: http.StatusOK,
		},
		{
			name:         "case1",
			mType:        "counter",
			mName:        "Counter2",
			mValue:       "0",
			expectedCode: http.StatusOK,
		},
		{
			name:         "case2",
			mType:        "gauge",
			mName:        "Gauge1",
			mValue:       "100500.2780001",
			expectedCode: http.StatusOK,
		},
		{
			name:         "case3",
			mType:        "gauge",
			mName:        "Gauge2",
			mValue:       "100500.278000100",
			expectedCode: http.StatusOK,
		},
		{
			name:         "case4",
			mType:        "gauge",
			mName:        "Gauge3",
			mValue:       "100500",
			expectedCode: http.StatusOK,
		},
		{
			name:         "case5",
			mType:        "gauge",
			mName:        "Gauge4",
			mValue:       "",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "case6",
			mType:        "invalid",
			mName:        "Gauge4",
			mValue:       "",
			expectedCode: http.StatusNotImplemented,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := resty.New().R().
				SetHeader("Content-Type", "text/plain").
				Get(fmt.Sprintf("%s/value/%s/%s", srv.URL, tt.mType, tt.mName))

			assert.NoError(t, err)
			assert.Equal(t, resp.StatusCode(), tt.expectedCode)
			assert.Equal(t, string(resp.Body()), tt.mValue)
		})
	}
}

func TestGetMetricFromJSON(t *testing.T) {
	r := chi.NewRouter()
	h := handler{}
	r.Post("/update/{type}/{name}/{value}", h.SaveMetric)
	r.Post("/value/", h.GetMetricFromJSON)
	srv := httptest.NewServer(r)
	defer srv.Close()

	client := resty.New()
	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/counter/Counter3/15", srv.URL))
	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/counter/Counter2/0", srv.URL))

	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/gauge/Gauge1/100500.2780001", srv.URL))
	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/gauge/Gauge2/100500.278000100", srv.URL))
	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/gauge/Gauge3/100500", srv.URL))

	testCases := []struct {
		name          string
		mType         string
		mName         string
		mValue        float64
		mDelta        int64
		expectedCode  int
		expectedError error
	}{
		{
			name:         "case0",
			mType:        "counter",
			mName:        "Counter3",
			mDelta:       15,
			expectedCode: http.StatusOK,
		},
		{
			name:         "case1",
			mType:        "counter",
			mName:        "Counter2",
			mDelta:       0,
			expectedCode: http.StatusOK,
		},
		{
			name:         "case2",
			mType:        "gauge",
			mName:        "Gauge1",
			mValue:       100500.2780001,
			expectedCode: http.StatusOK,
		},
		{
			name:         "case3",
			mType:        "gauge",
			mName:        "Gauge2",
			mValue:       100500.278000100,
			expectedCode: http.StatusOK,
		},
		{
			name:         "case4",
			mType:        "gauge",
			mName:        "Gauge3",
			mValue:       100500,
			expectedCode: http.StatusOK,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			body := collector2.MetricRequest{
				ID:    tt.mName,
				MType: tt.mType,
			}
			resBody, err := json.Marshal(body)
			assert.NoError(t, err)

			resp, err := resty.New().R().
				SetBody(resBody).
				Post(fmt.Sprintf("%s/value/", srv.URL))

			assert.NoError(t, err)
			assert.Equal(t, resp.StatusCode(), tt.expectedCode)
		})
	}
}

func TestShowMetrics(t *testing.T) {
	r := chi.NewRouter()
	h := handler{}
	r.Post("/update/{type}/{name}/{value}", h.SaveMetric)
	r.Get("/", h.ShowMetrics)
	srv := httptest.NewServer(r)
	defer srv.Close()

	client := resty.New()
	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/counter/Counter3/15", srv.URL))
	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/counter/Counter2/0", srv.URL))

	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/gauge/Gauge1/100500.2780001", srv.URL))
	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/gauge/Gauge2/100500.278000100", srv.URL))
	_, _ = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%s/update/gauge/Gauge3/100500", srv.URL))

	testCases := []struct {
		name         string
		expectedPage string
		expectedCode int
	}{
		{
			name:         "case0",
			expectedCode: http.StatusOK,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := resty.New().R().
				SetHeader("Content-Type", "text/plain").
				Get(fmt.Sprintf("%s/", srv.URL))

			assert.NoError(t, err)
			assert.Equal(t, resp.StatusCode(), tt.expectedCode)
		})
	}
}
