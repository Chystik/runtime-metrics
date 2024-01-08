package service

import (
	"context"
	"net/http"

	"github.com/Chystik/runtime-metrics/pkg/logger"
)

type DBClient interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Migrate() error
	Ping(ctx context.Context) error
}

type AppLogger = logger.Logger

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type ConnectionRetrier interface {
	DoWithRetry(retryableFunc func() error) error
}

type ConnectionRetrierFn interface {
	DoWithRetryFn() error
}
