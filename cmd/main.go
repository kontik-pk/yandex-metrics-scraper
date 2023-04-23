package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/collector"
)

func main() {
	client := resty.New()
	r, err := client.R().SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(`{"id": "RandomValue", "type": "gauge"}`).
		Post("http://localhost:8080/value/")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(r.Body()))
	var m collector.MetricJSON
	if err := json.Unmarshal(r.Body(), &m); err != nil {
		panic(err)
	}
}
