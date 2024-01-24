package grpcclient

import (
	"context"
	"errors"
	"net"
	"syscall"

	"github.com/Chystik/runtime-metrics/internal/models"
	"github.com/Chystik/runtime-metrics/internal/service"
	pb "github.com/Chystik/runtime-metrics/protobuf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type agentAPIClient struct {
	conn   *grpc.ClientConn
	client pb.MetricsServiceClient
	service.AgentAPIClient
}

func New(addr string) (*agentAPIClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := pb.NewMetricsServiceClient(conn)

	return &agentAPIClient{
		conn:   conn,
		client: c,
	}, nil
}

func (ac *agentAPIClient) ReportMetricsBatch(ctx context.Context, m map[string]models.Metric) error {
	req := &pb.UpdateMetricsRequest{
		Metrics: fromDomainMetrics(m),
	}

	res, err := ac.client.UpdateMetrics(ctx, req)
	if err != nil {
		if status.Code(err) == codes.Unavailable {
			return &net.OpError{Err: syscall.ECONNREFUSED}
		}
		return err
	}

	if res.Error != nil {
		return errors.New(res.Error.String())
	}

	_ = res

	return nil
}

func (ac *agentAPIClient) ConnClose() error {
	return ac.conn.Close()
}
