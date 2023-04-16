package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/flags"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/handlers"
	"log"
	"net/http"
)

func main() {
	params := flags.Init(flags.WithAddr())
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", handlers.SaveMetric)
	r.Get("/value/{type}/{name}", handlers.GetMetric)
	r.Get("/", handlers.ShowMetrics)

	log.Fatal(http.ListenAndServe(params.FlagRunAddr, r))
}
