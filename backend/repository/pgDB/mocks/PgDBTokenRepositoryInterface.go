// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/jerry0420/queue-system/backend/domain"
	mock "github.com/stretchr/testify/mock"
)

// PgDBTokenRepositoryInterface is an autogenerated mock type for the PgDBTokenRepositoryInterface type
type PgDBTokenRepositoryInterface struct {
	mock.Mock
}

// CreateToken provides a mock function with given fields: ctx, token
func (_m *PgDBTokenRepositoryInterface) CreateToken(ctx context.Context, token *domain.Token) error {
	ret := _m.Called(ctx, token)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Token) error); ok {
		r0 = rf(ctx, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveTokenByToken provides a mock function with given fields: ctx, token, tokenType
func (_m *PgDBTokenRepositoryInterface) RemoveTokenByToken(ctx context.Context, token string, tokenType string) error {
	ret := _m.Called(ctx, token, tokenType)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, token, tokenType)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
