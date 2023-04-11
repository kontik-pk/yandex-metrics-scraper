package handlers

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/collector"
	"io"
	"net/http"
	"strconv"
)

func SaveMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	metricValue := chi.URLParam(r, "value")

	err := collector.Collector.Collect(metricName, metricType, metricValue)
	if errors.Is(err, collector.ErrBadRequest) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if errors.Is(err, collector.ErrNotImplemented) {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	io.WriteString(w, "")
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.Header().Set("content-length", strconv.Itoa(len(metricName)))
	w.WriteHeader(http.StatusOK)
}

func GetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	value, err := collector.Collector.GetMetric(metricName, metricType)
	if errors.Is(err, collector.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errors.Is(err, collector.ErrNotImplemented) {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	io.WriteString(w, "")
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.Header().Set("content-length", strconv.Itoa(len(value)))
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, value)
}

func ShowMetrics(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	page := `
<html> 
   <head> 
   </head> 
   <body> 
`
	for _, n := range collector.Collector.GetAvailableMetrics() {
		page += fmt.Sprintf(`<h3>%s   </h3>`, n)
	}
	page += `
   </body> 
</html>
`
	w.Header().Set("content-type", "Content-Type: text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(page))
}
