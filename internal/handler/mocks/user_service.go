// Code generated by mockery v2.18.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/artnikel/ProfileService/internal/model"

	service "github.com/artnikel/ProfileService/internal/service"

	uuid "github.com/google/uuid"
)

// UserService is an autogenerated mock type for the UserService type
type UserService struct {
	mock.Mock
}

// DeleteAccount provides a mock function with given fields: ctx, id
func (_m *UserService) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Login provides a mock function with given fields: ctx, user
func (_m *UserService) Login(ctx context.Context, user *model.User) (*service.TokenPair, error) {
	ret := _m.Called(ctx, user)

	var r0 *service.TokenPair
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) *service.TokenPair); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*service.TokenPair)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *model.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Refresh provides a mock function with given fields: ctx, tokenPair
func (_m *UserService) Refresh(ctx context.Context, tokenPair service.TokenPair) (*service.TokenPair, error) {
	ret := _m.Called(ctx, tokenPair)

	var r0 *service.TokenPair
	if rf, ok := ret.Get(0).(func(context.Context, service.TokenPair) *service.TokenPair); ok {
		r0 = rf(ctx, tokenPair)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*service.TokenPair)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, service.TokenPair) error); ok {
		r1 = rf(ctx, tokenPair)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SignUp provides a mock function with given fields: ctx, user
func (_m *UserService) SignUp(ctx context.Context, user *model.User) error {
	ret := _m.Called(ctx, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewUserService interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserService creates a new instance of UserService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserService(t mockConstructorTestingTNewUserService) *UserService {
	mock := &UserService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
