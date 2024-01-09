package grpcclient

import (
	"github.com/Chystik/runtime-metrics/internal/models"
	pb "github.com/Chystik/runtime-metrics/protobuf"
)

func fromDomainMetrics(m map[string]models.Metric) []*pb.Metric {
	res := make([]*pb.Metric, 0, len(m))

	for _, v := range m {
		res = append(res, fromDomainMetric(v))
	}

	return res
}

func fromDomainMetric(m models.Metric) *pb.Metric {
	res := &pb.Metric{
		Id:   m.ID,
		Type: m.MType,
	}

	if m.Delta != nil {
		res.Delta = *m.Delta
	}
	if m.Value != nil {
		res.Value = *m.Value
	}

	return res
}
