// Code generated by mockery v1.0.0. DO NOT EDIT.
package mocks

import mock "github.com/stretchr/testify/mock"

// Auth is an autogenerated mock type for the Auth type
type Auth struct {
	mock.Mock
}

// LoginK8s provides a mock function with given fields: role, jwt, path
func (_m *Auth) LoginK8s(role string, jwt string, path string) (string, error) {
	ret := _m.Called(role, jwt, path)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, string) string); ok {
		r0 = rf(role, jwt, path)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(role, jwt, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
