package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

func main() {
	client := resty.New()
	r, err := client.R().SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		//SetBody(`{"id": "RandomValue", "type": "gauge"}`).
		Post("http://localhost:8080/")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(r.Body()))
	fmt.Println(r.Header())
	//var m collector.MetricJSON
	//if err := json.Unmarshal(r.Body(), &m); err != nil {
	//	panic(err)
	//}
}
