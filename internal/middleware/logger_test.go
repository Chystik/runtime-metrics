package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Chystik/runtime-metrics/internal/models"
	"github.com/Chystik/runtime-metrics/pkg/logger"
)

var nextLoggerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
})

func makeLoggerRequest(h http.Handler, b []byte) int {
	var body bytes.Buffer
	_, _ = body.Write(b)

	req := httptest.NewRequest(http.MethodPost, "http://testing", &body)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	res := rec.Result()

	defer res.Body.Close()

	return res.StatusCode
}

func Test_midLogger_WithLogging(t *testing.T) {
	logger, _ := logger.Initialize("info")

	type args struct {
		handlerToTest http.Handler
		testData      []models.Metric
	}
	tests := []struct {
		name       string
		l          *midLogger
		args       args
		wantStatus int
	}{
		{
			name: "logger returns 200",
			l:    MidLogger(logger),
			args: args{
				handlerToTest: nextLoggerHandler,
				testData:      generateMetrics(10),
			},
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, d := range tt.args.testData {
				var buf bytes.Buffer
				if err := json.NewEncoder(&buf).Encode(d); err != nil {
					t.Errorf(err.Error())
				}

				if got := makeLoggerRequest(tt.l.WithLogging(tt.args.handlerToTest), buf.Bytes()); got != tt.wantStatus {
					t.Errorf("midLogger.WithLogging() = %v, want %v", got, tt.wantStatus)
				}
			}
		})
	}
}
