// Code generated by mockery. DO NOT EDIT.

package runner

import (
	context "context"

	collector "github.com/kontik-pk/yandex-metrics-scraper/internal/collector"

	mock "github.com/stretchr/testify/mock"
)

// mockSaver is an autogenerated mock type for the saver type
type mockSaver struct {
	mock.Mock
}

// Restore provides a mock function with given fields: ctx
func (_m *mockSaver) Restore(ctx context.Context) ([]collector.StoredMetric, error) {
	ret := _m.Called(ctx)

	var r0 []collector.StoredMetric
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]collector.StoredMetric, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []collector.StoredMetric); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]collector.StoredMetric)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, metrics
func (_m *mockSaver) Save(ctx context.Context, metrics []collector.StoredMetric) error {
	ret := _m.Called(ctx, metrics)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []collector.StoredMetric) error); ok {
		r0 = rf(ctx, metrics)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTnewMockSaver interface {
	mock.TestingT
	Cleanup(func())
}

// newMockSaver creates a new instance of mockSaver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockSaver(t mockConstructorTestingTnewMockSaver) *mockSaver {
	mock := &mockSaver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
