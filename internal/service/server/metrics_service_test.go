package metricsservice

import (
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
					Id:    "test1",
					Value: createValue(11),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mks.repo.On("UpdateGauge", tt.args.metric).Return()
			tt.ss.UpdateGauge(tt.args.metric)
		})
	}
}

func TestUpdateCounter_WhenRepoReturnsResult(t *testing.T) {
	t.Parallel()
	service, mks := getMetricsServiceMocks()

	metric := models.Metric{
		Id:    "test",
		Delta: createDelta(1),
	}
	expMetric := metric
	*expMetric.Delta += 1

	mks.repo.On("Get", mock.Anything).Return(metric, nil)
	mks.repo.On("UpdateCounter", mock.Anything).Return()

	service.UpdateCounter(expMetric)
}

func TestGetMetric_WhenRepoReturnResult(t *testing.T) {
	t.Parallel()
	var repoMock mocks.MetricsRepository
	service := New(&repoMock)

	expected := models.Metric{
		Id:    "test",
		Value: createValue(1),
		Delta: createDelta(2),
	}

	repoMock.On("Get", mock.Anything).Return(expected, nil)
	actual, actualErr := service.GetMetric(expected.Id)

	assert.NoError(t, actualErr)
	assert.Equal(t, expected, actual)
}

func TestGetMetric_WhenRepoReturnError(t *testing.T) {
	t.Parallel()
	var repoMock mocks.MetricsRepository
	service := New(&repoMock)

	expected := models.Metric{}
	expError := errors.New("some error")

	repoMock.On("Get", mock.Anything).Return(expected, expError)
	actual, actualErr := service.GetMetric("some name")

	assert.ErrorIs(t, expError, actualErr)
	assert.Equal(t, expected, actual)
}

func TestGetAllMetrics(t *testing.T) {
	t.Parallel()
	var repoMock mocks.MetricsRepository
	service := New(&repoMock)

	expected := []models.Metric{}

	repoMock.On("GetAll", mock.Anything).Return(expected)
	actual := service.GetAllMetrics()

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
