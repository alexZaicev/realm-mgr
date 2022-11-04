// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	config "github.com/alexZaicev/realm-mgr/internal/drivers/config"
	mock "github.com/stretchr/testify/mock"
)

// koanfSource is an autogenerated mock type for the koanfSource type
type koanfSource struct {
	mock.Mock
}

// Load provides a mock function with given fields: k
func (_m *koanfSource) Load(k *config.KoanfConfig) error {
	ret := _m.Called(k)

	var r0 error
	if rf, ok := ret.Get(0).(func(*config.KoanfConfig) error); ok {
		r0 = rf(k)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTnewKoanfSource interface {
	mock.TestingT
	Cleanup(func())
}

// newKoanfSource creates a new instance of koanfSource. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newKoanfSource(t mockConstructorTestingTnewKoanfSource) *koanfSource {
	mock := &koanfSource{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}