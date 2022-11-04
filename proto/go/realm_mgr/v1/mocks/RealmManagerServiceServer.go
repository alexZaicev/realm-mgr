// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	realm_mgr_v1 "github.com/alexZaicev/realm-mgr/proto/go/realm_mgr/v1"
	mock "github.com/stretchr/testify/mock"
)

// RealmManagerServiceServer is an autogenerated mock type for the RealmManagerServiceServer type
type RealmManagerServiceServer struct {
	mock.Mock
}

// CreateRealm provides a mock function with given fields: _a0, _a1
func (_m *RealmManagerServiceServer) CreateRealm(_a0 context.Context, _a1 *realm_mgr_v1.CreateRealmRequest) (*realm_mgr_v1.CreateRealmResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *realm_mgr_v1.CreateRealmResponse
	if rf, ok := ret.Get(0).(func(context.Context, *realm_mgr_v1.CreateRealmRequest) *realm_mgr_v1.CreateRealmResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*realm_mgr_v1.CreateRealmResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *realm_mgr_v1.CreateRealmRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRealm provides a mock function with given fields: _a0, _a1
func (_m *RealmManagerServiceServer) GetRealm(_a0 context.Context, _a1 *realm_mgr_v1.GetRealmRequest) (*realm_mgr_v1.GetRealmResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *realm_mgr_v1.GetRealmResponse
	if rf, ok := ret.Get(0).(func(context.Context, *realm_mgr_v1.GetRealmRequest) *realm_mgr_v1.GetRealmResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*realm_mgr_v1.GetRealmResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *realm_mgr_v1.GetRealmRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReleaseRealm provides a mock function with given fields: _a0, _a1
func (_m *RealmManagerServiceServer) ReleaseRealm(_a0 context.Context, _a1 *realm_mgr_v1.ReleaseRealmRequest) (*realm_mgr_v1.ReleaseRealmResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *realm_mgr_v1.ReleaseRealmResponse
	if rf, ok := ret.Get(0).(func(context.Context, *realm_mgr_v1.ReleaseRealmRequest) *realm_mgr_v1.ReleaseRealmResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*realm_mgr_v1.ReleaseRealmResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *realm_mgr_v1.ReleaseRealmRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mustEmbedUnimplementedRealmManagerServiceServer provides a mock function with given fields:
func (_m *RealmManagerServiceServer) mustEmbedUnimplementedRealmManagerServiceServer() {
	_m.Called()
}

type mockConstructorTestingTNewRealmManagerServiceServer interface {
	mock.TestingT
	Cleanup(func())
}

// NewRealmManagerServiceServer creates a new instance of RealmManagerServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRealmManagerServiceServer(t mockConstructorTestingTNewRealmManagerServiceServer) *RealmManagerServiceServer {
	mock := &RealmManagerServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
