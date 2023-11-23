package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Chystik/runtime-metrics/internal/models"
	"github.com/Chystik/runtime-metrics/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewMetricHandlers(t *testing.T) {
	t.Parallel()
	handlers, _ := getMetricsHandlersMocks()

	assert.NotNil(t, handlers)
}

func TestUpdateGaugeMetric(t *testing.T) {
	t.Parallel()
	handlers, mks := getMetricsHandlersMocks()

	expStatus := http.StatusOK
	expContextType := "text/plain"

	req := httptest.NewRequest(http.MethodPost, "/update/gauge/test1/25", nil)
	rec := httptest.NewRecorder()

	mks.metricsService.EXPECT().UpdateGauge(mock.Anything, mock.Anything).Return(nil)
	handlers.UpdateMetric(rec, req)

	res := rec.Result()

	assert.Equal(t, expStatus, res.StatusCode)

	defer res.Body.Close()
	_, err := io.ReadAll(res.Body)

	assert.NoError(t, err)
	assert.Equal(t, expContextType, res.Header.Get("Content-Type"))

}

func TestUpdateGaugeMetric_ServiceReturnsError(t *testing.T) {
	t.Parallel()
	handlers, mks := getMetricsHandlersMocks()

	expStatus := http.StatusInternalServerError

	req := httptest.NewRequest(http.MethodPost, "/update/gauge/test1/25", nil)
	rec := httptest.NewRecorder()

	mks.metricsService.EXPECT().UpdateGauge(mock.Anything, mock.Anything).Return(errors.New("error"))
	handlers.UpdateMetric(rec, req)
	res := rec.Result()

	assert.Equal(t, expStatus, res.StatusCode)

	defer res.Body.Close()
	_, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
}

func Test_metricsHandlers_UpdateMetric(t *testing.T) {
	t.Parallel()
	type metric struct {
		mType string
		mName string
		name  string
		value string
	}

	tests := []struct {
		name           string
		metric         metric
		pathPattern    string
		reqMethod      string
		expStatus      int
		expContentType string
		wantServiceErr bool
	}{
		{
			name: "gauge",
			metric: metric{
				mType: "gauge",
				mName: "Gauge",
				name:  "test1",
				value: "42",
			},
			pathPattern:    "/update/%s/%s/%s",
			reqMethod:      http.MethodPost,
			expStatus:      http.StatusOK,
			expContentType: "text/plain",
		},
		{
			name: "counter",
			metric: metric{
				mType: "counter",
				mName: "Counter",
				name:  "test2",
				value: "42",
			},
			pathPattern:    "/update/%s/%s/%s",
			reqMethod:      http.MethodPost,
			expStatus:      http.StatusOK,
			expContentType: "text/plain",
		},
		{
			name:        "wrong method",
			pathPattern: "/update/%s/%s/%s",
			reqMethod:   http.MethodGet,
			expStatus:   http.StatusBadRequest,
		},
		{
			name: "wrong path",
			metric: metric{
				mType: "counter",
				mName: "Counter",
				value: "42",
			},
			pathPattern: "/update/%s%s%s",
			reqMethod:   http.MethodPost,
			expStatus:   http.StatusNotFound,
		},
		{
			name: "gauge not float",
			metric: metric{
				mType: "gauge",
				mName: "Gauge",
				name:  "test1",
				value: "42aa",
			},
			pathPattern: "/update/%s/%s/%s",
			reqMethod:   http.MethodPost,
			expStatus:   http.StatusBadRequest,
		},
		{
			name: "counter not int",
			metric: metric{
				mType: "counter",
				mName: "Counter",
				name:  "test1",
				value: "42aa",
			},
			pathPattern: "/update/%s/%s/%s",
			reqMethod:   http.MethodPost,
			expStatus:   http.StatusBadRequest,
		},
		{
			name: "wrong metric type",
			metric: metric{
				mType: "wrong",
				mName: "wrong",
				name:  "test1",
				value: "42aa",
			},
			pathPattern: "/update/%s/%s/%s",
			reqMethod:   http.MethodPost,
			expStatus:   http.StatusBadRequest,
		},
		{
			name: "metricsService.UpdateGauge returns error",
			metric: metric{
				mType: "gauge",
				mName: "Gauge",
				name:  "test1",
				value: "42",
			},
			pathPattern:    "/update/%s/%s/%s",
			reqMethod:      http.MethodPost,
			expStatus:      http.StatusInternalServerError,
			wantServiceErr: true,
		},
		{
			name: "metricsService.UpdateCounter returns error",
			metric: metric{
				mType: "counter",
				mName: "Counter",
				name:  "test1",
				value: "42",
			},
			pathPattern:    "/update/%s/%s/%s",
			reqMethod:      http.MethodPost,
			expStatus:      http.StatusInternalServerError,
			wantServiceErr: true,
		},
	}
	for _, tt := range tests {
		handlers, mks := getMetricsHandlersMocks()
		t.Run(tt.name, func(t *testing.T) {
			target := fmt.Sprintf(tt.pathPattern, tt.metric.mType, tt.metric.name, tt.metric.value)
			req := httptest.NewRequest(tt.reqMethod, target, nil)
			rec := httptest.NewRecorder()

			methodName := fmt.Sprintf("Update%s", tt.metric.mName)

			var err error
			if tt.wantServiceErr {
				err = errors.New("some error")
			}

			mks.metricsService.On(methodName, mock.Anything, mock.Anything).Return(err)
			handlers.UpdateMetric(rec, req)

			res := rec.Result()

			assert.Equal(t, tt.expStatus, res.StatusCode)

			defer res.Body.Close()
			_, err = io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, res.Header.Get("Content-Type"), tt.expContentType)
		})
	}
}

func Test_metricsHandlers_UpdateMetricJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		metric            any
		mockServiceMethod string
		expStatus         int
		expContentType    string
		wantServiceErr    bool
	}{
		{
			name:              "update gauge",
			metric:            generateMetric("gauge", "test1"),
			mockServiceMethod: "UpdateGauge",
			expStatus:         http.StatusOK,
			expContentType:    "application/json",
		},
		{
			name:              "update gauge, metricsService.UpdateGauge returns error",
			metric:            generateMetric("gauge", "test1"),
			mockServiceMethod: "UpdateGauge",
			expStatus:         http.StatusInternalServerError,
			expContentType:    "",
			wantServiceErr:    true,
		},
		{
			name:              "update counter",
			metric:            generateMetric("counter", "test1"),
			mockServiceMethod: "UpdateCounter",
			expStatus:         http.StatusOK,
			expContentType:    "application/json",
		},
		{
			name:              "update counter, metricsService.UpdateCounter returns error",
			metric:            generateMetric("counter", "test1"),
			mockServiceMethod: "UpdateCounter",
			expStatus:         http.StatusInternalServerError,
			expContentType:    "",
			wantServiceErr:    true,
		},
		{
			name:              "decoder returns error",
			metric:            []any{"wrong", "metric"},
			mockServiceMethod: "UpdateCounter",
			expStatus:         http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		handlers, mks := getMetricsHandlersMocks()
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			errEncode := json.NewEncoder(&buf).Encode(tt.metric)
			assert.NoError(t, errEncode)

			req := httptest.NewRequest(http.MethodPost, "/update", &buf)
			rec := httptest.NewRecorder()

			var err error
			if tt.wantServiceErr {
				err = errors.New("some error")
			}

			mks.metricsService.On(tt.mockServiceMethod, mock.Anything, mock.Anything).Return(err)
			handlers.UpdateMetricJSON(rec, req)

			res := rec.Result()

			assert.Equal(t, tt.expStatus, res.StatusCode)

			defer res.Body.Close()
			_, err = io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.expContentType, res.Header.Get("Content-Type"))
		})
	}
}

func Test_metricsHandlers_UpdateMetricsJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		metrics        any
		expStatus      int
		expContentType string
		wantServiceErr bool
	}{
		{
			name:           "update 10 metrics",
			metrics:        generateMetrics(10),
			expStatus:      http.StatusOK,
			expContentType: "application/json",
		},
		{
			name:           "update 10 metrics, metricsService.UpdateAll returns error",
			metrics:        generateMetrics(10),
			expStatus:      http.StatusInternalServerError,
			expContentType: "",
			wantServiceErr: true,
		},
		{
			name:      "decoder returns error",
			metrics:   []any{"wrong", "metric"},
			expStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		handlers, mks := getMetricsHandlersMocks()
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			errEncode := json.NewEncoder(&buf).Encode(tt.metrics)
			assert.NoError(t, errEncode)

			req := httptest.NewRequest(http.MethodPost, "/updates", &buf)
			rec := httptest.NewRecorder()

			var err error
			if tt.wantServiceErr {
				err = errors.New("some error")
			}

			mks.metricsService.On("UpdateList", mock.Anything, mock.Anything).Return(err)
			handlers.UpdateMetricsJSON(rec, req)

			res := rec.Result()

			assert.Equal(t, tt.expStatus, res.StatusCode)

			defer res.Body.Close()
			_, err = io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.expContentType, res.Header.Get("Content-Type"))
		})
	}
}

func Test_metricsHandlers_GetMetricJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		metric         any
		expStatus      int
		expContentType string
		wantServiceErr bool
	}{
		{
			name:           "get metric",
			metric:         generateMetric("gauge", "test1"),
			expStatus:      http.StatusOK,
			expContentType: "application/json",
		},
		{
			name:           "metricsService.GetMetric returns error",
			metric:         generateMetric("gauge", "test1"),
			expStatus:      http.StatusNotFound,
			expContentType: "",
			wantServiceErr: true,
		},
		{
			name:      "decoder returns error",
			metric:    []any{"wrong", "metric"},
			expStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		handlers, mks := getMetricsHandlersMocks()
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			errEncode := json.NewEncoder(&buf).Encode(tt.metric)
			assert.NoError(t, errEncode)

			req := httptest.NewRequest(http.MethodPost, "/value", &buf)
			rec := httptest.NewRecorder()

			var err error
			if tt.wantServiceErr {
				err = errors.New("some error")
			}

			mks.metricsService.On("Get", mock.Anything, mock.Anything).Return(tt.metric, err)
			handlers.GetMetricJSON(rec, req)

			res := rec.Result()

			assert.Equal(t, tt.expStatus, res.StatusCode)

			defer res.Body.Close()
			_, err = io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.expContentType, res.Header.Get("Content-Type"))
		})
	}
}

func Test_metricsHandlers_GetMetric(t *testing.T) {
	t.Parallel()
	type metric struct {
		mType string
		name  string
	}

	tests := []struct {
		name           string
		metric         metric
		pathPattern    string
		reqMethod      string
		expStatus      int
		expContentType string
		wantServiceErr bool
	}{
		{
			name: "gauge",
			metric: metric{
				mType: "gauge",
				name:  "test1",
			},
			pathPattern:    "/value/%s/%s",
			reqMethod:      http.MethodGet,
			expStatus:      http.StatusOK,
			expContentType: "text/plain",
		},
		{
			name: "counter",
			metric: metric{
				mType: "counter",
				name:  "test2",
			},
			pathPattern:    "/value/%s/%s",
			reqMethod:      http.MethodGet,
			expStatus:      http.StatusOK,
			expContentType: "text/plain",
		},
		{
			name:        "wrong method",
			pathPattern: "/value/%s/%s",
			reqMethod:   http.MethodPost,
			expStatus:   http.StatusBadRequest,
		},
		{
			name: "wrong path",
			metric: metric{
				mType: "counter",
				name:  "test3",
			},
			pathPattern: "/value/%s%s",
			reqMethod:   http.MethodGet,
			expStatus:   http.StatusNotFound,
		},
		{
			name: "wrong metric type",
			metric: metric{
				mType: "wrong",
				name:  "test1",
			},
			pathPattern: "/value/%s/%s",
			reqMethod:   http.MethodGet,
			expStatus:   http.StatusBadRequest,
		},
		{
			name: "metricsService.GetMetric returns error",
			metric: metric{
				mType: "counter",
				name:  "test1",
			},
			pathPattern:    "/value/%s/%s",
			reqMethod:      http.MethodGet,
			expStatus:      http.StatusNotFound,
			wantServiceErr: true,
		},
	}
	for _, tt := range tests {
		handlers, mks := getMetricsHandlersMocks()
		t.Run(tt.name, func(t *testing.T) {
			target := fmt.Sprintf(tt.pathPattern, tt.metric.mType, tt.metric.name)
			req := httptest.NewRequest(tt.reqMethod, target, nil)
			rec := httptest.NewRecorder()

			var err error
			if tt.wantServiceErr {
				err = errors.New("error")
			}

			mks.metricsService.On("Get", mock.Anything, mock.Anything).Return(models.Metric{Value: new(float64), Delta: new(int64)}, err)
			handlers.GetMetric(rec, req)

			res := rec.Result()

			assert.Equal(t, tt.expStatus, res.StatusCode)

			defer res.Body.Close()
			_, err = io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, res.Header.Get("Content-Type"), tt.expContentType)
		})
	}
}

func Test_metricsHandlers_AllMetrics(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name           string
		args           args
		expStatus      int
		expContentType string
	}{
		{
			name: "get all",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			expStatus:      http.StatusOK,
			expContentType: "text/html; charset=utf-8",
		},
	}
	for _, tt := range tests {
		handlers, mks := getMetricsHandlersMocks()
		t.Run(tt.name, func(t *testing.T) {
			expMetrics := generateMetrics(10)

			mks.metricsService.On("GetAll", mock.Anything).Return(expMetrics, nil)
			handlers.AllMetrics(tt.args.w, tt.args.r)

			res := tt.args.w.Result()

			assert.Equal(t, tt.expStatus, res.StatusCode)

			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, res.Header.Get("Content-Type"), tt.expContentType)
		})
	}
}

type metricsHandlersMocks struct {
	metricsService *mocks.MetricsService
}

func getMetricsHandlersMocks() (*metricsHandlers, *metricsHandlersMocks) {
	m := &metricsHandlersMocks{
		metricsService: &mocks.MetricsService{},
	}
	handlers := NewMetricsHandlers(m.metricsService)

	return handlers, m
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
