// Code generated by mockery v2.23.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// PostgresClient is an autogenerated mock type for the PostgresClient type
type PostgresClient struct {
	mock.Mock
}

type PostgresClient_Expecter struct {
	mock *mock.Mock
}

func (_m *PostgresClient) EXPECT() *PostgresClient_Expecter {
	return &PostgresClient_Expecter{mock: &_m.Mock}
}

// Connect provides a mock function with given fields: ctx
func (_m *PostgresClient) Connect(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PostgresClient_Connect_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Connect'
type PostgresClient_Connect_Call struct {
	*mock.Call
}

// Connect is a helper method to define mock.On call
//   - ctx context.Context
func (_e *PostgresClient_Expecter) Connect(ctx interface{}) *PostgresClient_Connect_Call {
	return &PostgresClient_Connect_Call{Call: _e.mock.On("Connect", ctx)}
}

func (_c *PostgresClient_Connect_Call) Run(run func(ctx context.Context)) *PostgresClient_Connect_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *PostgresClient_Connect_Call) Return(_a0 error) *PostgresClient_Connect_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PostgresClient_Connect_Call) RunAndReturn(run func(context.Context) error) *PostgresClient_Connect_Call {
	_c.Call.Return(run)
	return _c
}

// Disconnect provides a mock function with given fields: ctx
func (_m *PostgresClient) Disconnect(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PostgresClient_Disconnect_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Disconnect'
type PostgresClient_Disconnect_Call struct {
	*mock.Call
}

// Disconnect is a helper method to define mock.On call
//   - ctx context.Context
func (_e *PostgresClient_Expecter) Disconnect(ctx interface{}) *PostgresClient_Disconnect_Call {
	return &PostgresClient_Disconnect_Call{Call: _e.mock.On("Disconnect", ctx)}
}

func (_c *PostgresClient_Disconnect_Call) Run(run func(ctx context.Context)) *PostgresClient_Disconnect_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *PostgresClient_Disconnect_Call) Return(_a0 error) *PostgresClient_Disconnect_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PostgresClient_Disconnect_Call) RunAndReturn(run func(context.Context) error) *PostgresClient_Disconnect_Call {
	_c.Call.Return(run)
	return _c
}

// Migrate provides a mock function with given fields:
func (_m *PostgresClient) Migrate() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PostgresClient_Migrate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Migrate'
type PostgresClient_Migrate_Call struct {
	*mock.Call
}

// Migrate is a helper method to define mock.On call
func (_e *PostgresClient_Expecter) Migrate() *PostgresClient_Migrate_Call {
	return &PostgresClient_Migrate_Call{Call: _e.mock.On("Migrate")}
}

func (_c *PostgresClient_Migrate_Call) Run(run func()) *PostgresClient_Migrate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *PostgresClient_Migrate_Call) Return(_a0 error) *PostgresClient_Migrate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PostgresClient_Migrate_Call) RunAndReturn(run func() error) *PostgresClient_Migrate_Call {
	_c.Call.Return(run)
	return _c
}

// Ping provides a mock function with given fields: ctx
func (_m *PostgresClient) Ping(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PostgresClient_Ping_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Ping'
type PostgresClient_Ping_Call struct {
	*mock.Call
}

// Ping is a helper method to define mock.On call
//   - ctx context.Context
func (_e *PostgresClient_Expecter) Ping(ctx interface{}) *PostgresClient_Ping_Call {
	return &PostgresClient_Ping_Call{Call: _e.mock.On("Ping", ctx)}
}

func (_c *PostgresClient_Ping_Call) Run(run func(ctx context.Context)) *PostgresClient_Ping_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *PostgresClient_Ping_Call) Return(_a0 error) *PostgresClient_Ping_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PostgresClient_Ping_Call) RunAndReturn(run func(context.Context) error) *PostgresClient_Ping_Call {
	_c.Call.Return(run)
	return _c
}

// NewPostgresClient creates a new instance of PostgresClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPostgresClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *PostgresClient {
	mock := &PostgresClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
