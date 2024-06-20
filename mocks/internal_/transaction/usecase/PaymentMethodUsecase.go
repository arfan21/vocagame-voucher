// Code generated by mockery. DO NOT EDIT.

package transactionuc

import (
	context "context"

	model "github.com/arfan21/vocagame/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// PaymentMethodUsecase is an autogenerated mock type for the PaymentMethodUsecase type
type PaymentMethodUsecase struct {
	mock.Mock
}

type PaymentMethodUsecase_Expecter struct {
	mock *mock.Mock
}

func (_m *PaymentMethodUsecase) EXPECT() *PaymentMethodUsecase_Expecter {
	return &PaymentMethodUsecase_Expecter{mock: &_m.Mock}
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *PaymentMethodUsecase) GetByID(ctx context.Context, id int) (model.PaymentMethodResponse, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 model.PaymentMethodResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (model.PaymentMethodResponse, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) model.PaymentMethodResponse); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(model.PaymentMethodResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PaymentMethodUsecase_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type PaymentMethodUsecase_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id int
func (_e *PaymentMethodUsecase_Expecter) GetByID(ctx interface{}, id interface{}) *PaymentMethodUsecase_GetByID_Call {
	return &PaymentMethodUsecase_GetByID_Call{Call: _e.mock.On("GetByID", ctx, id)}
}

func (_c *PaymentMethodUsecase_GetByID_Call) Run(run func(ctx context.Context, id int)) *PaymentMethodUsecase_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *PaymentMethodUsecase_GetByID_Call) Return(_a0 model.PaymentMethodResponse, _a1 error) *PaymentMethodUsecase_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PaymentMethodUsecase_GetByID_Call) RunAndReturn(run func(context.Context, int) (model.PaymentMethodResponse, error)) *PaymentMethodUsecase_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// NewPaymentMethodUsecase creates a new instance of PaymentMethodUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPaymentMethodUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *PaymentMethodUsecase {
	mock := &PaymentMethodUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
