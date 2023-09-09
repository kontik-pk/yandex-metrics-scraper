// Code generated by mockery. DO NOT EDIT.

package runner

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// mockHttpServer is an autogenerated mock type for the httpServer type
type mockHttpServer struct {
	mock.Mock
}

// ListenAndServe provides a mock function with given fields:
func (_m *mockHttpServer) ListenAndServe() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Shutdown provides a mock function with given fields: ctx
func (_m *mockHttpServer) Shutdown(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTnewMockHttpServer interface {
	mock.TestingT
	Cleanup(func())
}

// newMockHttpServer creates a new instance of mockHttpServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockHttpServer(t mockConstructorTestingTnewMockHttpServer) *mockHttpServer {
	mock := &mockHttpServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
