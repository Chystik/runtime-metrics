package interceptors

import (
	"context"
	"time"

	"github.com/Chystik/runtime-metrics/internal/service"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func UnaryServerLogger(l service.AppLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()

		l.Info(
			"GRPC request started",
			zap.String("method", info.FullMethod),
		)

		m, err := handler(ctx, req)
		if err != nil {
			l.Info(
				"gRPC failed with error",
				zap.String("err", err.Error()),
			)
		}

		l.Info(
			"GRPC response completed",
			zap.Duration("duration", time.Since(start)),
		)

		return m, err
	}
}
