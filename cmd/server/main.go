package main

import (
	"github.com/kontik-pk/yandex-metrics-scraper/internal/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/update/", handlers.SaveMetric)

	if err := http.ListenAndServe(`:8080`, nil); err != nil {
		panic(err)
	}
}
