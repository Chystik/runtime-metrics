// Code generated by mockery v2.23.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// DBClient is an autogenerated mock type for the DBClient type
type DBClient struct {
	mock.Mock
}

type DBClient_Expecter struct {
	mock *mock.Mock
}

func (_m *DBClient) EXPECT() *DBClient_Expecter {
	return &DBClient_Expecter{mock: &_m.Mock}
}

// Connect provides a mock function with given fields: ctx
func (_m *DBClient) Connect(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DBClient_Connect_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Connect'
type DBClient_Connect_Call struct {
	*mock.Call
}

// Connect is a helper method to define mock.On call
//   - ctx context.Context
func (_e *DBClient_Expecter) Connect(ctx interface{}) *DBClient_Connect_Call {
	return &DBClient_Connect_Call{Call: _e.mock.On("Connect", ctx)}
}

func (_c *DBClient_Connect_Call) Run(run func(ctx context.Context)) *DBClient_Connect_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *DBClient_Connect_Call) Return(_a0 error) *DBClient_Connect_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DBClient_Connect_Call) RunAndReturn(run func(context.Context) error) *DBClient_Connect_Call {
	_c.Call.Return(run)
	return _c
}

// Disconnect provides a mock function with given fields: ctx
func (_m *DBClient) Disconnect(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DBClient_Disconnect_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Disconnect'
type DBClient_Disconnect_Call struct {
	*mock.Call
}

// Disconnect is a helper method to define mock.On call
//   - ctx context.Context
func (_e *DBClient_Expecter) Disconnect(ctx interface{}) *DBClient_Disconnect_Call {
	return &DBClient_Disconnect_Call{Call: _e.mock.On("Disconnect", ctx)}
}

func (_c *DBClient_Disconnect_Call) Run(run func(ctx context.Context)) *DBClient_Disconnect_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *DBClient_Disconnect_Call) Return(_a0 error) *DBClient_Disconnect_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DBClient_Disconnect_Call) RunAndReturn(run func(context.Context) error) *DBClient_Disconnect_Call {
	_c.Call.Return(run)
	return _c
}

// Migrate provides a mock function with given fields:
func (_m *DBClient) Migrate() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DBClient_Migrate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Migrate'
type DBClient_Migrate_Call struct {
	*mock.Call
}

// Migrate is a helper method to define mock.On call
func (_e *DBClient_Expecter) Migrate() *DBClient_Migrate_Call {
	return &DBClient_Migrate_Call{Call: _e.mock.On("Migrate")}
}

func (_c *DBClient_Migrate_Call) Run(run func()) *DBClient_Migrate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DBClient_Migrate_Call) Return(_a0 error) *DBClient_Migrate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DBClient_Migrate_Call) RunAndReturn(run func() error) *DBClient_Migrate_Call {
	_c.Call.Return(run)
	return _c
}

// Ping provides a mock function with given fields: ctx
func (_m *DBClient) Ping(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DBClient_Ping_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Ping'
type DBClient_Ping_Call struct {
	*mock.Call
}

// Ping is a helper method to define mock.On call
//   - ctx context.Context
func (_e *DBClient_Expecter) Ping(ctx interface{}) *DBClient_Ping_Call {
	return &DBClient_Ping_Call{Call: _e.mock.On("Ping", ctx)}
}

func (_c *DBClient_Ping_Call) Run(run func(ctx context.Context)) *DBClient_Ping_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *DBClient_Ping_Call) Return(_a0 error) *DBClient_Ping_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DBClient_Ping_Call) RunAndReturn(run func(context.Context) error) *DBClient_Ping_Call {
	_c.Call.Return(run)
	return _c
}

// NewDBClient creates a new instance of DBClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDBClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *DBClient {
	mock := &DBClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
