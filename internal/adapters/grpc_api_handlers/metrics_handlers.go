package grpcapihandlers

import (
	"context"

	"github.com/Chystik/runtime-metrics/internal/service"
	pb "github.com/Chystik/runtime-metrics/protobuf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type metricsHandlers struct {
	metricsService service.MetricsService
	pb.UnimplementedMetricsServiceServer
}

func NewMetricsHandlers(ms service.MetricsService) *metricsHandlers {
	return &metricsHandlers{
		metricsService: ms,
	}
}

func (mh *metricsHandlers) UpdateMetrics(ctx context.Context, m *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	var response pb.UpdateMetricsResponse

	err := mh.metricsService.UpdateList(ctx, toDomainMetrics(m.Metrics))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update list error: %s", err.Error())
	}

	return &response, nil
}

func (mh *metricsHandlers) UpdateMetric(ctx context.Context, m *pb.UpdateMetricRequest) (*pb.UpdateMetricResponse, error) {
	var response pb.UpdateMetricResponse
	var err error

	switch m.Metric.Type {
	case "gauge":
		err = mh.metricsService.UpdateGauge(ctx, toDomainMetric(m.Metric))
	case "counter":
		err = mh.metricsService.UpdateCounter(ctx, toDomainMetric(m.Metric))
	default:
		return nil, status.Errorf(codes.NotFound, "unknown metric type: %s", m.Metric.Type)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "update %s error: %s", m.Metric.Type, err.Error())
	}

	return &response, nil
}

func (mh *metricsHandlers) GetMetric(ctx context.Context, req *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	var response pb.GetMetricResponse
	var err error

	m, err := mh.metricsService.Get(ctx, toDomainMetric(req.Metric))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "can't find metric with id %s", req.Metric.Id)
	}
	response.Metric = fromDomainMetric(m)

	return &response, nil
}
