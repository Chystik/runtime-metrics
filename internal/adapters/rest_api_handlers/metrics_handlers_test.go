package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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

	mks.metricsService.On("UpdateGauge", mock.Anything).Return()
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

type metricsHandlersMocks struct {
	metricsService *mocks.MetricsService
}

func getMetricsHandlersMocks() (MetricsHandlers, *metricsHandlersMocks) {
	m := &metricsHandlersMocks{
		metricsService: &mocks.MetricsService{},
	}
	handlers := NewMetricHandlers(m.metricsService)

	return handlers, m
}
