// Code generated by mockery. DO NOT EDIT.

package http

import (
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// Reader is an autogenerated mock type for the Reader type
type Reader struct {
	mock.Mock
}

type Reader_Expecter struct {
	mock *mock.Mock
}

func (_m *Reader) EXPECT() *Reader_Expecter {
	return &Reader_Expecter{mock: &_m.Mock}
}

// ContentType provides a mock function with given fields:
func (_m *Reader) ContentType() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ContentType")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Reader_ContentType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ContentType'
type Reader_ContentType_Call struct {
	*mock.Call
}

// ContentType is a helper method to define mock.On call
func (_e *Reader_Expecter) ContentType() *Reader_ContentType_Call {
	return &Reader_ContentType_Call{Call: _e.mock.On("ContentType")}
}

func (_c *Reader_ContentType_Call) Run(run func()) *Reader_ContentType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Reader_ContentType_Call) Return(_a0 string) *Reader_ContentType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Reader_ContentType_Call) RunAndReturn(run func() string) *Reader_ContentType_Call {
	_c.Call.Return(run)
	return _c
}

// Reader provides a mock function with given fields:
func (_m *Reader) Reader() io.Reader {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Reader")
	}

	var r0 io.Reader
	if rf, ok := ret.Get(0).(func() io.Reader); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.Reader)
		}
	}

	return r0
}

// Reader_Reader_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Reader'
type Reader_Reader_Call struct {
	*mock.Call
}

// Reader is a helper method to define mock.On call
func (_e *Reader_Expecter) Reader() *Reader_Reader_Call {
	return &Reader_Reader_Call{Call: _e.mock.On("Reader")}
}

func (_c *Reader_Reader_Call) Run(run func()) *Reader_Reader_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Reader_Reader_Call) Return(_a0 io.Reader) *Reader_Reader_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Reader_Reader_Call) RunAndReturn(run func() io.Reader) *Reader_Reader_Call {
	_c.Call.Return(run)
	return _c
}

// NewReader creates a new instance of Reader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewReader(t interface {
	mock.TestingT
	Cleanup(func())
}) *Reader {
	mock := &Reader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}