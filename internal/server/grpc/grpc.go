package grpc

import (
	"context"
	collector2 "github.com/kontik-pk/yandex-metrics-scraper/internal/agent/collector"
	pb "github.com/kontik-pk/yandex-metrics-scraper/proto"
	"log"
	"strconv"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServer
}

func (s *MetricsServer) SaveMetricFromJSON(ctx context.Context, in *pb.MetricRequest) (*pb.SaveMetricResponse, error) {
	c := collector2.Collector()
	metric := collector2.MetricRequest{
		ID:    in.ID,
		MType: in.MType,
	}

	// get metric value
	var metricValue string
	switch in.MType {
	case collector2.Counter:
		metricValue = strconv.Itoa(int(in.Delta))
		metric.Delta = &in.Delta
	case collector2.Gauge:
		metricValue = strconv.FormatFloat(in.Value, 'f', 11, 64)
		metric.Value = &in.Value
	default:
		return &pb.SaveMetricResponse{
			ResultJSON: nil,
			Error:      collector2.ErrNotImplemented.Error(),
		}, collector2.ErrNotImplemented
	}

	if err := c.Collect(metric, metricValue); err != nil {
		return &pb.SaveMetricResponse{
			ResultJSON: nil,
			Error:      err.Error(),
		}, err
	}

	// get saved metric in JSON format for response
	resultJSON, err := c.GetMetricJSON(metric.ID)
	if err != nil {
		return &pb.SaveMetricResponse{
			ResultJSON: nil,
			Error:      err.Error(),
		}, err
	}
	log.Println(string(resultJSON))
	return &pb.SaveMetricResponse{
		ResultJSON: resultJSON,
		Error:      "",
	}, nil
}
