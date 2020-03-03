// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import amqp "github.com/streadway/amqp"

import mock "github.com/stretchr/testify/mock"

// AMQPChannel is an autogenerated mock type for the AMQPChannel type
type AMQPChannel struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *AMQPChannel) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Consume provides a mock function with given fields: queue, consumer, autoAck, exclusive, noLocal, noWait, args
func (_m *AMQPChannel) Consume(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	ret := _m.Called(queue, consumer, autoAck, exclusive, noLocal, noWait, args)

	var r0 <-chan amqp.Delivery
	if rf, ok := ret.Get(0).(func(string, string, bool, bool, bool, bool, amqp.Table) <-chan amqp.Delivery); ok {
		r0 = rf(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan amqp.Delivery)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, bool, bool, bool, bool, amqp.Table) error); ok {
		r1 = rf(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExchangeBind provides a mock function with given fields: destination, key, source, noWait, args
func (_m *AMQPChannel) ExchangeBind(destination string, key string, source string, noWait bool, args amqp.Table) error {
	ret := _m.Called(destination, key, source, noWait, args)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, bool, amqp.Table) error); ok {
		r0 = rf(destination, key, source, noWait, args)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ExchangeDeclare provides a mock function with given fields: name, kind, durable, autoDelete, internal, noWait, args
func (_m *AMQPChannel) ExchangeDeclare(name string, kind string, durable bool, autoDelete bool, internal bool, noWait bool, args amqp.Table) error {
	ret := _m.Called(name, kind, durable, autoDelete, internal, noWait, args)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, bool, bool, bool, bool, amqp.Table) error); ok {
		r0 = rf(name, kind, durable, autoDelete, internal, noWait, args)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Publish provides a mock function with given fields: exchange, key, mandatory, immediate, msg
func (_m *AMQPChannel) Publish(exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error {
	ret := _m.Called(exchange, key, mandatory, immediate, msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, bool, bool, amqp.Publishing) error); ok {
		r0 = rf(exchange, key, mandatory, immediate, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// QueueBind provides a mock function with given fields: name, key, exchange, noWait, args
func (_m *AMQPChannel) QueueBind(name string, key string, exchange string, noWait bool, args amqp.Table) error {
	ret := _m.Called(name, key, exchange, noWait, args)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, bool, amqp.Table) error); ok {
		r0 = rf(name, key, exchange, noWait, args)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// QueueDeclare provides a mock function with given fields: name, durable, autoDelete, exclusive, noWait, args
func (_m *AMQPChannel) QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error) {
	ret := _m.Called(name, durable, autoDelete, exclusive, noWait, args)

	var r0 amqp.Queue
	if rf, ok := ret.Get(0).(func(string, bool, bool, bool, bool, amqp.Table) amqp.Queue); ok {
		r0 = rf(name, durable, autoDelete, exclusive, noWait, args)
	} else {
		r0 = ret.Get(0).(amqp.Queue)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, bool, bool, bool, bool, amqp.Table) error); ok {
		r1 = rf(name, durable, autoDelete, exclusive, noWait, args)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Tx provides a mock function with given fields:
func (_m *AMQPChannel) Tx() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TxCommit provides a mock function with given fields:
func (_m *AMQPChannel) TxCommit() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TxRollback provides a mock function with given fields:
func (_m *AMQPChannel) TxRollback() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
