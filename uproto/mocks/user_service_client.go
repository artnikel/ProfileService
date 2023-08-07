// Code generated by mockery v2.18.0. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	uproto "github.com/artnikel/ProfileService/uproto"
)

// UserServiceClient is an autogenerated mock type for the UserServiceClient type
type UserServiceClient struct {
	mock.Mock
}

// AddRefreshToken provides a mock function with given fields: ctx, in, opts
func (_m *UserServiceClient) AddRefreshToken(ctx context.Context, in *uproto.AddRefreshTokenRequest, opts ...grpc.CallOption) (*uproto.AddRefreshTokenResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *uproto.AddRefreshTokenResponse
	if rf, ok := ret.Get(0).(func(context.Context, *uproto.AddRefreshTokenRequest, ...grpc.CallOption) *uproto.AddRefreshTokenResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*uproto.AddRefreshTokenResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *uproto.AddRefreshTokenRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAccount provides a mock function with given fields: ctx, in, opts
func (_m *UserServiceClient) DeleteAccount(ctx context.Context, in *uproto.DeleteAccountRequest, opts ...grpc.CallOption) (*uproto.DeleteAccountResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *uproto.DeleteAccountResponse
	if rf, ok := ret.Get(0).(func(context.Context, *uproto.DeleteAccountRequest, ...grpc.CallOption) *uproto.DeleteAccountResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*uproto.DeleteAccountResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *uproto.DeleteAccountRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByLogin provides a mock function with given fields: ctx, in, opts
func (_m *UserServiceClient) GetByLogin(ctx context.Context, in *uproto.GetByLoginRequest, opts ...grpc.CallOption) (*uproto.GetByLoginResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *uproto.GetByLoginResponse
	if rf, ok := ret.Get(0).(func(context.Context, *uproto.GetByLoginRequest, ...grpc.CallOption) *uproto.GetByLoginResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*uproto.GetByLoginResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *uproto.GetByLoginRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRefreshTokenByID provides a mock function with given fields: ctx, in, opts
func (_m *UserServiceClient) GetRefreshTokenByID(ctx context.Context, in *uproto.GetRefreshTokenByIDRequest, opts ...grpc.CallOption) (*uproto.GetRefreshTokenByIDResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *uproto.GetRefreshTokenByIDResponse
	if rf, ok := ret.Get(0).(func(context.Context, *uproto.GetRefreshTokenByIDRequest, ...grpc.CallOption) *uproto.GetRefreshTokenByIDResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*uproto.GetRefreshTokenByIDResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *uproto.GetRefreshTokenByIDRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SignUp provides a mock function with given fields: ctx, in, opts
func (_m *UserServiceClient) SignUp(ctx context.Context, in *uproto.SignUpRequest, opts ...grpc.CallOption) (*uproto.SignUpResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *uproto.SignUpResponse
	if rf, ok := ret.Get(0).(func(context.Context, *uproto.SignUpRequest, ...grpc.CallOption) *uproto.SignUpResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*uproto.SignUpResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *uproto.SignUpRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUserServiceClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserServiceClient creates a new instance of UserServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserServiceClient(t mockConstructorTestingTNewUserServiceClient) *UserServiceClient {
	mock := &UserServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
