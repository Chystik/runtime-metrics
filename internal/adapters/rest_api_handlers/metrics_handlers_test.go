package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Chystik/runtime-metrics/internal/adapters"
	"github.com/Chystik/runtime-metrics/internal/models"
	"github.com/Chystik/runtime-metrics/internal/service/server/mocks"

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

	mks.metricsService.EXPECT().UpdateGauge(mock.Anything).Return()
	handlers.UpdateMetric(rec, req)

	res := rec.Result()

	assert.Equal(t, res.StatusCode, expStatus)

	defer res.Body.Close()
	_, err := io.ReadAll(res.Body)

	require.NoError(t, err)
	assert.Equal(t, res.Header.Get("Content-Type"), expContextType)
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
		expContextType string
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
			expContextType: "text/plain",
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
			expContextType: "text/plain",
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
	}
	for _, tt := range tests {
		handlers, mks := getMetricsHandlersMocks()
		t.Run(tt.name, func(t *testing.T) {
			target := fmt.Sprintf(tt.pathPattern, tt.metric.mType, tt.metric.name, tt.metric.value)
			req := httptest.NewRequest(tt.reqMethod, target, nil)
			rec := httptest.NewRecorder()

			methodName := fmt.Sprintf("Update%s", tt.metric.mName)
			mks.metricsService.On(methodName, mock.Anything).Return()
			handlers.UpdateMetric(rec, req)

			res := rec.Result()

			assert.Equal(t, tt.expStatus, res.StatusCode)

			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, res.Header.Get("Content-Type"), tt.expContextType)
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
	}
	for _, tt := range tests {
		handlers, mks := getMetricsHandlersMocks()
		t.Run(tt.name, func(t *testing.T) {
			target := fmt.Sprintf(tt.pathPattern, tt.metric.mType, tt.metric.name)
			req := httptest.NewRequest(tt.reqMethod, target, nil)
			rec := httptest.NewRecorder()

			mks.metricsService.On("GetMetric", mock.Anything).Return(models.Metric{Value: new(float64), Delta: new(int64)}, nil)
			handlers.GetMetric(rec, req)

			res := rec.Result()

			assert.Equal(t, tt.expStatus, res.StatusCode)

			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)

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
			mks.metricsService.On("GetAllMetrics").Return([]models.Metric{})
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

func getMetricsHandlersMocks() (adapters.MetricsHandlers, *metricsHandlersMocks) {
	m := &metricsHandlersMocks{
		metricsService: &mocks.MetricsService{},
	}
	handlers := NewMetricsHandlers(m.metricsService)

	return handlers, m
}
