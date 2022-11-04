// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	uuid "github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

// Generator is an autogenerated mock type for the Generator type
type Generator struct {
	mock.Mock
}

// New provides a mock function with given fields:
func (_m *Generator) New() (uuid.UUID, error) {
	ret := _m.Called()

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func() uuid.UUID); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewGenerator interface {
	mock.TestingT
	Cleanup(func())
}

// NewGenerator creates a new instance of Generator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGenerator(t mockConstructorTestingTNewGenerator) *Generator {
	mock := &Generator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}