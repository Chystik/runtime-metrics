package interceptors

import (
	"runtime/debug"

	"github.com/Chystik/runtime-metrics/internal/service"
	"go.uber.org/zap"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryServerRecoverer(l service.AppLogger) grpc.UnaryServerInterceptor {
	var grpcPanicRecoveryHandler = func(p any) error {
		l.Error(
			"recovered from panic",
			zap.Any("panic", p),
			zap.String("stack", string(debug.Stack())),
		)
		return status.Errorf(codes.Internal, "%s", p)
	}
	return recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler))
}
