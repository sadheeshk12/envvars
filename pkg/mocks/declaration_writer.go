// Code generated by mockery v1.0.0. DO NOT EDIT.
package mocks

import envvars "github.com/flemay/envvars/pkg/envvars"
import mock "github.com/stretchr/testify/mock"

// DeclarationWriter is an autogenerated mock type for the DeclarationWriter type
type DeclarationWriter struct {
	mock.Mock
}

// Write provides a mock function with given fields: d, overwrite
func (_m *DeclarationWriter) Write(d *envvars.Declaration, overwrite bool) error {
	ret := _m.Called(d, overwrite)

	var r0 error
	if rf, ok := ret.Get(0).(func(*envvars.Declaration, bool) error); ok {
		r0 = rf(d, overwrite)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
