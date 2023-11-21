package middleware

import (
	"net/http"
	"time"

	"github.com/Chystik/runtime-metrics/internal/service"
	"go.uber.org/zap"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}

	midLogger struct {
		service.AppLogger
	}
)

func MidLogger(l service.AppLogger) *midLogger {
	return &midLogger{l}
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	if err != nil {
		return 0, err
	}

	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func (l *midLogger) WithLogging(next http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		l.Info(
			"request started",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
		)

		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		l.Info(
			"response completed",
			zap.Int("status", responseData.status),
			zap.Duration("duration", duration),
			zap.Int("size", responseData.size),
		)
	}
	return http.HandlerFunc(logFn)
}
