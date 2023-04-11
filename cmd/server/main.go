package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/handlers"
	"log"
	"net/http"
	"os"
)

var flagRunAddr string

func parseFlags() {
	cnvFlags := flag.NewFlagSet("cnv", flag.ContinueOnError)
	cnvFlags.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	err := cnvFlags.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
}

func main() {
	parseFlags()
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", handlers.SaveMetric)
	r.Get("/value/{type}/{name}", handlers.GetMetric)
	r.Get("/", handlers.ShowMetrics)

	log.Fatal(http.ListenAndServe(flagRunAddr, r))
}
