package main

import (
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type MockExpectTestsGenerator struct{}

func (MockExpectTestsGenerator) Checker() checkers.Checker {
	return checkers.NewMockExpect()
}

func (g MockExpectTestsGenerator) TemplateData() any {
	var (
		checker = g.Checker().Name()
		report  = checker + `: use u\\.EXPECT\\(\\)\\.%s\\(\\.\\.\\.\\)`
	)

	return struct {
		CheckerName CheckerName
		Report      string
	}{
		CheckerName: CheckerName(checker),
		Report:      report,
	}
}

func (MockExpectTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("MockExpectTestsGenerator.ErroredTemplate").
		Parse(mockExpectTestTmpl))
}

func (MockExpectTestsGenerator) GoldenTemplate() Executor {
	return template.Must(template.New("MockExpectTestsGenerator.GoldenTemplate").
		Parse(mockExpectGoldenTmpl))
}

const mockExpectTestTmpl = boilerPlate + `
func {{ .CheckerName.AsTestName }}(t *testing.T) {
	u := NewMockUserIFace(t)

	// Invalid.
	{
		u.On("CreateUser", mock.Anything, User{}).Return(nil)    // want "{{ printf $.Report "CreateUser" }}"
		u.On("GetUser", t.Context(), "test").Return(User{}, nil) // want "{{ printf $.Report "GetUser" }}"
		u.On("Void")                                             // want "{{ printf $.Report "Void" }}"
		u.On("Void").Once()                                      // want "{{ printf $.Report "Void" }}"
		u.On("CountUsers").Return(123)                           // want "{{ printf $.Report "CountUsers" }}"
	}

	// Ignored.
	{
		u.On("", mock.Anything, User{}).Return(nil)
		u.On("DoesNotExist", mock.Anything, User{}, 123).Return(nil)
	}
}`

const mockExpectGoldenTmpl = boilerPlate + `
func {{ .CheckerName.AsTestName }}(t *testing.T) {
	u := NewMockUserIFace(t)

	// Invalid.
	{
		u.EXPECT().CreateUser(mock.Anything, User{}).Return(nil)    // want "{{ printf $.Report "CreateUser" }}"
		u.EXPECT().GetUser(t.Context(), "test").Return(User{}, nil) // want "{{ printf $.Report "GetUser" }}"
		u.EXPECT().Void()                                           // want "{{ printf $.Report "Void" }}"
		u.EXPECT().Void().Once()                                    // want "{{ printf $.Report "Void" }}"
		u.EXPECT().CountUsers().Return(123)                         // want "{{ printf $.Report "CountUsers" }}"
	}

	// Ignored.
	{
		u.On("", mock.Anything, User{}).Return(nil)
		u.On("DoesNotExist", mock.Anything, User{}, 123).Return(nil)
	}
}
`

const boilerPlate = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockUserIFace struct {
	mock.Mock
}

type MockUserIFace_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUserIFace) EXPECT() *MockUserIFace_Expecter {
	return &MockUserIFace_Expecter{mock: &_m.Mock}
}

func (_m *MockUserIFace) CountUsers() int {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for CountUsers")
	}

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

type MockUserIFace_CountUsers_Call struct {
	*mock.Call
}

func (_e *MockUserIFace_Expecter) CountUsers() *MockUserIFace_CountUsers_Call {
	return &MockUserIFace_CountUsers_Call{Call: _e.mock.On("CountUsers")}
}

func (_c *MockUserIFace_CountUsers_Call) Run(run func()) *MockUserIFace_CountUsers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockUserIFace_CountUsers_Call) Return(_a0 int) *MockUserIFace_CountUsers_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockUserIFace_CountUsers_Call) RunAndReturn(run func() int) *MockUserIFace_CountUsers_Call {
	_c.Call.Return(run)
	return _c
}

func (_m *MockUserIFace) CreateUser(_a0 context.Context, _a1 User) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, User) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type MockUserIFace_CreateUser_Call struct {
	*mock.Call
}

func (_e *MockUserIFace_Expecter) CreateUser(_a0 interface{}, _a1 interface{}) *MockUserIFace_CreateUser_Call {
	return &MockUserIFace_CreateUser_Call{Call: _e.mock.On("CreateUser", _a0, _a1)}
}

func (_c *MockUserIFace_CreateUser_Call) Run(run func(_a0 context.Context, _a1 User)) *MockUserIFace_CreateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(User))
	})
	return _c
}

func (_c *MockUserIFace_CreateUser_Call) Return(_a0 error) *MockUserIFace_CreateUser_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockUserIFace_CreateUser_Call) RunAndReturn(run func(context.Context, User) error) *MockUserIFace_CreateUser_Call {
	_c.Call.Return(run)
	return _c
}

func (_m *MockUserIFace) GetUser(ctx context.Context, name string) (User, error) {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for GetUser")
	}

	var r0 User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (User, error)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) User); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Get(0).(User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type MockUserIFace_GetUser_Call struct {
	*mock.Call
}

func (_e *MockUserIFace_Expecter) GetUser(ctx interface{}, name interface{}) *MockUserIFace_GetUser_Call {
	return &MockUserIFace_GetUser_Call{Call: _e.mock.On("GetUser", ctx, name)}
}

func (_c *MockUserIFace_GetUser_Call) Run(run func(ctx context.Context, name string)) *MockUserIFace_GetUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockUserIFace_GetUser_Call) Return(_a0 User, _a1 error) *MockUserIFace_GetUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserIFace_GetUser_Call) RunAndReturn(run func(context.Context, string) (User, error)) *MockUserIFace_GetUser_Call {
	_c.Call.Return(run)
	return _c
}

func (_m *MockUserIFace) Void() {
	_m.Called()
}

type MockUserIFace_Void_Call struct {
	*mock.Call
}

func (_e *MockUserIFace_Expecter) Void() *MockUserIFace_Void_Call {
	return &MockUserIFace_Void_Call{Call: _e.mock.On("Void")}
}

func (_c *MockUserIFace_Void_Call) Run(run func()) *MockUserIFace_Void_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockUserIFace_Void_Call) Return() *MockUserIFace_Void_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockUserIFace_Void_Call) RunAndReturn(run func()) *MockUserIFace_Void_Call {
	_c.Run(run)
	return _c
}

func NewMockUserIFace(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUserIFace {
	mock := &MockUserIFace{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

type User struct {
	Name string
	Age  int
}
`
