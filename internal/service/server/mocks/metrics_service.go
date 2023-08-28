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

type MetricsService_Expecter struct {
	mock *mock.Mock
}

func (_m *MetricsService) EXPECT() *MetricsService_Expecter {
	return &MetricsService_Expecter{mock: &_m.Mock}
}

// GetAllMetrics provides a mock function with given fields:
func (_m *MetricsService) GetAllMetrics() []models.Metric {
	ret := _m.Called()

	var r0 []models.Metric
	if rf, ok := ret.Get(0).(func() []models.Metric); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Metric)
		}
	}

	return r0
}

// MetricsService_GetAllMetrics_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllMetrics'
type MetricsService_GetAllMetrics_Call struct {
	*mock.Call
}

// GetAllMetrics is a helper method to define mock.On call
func (_e *MetricsService_Expecter) GetAllMetrics() *MetricsService_GetAllMetrics_Call {
	return &MetricsService_GetAllMetrics_Call{Call: _e.mock.On("GetAllMetrics")}
}

func (_c *MetricsService_GetAllMetrics_Call) Run(run func()) *MetricsService_GetAllMetrics_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MetricsService_GetAllMetrics_Call) Return(_a0 []models.Metric) *MetricsService_GetAllMetrics_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MetricsService_GetAllMetrics_Call) RunAndReturn(run func() []models.Metric) *MetricsService_GetAllMetrics_Call {
	_c.Call.Return(run)
	return _c
}

// GetMetric provides a mock function with given fields: Id
func (_m *MetricsService) GetMetric(Id string) (models.Metric, error) {
	ret := _m.Called(Id)

	var r0 models.Metric
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (models.Metric, error)); ok {
		return rf(Id)
	}
	if rf, ok := ret.Get(0).(func(string) models.Metric); ok {
		r0 = rf(Id)
	} else {
		r0 = ret.Get(0).(models.Metric)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(Id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MetricsService_GetMetric_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMetric'
type MetricsService_GetMetric_Call struct {
	*mock.Call
}

// GetMetric is a helper method to define mock.On call
//   - Id string
func (_e *MetricsService_Expecter) GetMetric(Id interface{}) *MetricsService_GetMetric_Call {
	return &MetricsService_GetMetric_Call{Call: _e.mock.On("GetMetric", Id)}
}

func (_c *MetricsService_GetMetric_Call) Run(run func(Id string)) *MetricsService_GetMetric_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MetricsService_GetMetric_Call) Return(_a0 models.Metric, _a1 error) *MetricsService_GetMetric_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MetricsService_GetMetric_Call) RunAndReturn(run func(string) (models.Metric, error)) *MetricsService_GetMetric_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateCounter provides a mock function with given fields: _a0
func (_m *MetricsService) UpdateCounter(_a0 models.Metric) {
	_m.Called(_a0)
}

// MetricsService_UpdateCounter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateCounter'
type MetricsService_UpdateCounter_Call struct {
	*mock.Call
}

// UpdateCounter is a helper method to define mock.On call
//   - _a0 models.Metric
func (_e *MetricsService_Expecter) UpdateCounter(_a0 interface{}) *MetricsService_UpdateCounter_Call {
	return &MetricsService_UpdateCounter_Call{Call: _e.mock.On("UpdateCounter", _a0)}
}

func (_c *MetricsService_UpdateCounter_Call) Run(run func(_a0 models.Metric)) *MetricsService_UpdateCounter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(models.Metric))
	})
	return _c
}

func (_c *MetricsService_UpdateCounter_Call) Return() *MetricsService_UpdateCounter_Call {
	_c.Call.Return()
	return _c
}

func (_c *MetricsService_UpdateCounter_Call) RunAndReturn(run func(models.Metric)) *MetricsService_UpdateCounter_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateGauge provides a mock function with given fields: _a0
func (_m *MetricsService) UpdateGauge(_a0 models.Metric) {
	_m.Called(_a0)
}

// MetricsService_UpdateGauge_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateGauge'
type MetricsService_UpdateGauge_Call struct {
	*mock.Call
}

// UpdateGauge is a helper method to define mock.On call
//   - _a0 models.Metric
func (_e *MetricsService_Expecter) UpdateGauge(_a0 interface{}) *MetricsService_UpdateGauge_Call {
	return &MetricsService_UpdateGauge_Call{Call: _e.mock.On("UpdateGauge", _a0)}
}

func (_c *MetricsService_UpdateGauge_Call) Run(run func(_a0 models.Metric)) *MetricsService_UpdateGauge_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(models.Metric))
	})
	return _c
}

func (_c *MetricsService_UpdateGauge_Call) Return() *MetricsService_UpdateGauge_Call {
	_c.Call.Return()
	return _c
}

func (_c *MetricsService_UpdateGauge_Call) RunAndReturn(run func(models.Metric)) *MetricsService_UpdateGauge_Call {
	_c.Call.Return(run)
	return _c
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
