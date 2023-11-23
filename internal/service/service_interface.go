package service

import (
	"context"
	"net/http"
)

type DBClient interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Migrate() error
	Ping(ctx context.Context) error
}

type AppLogger interface {
	Debug(msg string, fields ...any)
	Error(msg string, fields ...any)
	Fatal(msg string, fields ...any)
	Info(msg string, fields ...any)
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type ConnectionRetrier interface {
	DoWithRetry(retryableFunc func() error) error
}

type ConnectionRetrierFn interface {
	DoWithRetryFn() error
}
