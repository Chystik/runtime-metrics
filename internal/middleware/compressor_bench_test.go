package middleware

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Chystik/runtime-metrics/internal/models"
)

/*
goos: windows
goarch: amd64
pkg: github.com/Chystik/runtime-metrics/internal/middleware
cpu: Intel(R) Core(TM) i5-10300H CPU @ 2.50GHz
BenchmarkWriters-8              	     429	   2882737 ns/op	   816320 B/op	      17 allocs/op
BenchmarkPoolWriters-8          	     422	   2592022 ns/op	     4415 B/op	       0 allocs/op
BenchmarkCompressor/Normal-8    	      13	 102234677 ns/op	457133233 B/op	   42296 allocs/op
BenchmarkCompressor/Pool-8      	      38	  30838611 ns/op	 62529156 B/op	   32636 allocs/op
BenchmarkCompressorParallel/Normal-8      30	  41232880 ns/op	456668210 B/op	   41498 allocs/op
BenchmarkCompressorParallel/Pool-8        62	  16883881 ns/op	113411487 B/op	   34112 allocs/op
*/
const (
	metricsCount = 500
)

type metricBody struct {
	body []byte
}

func BenchmarkWriters(b *testing.B) {
	w := io.Discard
	d := make([]byte, 1024*1024)
	for n := 0; n < b.N; n++ {
		z := gzip.NewWriter(w)
		z.Write(d)
		z.Close()
	}
}

func BenchmarkPoolWriters(b *testing.B) {
	w := io.Discard
	var zippers sync.Pool
	d := make([]byte, 1024*1024)
	for n := 0; n < b.N; n++ {
		z, _ := zippers.Get().(*gzip.Writer)
		if z == nil {
			z = gzip.NewWriter(w)
		} else {
			z.Reset(w)
		}
		z.Write(d)
		z.Close()
		zippers.Put(z)
	}
}

var nextCompressorHandlerB = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	contentEncoding := r.Header.Get("Content-Encoding")
	sendsGzip := strings.Contains(contentEncoding, "gzip")

	if !sendsGzip {
		panic("answer not compressed")
	}

	var m models.Metric

	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		panic(err)
	}

	var body bytes.Buffer

	err = json.NewEncoder(&body).Encode(m)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(body.Bytes())
})

// makeCompressorRequest makes test request using compressor middleware.
// it panics if error occurs
func makeCompressorRequestB(h http.Handler, b []byte) {
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
		panic("expect http.StatusOK")
	}

	reader, err := gzip.NewReader(res.Body)
	if err != nil {
		panic(err)
	}

	var p bytes.Buffer
	_, err = io.Copy(&p, reader)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	// does something with decompressed response
	_ = p.String
}

func BenchmarkCompressor(b *testing.B) {
	b.Run("Normal", func(b *testing.B) {
		handlerToTest := GzipMiddleware(nextCompressorHandlerB)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			b.StopTimer()
			reqMetrics := compressMetrics(generateMetrics(metricsCount))
			b.StartTimer()

			for i := range reqMetrics {
				makeCompressorRequestB(handlerToTest, reqMetrics[i].body)
			}
		}
	})

	b.Run("Pool", func(b *testing.B) {
		gz := GzipPoolMiddleware()
		handlerToTest := gz(nextCompressorHandlerB)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			b.StopTimer()
			reqMetrics := compressMetrics(generateMetrics(metricsCount))
			b.StartTimer()

			for i := range reqMetrics {
				makeCompressorRequestB(handlerToTest, reqMetrics[i].body)
			}
		}
	})
}

func BenchmarkCompressorParallel(b *testing.B) {
	b.Run("Normal", func(b *testing.B) {
		handlerToTest := GzipMiddleware(nextCompressorHandlerB)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			b.StopTimer()
			reqMetrics := compressMetrics(generateMetrics(metricsCount))
			var wg sync.WaitGroup
			wg.Add(metricsCount)
			b.StartTimer()

			for i := range reqMetrics {
				i := i
				go func() {
					defer wg.Done()
					makeCompressorRequestB(handlerToTest, reqMetrics[i].body)

				}()
			}
			wg.Wait()
		}
	})

	b.Run("Pool", func(b *testing.B) {
		gz := GzipPoolMiddleware()
		handlerToTest := gz(nextCompressorHandlerB)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			b.StopTimer()
			reqMetrics := compressMetrics(generateMetrics(metricsCount))
			var wg sync.WaitGroup
			wg.Add(metricsCount)
			b.StartTimer()

			for i := range reqMetrics {
				i := i
				go func() {
					defer wg.Done()
					makeCompressorRequestB(handlerToTest, reqMetrics[i].body)
				}()
			}
			wg.Wait()
		}
	})
}

func generateMetrics(count int) []models.Metric {
	m := make([]models.Metric, count)

	rand.New(rand.NewSource(time.Now().UnixNano()))

	randMetricType := [2]string{"gauge", "counter"}

	for i := range m {
		mName := fmt.Sprintf("TestMetric%d", i)
		mType := randMetricType[rand.Intn(2)]
		m[i] = generateMetric(mType, mName)
	}

	return m
}

func generateMetric(metricType string, metricName string) models.Metric {
	var m models.Metric

	min := 1e1
	max := 1e3

	m.ID = metricName
	m.MType = metricType

	switch metricType {
	case "gauge":
		m.Value = new(float64)
		*m.Value = min + rand.Float64()*(max-min)
	case "counter":
		m.Delta = new(int64)
		*m.Delta = int64(min + rand.Float64()*(max-min))
	}

	return m
}

func compressMetric(metric models.Metric) metricBody {
	var m metricBody

	var buf, cm bytes.Buffer

	err := json.NewEncoder(&buf).Encode(metric)
	if err != nil {
		log.Fatal(err)
	}

	gz := gzip.NewWriter(&cm)
	_, err = gz.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	err = gz.Close()
	if err != nil {
		log.Fatal(err)
	}

	m.body = cm.Bytes()

	return m
}

func compressMetrics(metrics []models.Metric) []metricBody {
	m := make([]metricBody, len(metrics))

	for i := range metrics {
		m[i] = compressMetric(metrics[i])
	}

	return m
}

func encodeMetric(metric models.Metric) metricBody {
	var m metricBody
	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(metric)
	if err != nil {
		log.Fatal(err)
	}

	m.body = buf.Bytes()

	return m
}

func encodeMetrics(metrics []models.Metric) []metricBody {
	m := make([]metricBody, len(metrics))

	for i := range metrics {
		m[i] = encodeMetric(metrics[i])
	}

	return m
}
