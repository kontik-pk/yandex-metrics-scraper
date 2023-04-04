package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type metric struct {
	value any
	mType string
}

type MemStorage struct {
	metrics map[string]metric
}

var m MemStorage

func SaveMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if len(m.metrics) == 0 {
		m.metrics = make(map[string]metric, 0)
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
		m.metrics[url[2]] = metric{value, url[2]}
	case "gauge":
		value, err := strconv.ParseFloat(url[4], 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		m.metrics[url[3]] = metric{value, url[2]}
	default:
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	io.WriteString(w, "")
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.Header().Set("content-length", strconv.Itoa(len(url[3])))
	w.WriteHeader(http.StatusOK)
	fmt.Println(m.metrics)
	fmt.Println(r.URL)
}
