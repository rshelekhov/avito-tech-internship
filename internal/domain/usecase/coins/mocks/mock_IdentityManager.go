// Code generated by mockery v2.50.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// IdentityManager is an autogenerated mock type for the IdentityManager type
type IdentityManager struct {
	mock.Mock
}

type IdentityManager_Expecter struct {
	mock *mock.Mock
}

func (_m *IdentityManager) EXPECT() *IdentityManager_Expecter {
	return &IdentityManager_Expecter{mock: &_m.Mock}
}

// ExtractUserIDFromContext provides a mock function with given fields: ctx
func (_m *IdentityManager) ExtractUserIDFromContext(ctx context.Context) (string, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ExtractUserIDFromContext")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (string, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) string); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IdentityManager_ExtractUserIDFromContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExtractUserIDFromContext'
type IdentityManager_ExtractUserIDFromContext_Call struct {
	*mock.Call
}

// ExtractUserIDFromContext is a helper method to define mock.On call
//   - ctx context.Context
func (_e *IdentityManager_Expecter) ExtractUserIDFromContext(ctx interface{}) *IdentityManager_ExtractUserIDFromContext_Call {
	return &IdentityManager_ExtractUserIDFromContext_Call{Call: _e.mock.On("ExtractUserIDFromContext", ctx)}
}

func (_c *IdentityManager_ExtractUserIDFromContext_Call) Run(run func(ctx context.Context)) *IdentityManager_ExtractUserIDFromContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *IdentityManager_ExtractUserIDFromContext_Call) Return(_a0 string, _a1 error) *IdentityManager_ExtractUserIDFromContext_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IdentityManager_ExtractUserIDFromContext_Call) RunAndReturn(run func(context.Context) (string, error)) *IdentityManager_ExtractUserIDFromContext_Call {
	_c.Call.Return(run)
	return _c
}

// NewIdentityManager creates a new instance of IdentityManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIdentityManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *IdentityManager {
	mock := &IdentityManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
