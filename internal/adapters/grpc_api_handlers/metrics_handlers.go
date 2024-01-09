package grpcapihandlers

import (
	"context"

	"github.com/Chystik/runtime-metrics/internal/service"
	pb "github.com/Chystik/runtime-metrics/protobuf"
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
		response.Error.Message = err.Error()
	}

	return &response, err
}
