package middleware

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Chystik/runtime-metrics/internal/models"
)

var nextCompressorHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	contentEncoding := r.Header.Get("Content-Encoding")
	sendsGzip := strings.Contains(contentEncoding, "gzip")

	if !sendsGzip {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("answer not compressed"))
		return
	}

	var m models.Metric

	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var body bytes.Buffer

	err = json.NewEncoder(&body).Encode(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(body.Bytes())
})

// makeCompressorRequest makes test request using compressor middleware.
func makeCompressorRequest(h http.Handler, b []byte) error {
	var body bytes.Buffer
	_, _ = body.Write(b)

	req := httptest.NewRequest(http.MethodPost, "http://testing", &body)
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	res := rec.Result()

	defer res.Body.Close()

	if http.StatusOK != res.StatusCode {
		return fmt.Errorf("expect http.StatusOK, got %d", res.StatusCode)
	}

	reader, err := gzip.NewReader(res.Body)
	if err != nil {
		return err
	}

	var p bytes.Buffer
	_, err = io.Copy(&p, reader)
	if err != nil {
		return err
	}
	defer reader.Close()

	// does something with decompressed response
	_ = p.String()

	return nil
}

func TestGzipMiddleware(t *testing.T) {
	type args struct {
		handlerToTest http.Handler
		testData      []metricBody
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test compressor with 1000 requests",
			args: args{
				handlerToTest: GzipMiddleware(nextCompressorHandler),
				testData:      compressMetrics(generateMetrics(1000)),
			},
		},
		{
			name: "test compressor with uncompressed request",
			args: args{
				handlerToTest: GzipMiddleware(nextCompressorHandler),
				testData:      encodeMetrics(generateMetrics(10)),
			},
			wantErr: true,
		},
		{
			name: "test compressor pool with 1000 requests",
			args: args{
				handlerToTest: GzipPoolMiddleware()(nextCompressorHandler),
				testData:      compressMetrics(generateMetrics(1000)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, d := range tt.args.testData {
				if err := makeCompressorRequest(tt.args.handlerToTest, d.body); err != nil && !tt.wantErr {
					t.Errorf("test: %s, error: %v", tt.name, err)
				}
			}
		})
	}
}
