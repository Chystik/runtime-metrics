// Code generated by mockery v2.23.4. DO NOT EDIT.

package mocks

import (
	models "github.com/Chystik/runtime-metrics/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// MetricsService is an autogenerated mock type for the MetricsService type
type MetricsService struct {
	mock.Mock
}

// UpdateCounter provides a mock function with given fields: _a0
func (_m *MetricsService) UpdateCounter(_a0 models.Metric) {
	_m.Called(_a0)
}

// UpdateGauge provides a mock function with given fields: _a0
func (_m *MetricsService) UpdateGauge(_a0 models.Metric) {
	_m.Called(_a0)
}

// NewMetricsService creates a new instance of MetricsService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMetricsService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MetricsService {
	mock := &MetricsService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
