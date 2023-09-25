package metricsservice

import (
	"context"
	"errors"
	"testing"

	"github.com/Chystik/runtime-metrics/internal/models"
	"github.com/Chystik/runtime-metrics/internal/service/server/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_New(t *testing.T) {
	t.Parallel()
	service, _ := getMetricsServiceMocks()

	assert.NotNil(t, service)
}
func Test_metricsService_UpdateGauge(t *testing.T) {
	t.Parallel()
	service, mks := getMetricsServiceMocks()

	type args struct {
		metric models.Metric
	}
	tests := []struct {
		name string
		ss   MetricsService
		args args
	}{
		{
			name: "add gauge",
			ss:   service,
			args: args{
				metric: models.Metric{
					ID:    "test1",
					Value: createValue(11),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mks.repo.On("UpdateGauge", mock.Anything, tt.args.metric).Return(nil)
			err := tt.ss.UpdateGauge(context.Background(), tt.args.metric)
			assert.NoError(t, err)
		})
	}
}

func TestUpdateCounter_WhenRepoReturnsResult(t *testing.T) {
	t.Parallel()
	service, mks := getMetricsServiceMocks()

	metric := models.Metric{
		ID:    "test",
		Delta: createDelta(1),
	}
	expMetric := metric
	*expMetric.Delta += 1

	mks.repo.On("Get", mock.Anything, mock.Anything).Return(metric, nil)
	mks.repo.On("UpdateCounter", mock.Anything, mock.Anything).Return(nil)

	err := service.UpdateCounter(context.Background(), expMetric)
	assert.NoError(t, err)
}

func TestGetMetric_WhenRepoReturnResult(t *testing.T) {
	t.Parallel()
	var repoMock mocks.MetricsRepository
	service := New(&repoMock)

	expected := models.Metric{
		ID:    "test",
		Value: createValue(1),
		Delta: createDelta(2),
	}

	repoMock.On("Get", mock.Anything, mock.Anything).Return(expected, nil)
	actual, actualErr := service.GetMetric(context.Background(), expected)

	assert.NoError(t, actualErr)
	assert.Equal(t, expected, actual)
}

func TestGetMetric_WhenRepoReturnError(t *testing.T) {
	t.Parallel()
	var repoMock mocks.MetricsRepository
	service := New(&repoMock)

	expected := models.Metric{}
	expError := errors.New("some error")

	repoMock.On("Get", mock.Anything, mock.Anything).Return(expected, expError)
	actual, actualErr := service.GetMetric(context.Background(), models.Metric{ID: "some name"})

	assert.ErrorIs(t, expError, actualErr)
	assert.Equal(t, expected, actual)
}

func TestGetAllMetrics(t *testing.T) {
	t.Parallel()
	var repoMock mocks.MetricsRepository
	service := New(&repoMock)

	expected := []models.Metric{}

	repoMock.On("GetAll", mock.Anything, mock.Anything).Return(expected, nil)
	actual, err := service.GetAllMetrics(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

type metricsServiceMocks struct {
	repo *mocks.MetricsRepository
}

func getMetricsServiceMocks() (MetricsService, *metricsServiceMocks) {
	m := &metricsServiceMocks{
		repo: &mocks.MetricsRepository{},
	}
	service := New(m.repo)

	return service, m
}

func createValue(x float64) *float64 {
	return &x
}

func createDelta(x int64) *int64 {
	return &x
}
