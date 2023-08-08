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
					Name: "test1",
					MetricValue: models.MetricValue{
						Gauge: 11,
					},
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

func TestUpdateCounter_WhenRepoReturnsError(t *testing.T) {
	t.Parallel()
	service, mks := getMetricsServiceMocks()

	expError := errors.New("some error")
	metric := models.Metric{
		Name: "test",
		MetricValue: models.MetricValue{
			Counter: 1,
		},
	}

	mks.repo.On("Get", mock.Anything).Return(models.Metric{}, expError)
	mks.repo.On("UpdateCounter", mock.Anything).Return()

	service.UpdateCounter(metric)
}

func TestUpdateCounter_WhenRepoReturnsResult(t *testing.T) {
	t.Parallel()
	service, mks := getMetricsServiceMocks()

	metric := models.Metric{
		Name: "test",
		MetricValue: models.MetricValue{
			Counter: 1,
		},
	}
	expMetric := metric
	expMetric.Counter += 1

	mks.repo.On("Get", mock.Anything).Return(metric, nil)
	mks.repo.On("UpdateCounter", mock.Anything).Return()

	service.UpdateCounter(expMetric)
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

func TestGetMetric_WhenRepoReturnResult(t *testing.T) {
	t.Parallel()
	var repoMock mocks.MetricsRepository
	service := New(&repoMock)

	expected := models.Metric{
		Name: "test",
		MetricValue: models.MetricValue{
			Gauge:   1,
			Counter: 2,
		}}

	repoMock.On("Get", mock.Anything).Return(expected, nil)
	actual, actualErr := service.GetMetric(expected.Name)

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
