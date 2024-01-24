package grpcapihandlers

import (
	"github.com/Chystik/runtime-metrics/internal/models"
	pb "github.com/Chystik/runtime-metrics/protobuf"
)

func toDomainMetrics(m []*pb.Metric) []models.Metric {
	res := make([]models.Metric, len(m))

	for i := range m {
		res[i] = toDomainMetric(m[i])
	}

	return res
}

func toDomainMetric(m *pb.Metric) models.Metric {
	return models.Metric{
		ID:    m.Id,
		MType: m.Type,
		Delta: &m.Delta,
		Value: &m.Value,
	}
}

func fromDomainMetric(m models.Metric) *pb.Metric {
	return &pb.Metric{
		Id:    m.ID,
		Type:  m.MType,
		Delta: *m.Delta,
		Value: *m.Value,
	}
}
