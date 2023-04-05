package handlers

import (
	"fmt"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/domain"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func SaveMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if len(domain.MetricsStorage.Metrics) == 0 {
		domain.MetricsStorage.Metrics = make(map[string]domain.Metric, 0)
	}

	url := strings.Split(r.URL.String(), "/")
	if len(url) != 5 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if url[1] != "update" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch url[2] {
	case "counter":
		value, err := strconv.Atoi(url[4])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if domain.MetricsStorage.Metrics[url[3]].Value != nil {
			value += domain.MetricsStorage.Metrics[url[3]].Value.(int)
		}
		domain.MetricsStorage.Metrics[url[3]] = domain.Metric{Value: value, MType: url[2]}
	case "gauge":
		_, err := strconv.ParseFloat(url[4], 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		domain.MetricsStorage.Metrics[url[3]] = domain.Metric{Value: url[4], MType: url[2]}
	default:
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	io.WriteString(w, "")
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.Header().Set("content-length", strconv.Itoa(len(url[3])))
	w.WriteHeader(http.StatusOK)
}

func GetMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	url := strings.Split(r.URL.String(), "/")
	if len(url) != 4 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if url[1] != "value" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if _, ok := domain.MetricsStorage.Metrics[url[3]]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	value := domain.MetricsStorage.Metrics[url[3]].Value

	io.WriteString(w, "")
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.Header().Set("content-length", strconv.Itoa(len(url[3])))
	w.WriteHeader(http.StatusOK)

	switch value.(type) {
	case uint, uint64, int, int64:
		io.WriteString(w, strconv.Itoa(value.(int)))
	default:
		io.WriteString(w, value.(string))
	}
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
	for n := range domain.MetricsStorage.Metrics {
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
