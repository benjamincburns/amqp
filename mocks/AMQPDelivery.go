// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// AMQPDelivery is an autogenerated mock type for the AMQPDelivery type
type AMQPDelivery struct {
	mock.Mock
}

// Reject provides a mock function with given fields: requeue
func (_m *AMQPDelivery) Reject(requeue bool) error {
	ret := _m.Called(requeue)

	var r0 error
	if rf, ok := ret.Get(0).(func(bool) error); ok {
		r0 = rf(requeue)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
