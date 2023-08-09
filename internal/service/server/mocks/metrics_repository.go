// Code generated by mockery v2.23.4. DO NOT EDIT.

package mocks

import (
	models "github.com/Chystik/runtime-metrics/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// MetricsRepository is an autogenerated mock type for the MetricsRepository type
type MetricsRepository struct {
	mock.Mock
}

type MetricsRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MetricsRepository) EXPECT() *MetricsRepository_Expecter {
	return &MetricsRepository_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: name
func (_m *MetricsRepository) Get(name string) (models.Metric, error) {
	ret := _m.Called(name)

	var r0 models.Metric
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (models.Metric, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) models.Metric); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(models.Metric)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MetricsRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MetricsRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - name string
func (_e *MetricsRepository_Expecter) Get(name interface{}) *MetricsRepository_Get_Call {
	return &MetricsRepository_Get_Call{Call: _e.mock.On("Get", name)}
}

func (_c *MetricsRepository_Get_Call) Run(run func(name string)) *MetricsRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MetricsRepository_Get_Call) Return(_a0 models.Metric, _a1 error) *MetricsRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MetricsRepository_Get_Call) RunAndReturn(run func(string) (models.Metric, error)) *MetricsRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with given fields:
func (_m *MetricsRepository) GetAll() []models.Metric {
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

// MetricsRepository_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type MetricsRepository_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
func (_e *MetricsRepository_Expecter) GetAll() *MetricsRepository_GetAll_Call {
	return &MetricsRepository_GetAll_Call{Call: _e.mock.On("GetAll")}
}

func (_c *MetricsRepository_GetAll_Call) Run(run func()) *MetricsRepository_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MetricsRepository_GetAll_Call) Return(_a0 []models.Metric) *MetricsRepository_GetAll_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MetricsRepository_GetAll_Call) RunAndReturn(run func() []models.Metric) *MetricsRepository_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateCounter provides a mock function with given fields: _a0
func (_m *MetricsRepository) UpdateCounter(_a0 models.Metric) {
	_m.Called(_a0)
}

// MetricsRepository_UpdateCounter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateCounter'
type MetricsRepository_UpdateCounter_Call struct {
	*mock.Call
}

// UpdateCounter is a helper method to define mock.On call
//   - _a0 models.Metric
func (_e *MetricsRepository_Expecter) UpdateCounter(_a0 interface{}) *MetricsRepository_UpdateCounter_Call {
	return &MetricsRepository_UpdateCounter_Call{Call: _e.mock.On("UpdateCounter", _a0)}
}

func (_c *MetricsRepository_UpdateCounter_Call) Run(run func(_a0 models.Metric)) *MetricsRepository_UpdateCounter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(models.Metric))
	})
	return _c
}

func (_c *MetricsRepository_UpdateCounter_Call) Return() *MetricsRepository_UpdateCounter_Call {
	_c.Call.Return()
	return _c
}

func (_c *MetricsRepository_UpdateCounter_Call) RunAndReturn(run func(models.Metric)) *MetricsRepository_UpdateCounter_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateGauge provides a mock function with given fields: _a0
func (_m *MetricsRepository) UpdateGauge(_a0 models.Metric) {
	_m.Called(_a0)
}

// MetricsRepository_UpdateGauge_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateGauge'
type MetricsRepository_UpdateGauge_Call struct {
	*mock.Call
}

// UpdateGauge is a helper method to define mock.On call
//   - _a0 models.Metric
func (_e *MetricsRepository_Expecter) UpdateGauge(_a0 interface{}) *MetricsRepository_UpdateGauge_Call {
	return &MetricsRepository_UpdateGauge_Call{Call: _e.mock.On("UpdateGauge", _a0)}
}

func (_c *MetricsRepository_UpdateGauge_Call) Run(run func(_a0 models.Metric)) *MetricsRepository_UpdateGauge_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(models.Metric))
	})
	return _c
}

func (_c *MetricsRepository_UpdateGauge_Call) Return() *MetricsRepository_UpdateGauge_Call {
	_c.Call.Return()
	return _c
}

func (_c *MetricsRepository_UpdateGauge_Call) RunAndReturn(run func(models.Metric)) *MetricsRepository_UpdateGauge_Call {
	_c.Call.Return(run)
	return _c
}

// NewMetricsRepository creates a new instance of MetricsRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMetricsRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MetricsRepository {
	mock := &MetricsRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
