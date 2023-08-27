package handlers

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/collector"
)

// SaveMetric - a method for saving metric from url.
func (h *handler) SaveMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	metricValue := chi.URLParam(r, "value")

	if err := collector.Collector().Collect(
		collector.MetricRequest{
			ID:    metricName,
			MType: metricType,
		}, metricValue); err != nil {
		w.WriteHeader(h.getStatusOnError(err))
		return
	}

	if _, err := io.WriteString(w, fmt.Sprintf("saved metric %q with value %q", metricName, metricValue)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.Header().Set("content-length", strconv.Itoa(len(metricName)))
}

// SaveMetricFromJSON - a method for saving metric from JSON body of http request.
func (h *handler) SaveMetricFromJSON(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// decrypt message if crypto key was specified
	message := buf.Bytes()
	if h.cryptoKey != nil {
		encryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, h.cryptoKey, message)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		message = encryptedData
	}

	// unmarshall request body and get metric
	var metric collector.MetricRequest
	if err := json.Unmarshal(message, &metric); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// save metric
	resultJSON, err := h.collectMetric(metric)
	if err != nil {
		w.WriteHeader(h.getStatusOnError(err))
		return
	}

	if _, err = w.Write(resultJSON); err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-length", strconv.Itoa(len(metric.ID)))
	w.Header().Set("content-type", "application/json")
}

// SaveListMetricsFromJSON - a method for saving a list of metrics from JSON body of http request.
func (h *handler) SaveListMetricsFromJSON(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// unmarshall request body and get metric
	var metrics []collector.MetricRequest
	if err := json.Unmarshal(buf.Bytes(), &metrics); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var results []byte
	// save all metrics from request
	for _, metric := range metrics {
		resultJSON, err := h.collectMetric(metric)
		if err != nil {
			w.WriteHeader(h.getStatusOnError(err))
			return
		}
		results = append(results, resultJSON...)
	}
	if _, err := w.Write(results); err != nil {
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// GetMetricFromJSON - a method for getting metrics by JSON from http request.
func (h *handler) GetMetricFromJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// unmarshall body and get requested metric name
	var metric collector.MetricRequest
	if err := json.Unmarshal(buf.Bytes(), &metric); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get metric from collector
	resultJSON, err := collector.Collector().GetMetric(metric.ID)
	if err != nil {
		w.WriteHeader(h.getStatusOnError(err))
		return
	}
	// get metric value
	switch metric.MType {
	case collector.Counter:
		metric.Delta = resultJSON.CounterValue
	case collector.Gauge:
		metric.Value = resultJSON.GaugeValue
	}
	answer, err := json.Marshal(metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(answer); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-length", strconv.Itoa(len(metric.ID)))
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// GetMetric - a metric for getting metric from url.
func (h *handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	if metricType != collector.Counter && metricType != collector.Gauge {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	// get requested metric from collector
	value, err := collector.Collector().GetMetric(metricName)
	if err != nil {
		w.WriteHeader(h.getStatusOnError(err))
		return
	}

	if _, err = io.WriteString(w, *value.TextValue); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "text/plain; charset=utf-8")
}

// ShowMetrics - a method for getting all available metrics from server.
func (h *handler) ShowMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "Content-Type: text/html; charset=utf-8")
	if r.URL.Path != "/" {
		http.Error(w, fmt.Sprintf("wrong path %q", r.URL.Path), http.StatusNotFound)
		return
	}
	var page string
	for _, n := range collector.Collector().GetAvailableMetrics() {
		page += fmt.Sprintf("<h1>	%s</h1>", n)
	}
	tmpl, _ := template.New("data").Parse("<h1>AVAILABLE METRICS</h1>{{range .}}<h3>{{ .}}</h3>{{end}}")
	if err := tmpl.Execute(w, collector.Collector().GetAvailableMetrics()); err != nil {
		return
	}
	w.Header().Set("content-type", "Content-Type: text/html; charset=utf-8")
}

// Ping - a method for pinging server DB.
func (h *handler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	db, err := sql.Open("pgx", h.dbAddress)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()
	if err := db.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte("pong")); err != nil {
		return
	}
}

func (h *handler) CheckSubscription(hh http.Handler) http.Handler {
	checkFn := func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body.Close()

		buf := bytes.NewBuffer(bodyBytes)

		gotHash := r.Header.Get("HashSHA256")
		want := h.getHash(buf.Bytes())
		if gotHash != "" {
			w.Header().Set("HashSHA256", want)
		}
		if !h.checkSubscription(w, *buf, r.Header.Get("HashSHA256")) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		hh.ServeHTTP(w, r)
	}
	return http.HandlerFunc(checkFn)
}

func (h *handler) collectMetric(metric collector.MetricRequest) ([]byte, error) {
	c := collector.Collector()

	// get metric value
	var metricValue string
	switch metric.MType {
	case collector.Counter:
		metricValue = strconv.Itoa(int(*metric.Delta))
	case collector.Gauge:
		metricValue = strconv.FormatFloat(*metric.Value, 'f', 11, 64)
	default:
		return nil, collector.ErrNotImplemented
	}

	// save metric
	if err := c.Collect(metric, metricValue); err != nil {
		return nil, err
	}

	// get saved metric in JSON format for response
	resultJSON, err := c.GetMetricJSON(metric.ID)
	if err != nil {
		return nil, err
	}
	return resultJSON, err
}

func (h *handler) checkSubscription(w http.ResponseWriter, buf bytes.Buffer, header string) bool {
	want := h.getHash(buf.Bytes())
	if header != "" {
		w.Header().Set("HashSHA256", want)
	}
	if h.key != "" && len(want) != 0 && header != "" {
		//h := hmac.New(sha256.New, []byte(h.key))
		//h.Write(body)
		//dst := h.Sum(nil)
		//return fmt.Sprintf("%x", dst) == header
		return header == want
	}
	return true
}

func (h *handler) getStatusOnError(err error) int {
	switch {
	case errors.Is(err, collector.ErrBadRequest):
		return http.StatusBadRequest
	case errors.Is(err, collector.ErrNotImplemented):
		return http.StatusNotImplemented
	case errors.Is(err, collector.ErrNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// getHash - a method for getting hash from request body.
func (h *handler) getHash(body []byte) string {
	want := sha256.Sum256(body)
	wantDecoded := fmt.Sprintf("%x", want)
	return wantDecoded
}

func New(db string, key string, cryptoKey string) (*handler, error) {
	handler := &handler{
		dbAddress: db,
		key:       key,
	}
	if cryptoKey != "" {
		b, err := os.ReadFile(cryptoKey)
		if err != nil {
			return nil, fmt.Errorf("error while reading file with crypto private key: %w", err)
		}
		block, _ := pem.Decode(b)
		privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("error parsing private key: %w", err)
		}
		handler.cryptoKey = privateKey.(*rsa.PrivateKey)
	}
	return handler, nil
}

type handler struct {
	dbAddress string
	key       string
	cryptoKey *rsa.PrivateKey
}
