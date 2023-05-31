package handlers

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/collector"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (h *handler) SaveMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	metricValue := chi.URLParam(r, "value")

	metric := collector.MetricRequest{
		ID:    metricName,
		MType: metricType,
	}
	err := collector.Collector.Collect(metric, metricValue)
	if errors.Is(err, collector.ErrBadRequest) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if errors.Is(err, collector.ErrNotImplemented) {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := io.WriteString(w, fmt.Sprintf("inserted metric %q with value %q", metricName, metricValue)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.Header().Set("content-length", strconv.Itoa(len(metricName)))
}

func (h *handler) SaveMetricFromJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !h.checkSubscription(w, buf, r.Header.Get("HashSHA256")) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var metric collector.MetricRequest
	if err := json.Unmarshal(buf.Bytes(), &metric); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//TODO: не самое изящное архитектурное решение, как тут можно сделать лучше?
	metricValue := ""
	switch metric.MType {
	case collector.Counter:
		metricValue = strconv.Itoa(int(*metric.Delta))
	case collector.Gauge:
		metricValue = strconv.FormatFloat(*metric.Value, 'f', 11, 64)
	default:
	}

	err := collector.Collector.Collect(metric, metricValue)
	if errors.Is(err, collector.ErrBadRequest) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if errors.Is(err, collector.ErrNotImplemented) {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	resultJSON, err := collector.Collector.GetMetricJSON(metric.ID)
	if errors.Is(err, collector.ErrBadRequest) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if errors.Is(err, collector.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errors.Is(err, collector.ErrNotImplemented) {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	if _, err = w.Write(resultJSON); err != nil {
		return
	}
	w.Header().Set("content-length", strconv.Itoa(len(metric.ID)))
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *handler) SaveListMetricsFromJSON(w http.ResponseWriter, r *http.Request) {
	fmt.Println(collector.Collector.Metrics)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !h.checkSubscription(w, buf, r.Header.Get("HashSHA256")) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var metrics []collector.MetricRequest
	if err := json.Unmarshal(buf.Bytes(), &metrics); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var results []byte
	for _, metric := range metrics {
		//TODO: не самое изящное архитектурное решение, как тут можно сделать лучше?
		metricValue := ""
		switch metric.MType {
		case collector.Counter:
			metricValue = strconv.Itoa(int(*metric.Delta))
		case collector.Gauge:
			metricValue = strconv.FormatFloat(*metric.Value, 'f', 11, 64)
		default:
		}

		err := collector.Collector.Collect(metric, metricValue)
		if errors.Is(err, collector.ErrBadRequest) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if errors.Is(err, collector.ErrNotImplemented) {
			w.WriteHeader(http.StatusNotImplemented)
			return
		}

		resultJSON, err := collector.Collector.GetMetricJSON(metric.ID)
		if errors.Is(err, collector.ErrBadRequest) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if errors.Is(err, collector.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, collector.ErrNotImplemented) {
			w.WriteHeader(http.StatusNotImplemented)
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

func (h *handler) GetMetricFromJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	gotHash := r.Header.Get("HashSHA256")
	want := h.getHash(buf.Bytes())
	if gotHash != "" {
		w.Header().Set("HashSHA256", want)
	}
	if !h.checkSubscription(w, buf, r.Header.Get("HashSHA256")) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var metric collector.MetricRequest
	if err := json.Unmarshal(buf.Bytes(), &metric); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resultJSON, err := collector.Collector.GetMetric(metric.ID)
	if errors.Is(err, collector.ErrBadRequest) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if errors.Is(err, collector.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errors.Is(err, collector.ErrNotImplemented) {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	switch metric.MType {
	case collector.Counter:
		metric.Delta = resultJSON.CounterValue
	case collector.Gauge:
		metric.Value = resultJSON.GaugeValue
	}
	answer, _ := json.Marshal(metric)

	if _, err = w.Write(answer); err != nil {
		return
	}
	w.Header().Set("content-length", strconv.Itoa(len(metric.ID)))
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Println("HEADERS: ", w.Header())
}

func (h *handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	if metricType != collector.Counter && metricType != collector.Gauge {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	value, err := collector.Collector.GetMetric(metricName)
	if errors.Is(err, collector.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errors.Is(err, collector.ErrNotImplemented) {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = io.WriteString(w, *value.TextValue); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "text/plain; charset=utf-8")
}

func (h *handler) ShowMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "Content-Type: text/html; charset=utf-8")
	if r.URL.Path != "/" {
		http.Error(w, fmt.Sprintf("wrong path %q", r.URL.Path), http.StatusNotFound)
		return
	}
	page := ""
	for _, n := range collector.Collector.GetAvailableMetrics() {
		page += fmt.Sprintf("<h1>	%s</h1>", n)
	}
	tmpl, _ := template.New("data").Parse("<h1>AVAILABLE METRICS</h1>{{range .}}<h3>{{ .}}</h3>{{end}}")
	if err := tmpl.Execute(w, collector.Collector.GetAvailableMetrics()); err != nil {
		return
	}
	w.Header().Set("content-type", "Content-Type: text/html; charset=utf-8")
}

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
	} else {
		w.WriteHeader(http.StatusOK)
	}
	if _, err = w.Write([]byte("pong")); err != nil {
		return
	}
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

func (h *handler) getHash(body []byte) string {
	want := sha256.Sum256(body)
	wantDecoded := fmt.Sprintf("%x", want)
	return wantDecoded
}

func New(db string, key string) *handler {
	return &handler{
		dbAddress: db,
		key:       key,
	}
}

type handler struct {
	dbAddress string
	key       string
}
