package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kontik-pk/yandex-metrics-scraper/internal/handlers"
)

func main() {
	client := resty.New()
	r, err := client.R().SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(`{"id": "RandomValue", "type": "gauge", "value": 1}`).
		Post("http://localhost:8080/update/")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(r.Body()))
	var m handlers.Metrics
	if err := json.Unmarshal(r.Body(), &m); err != nil {
		panic(err)
	}
}
