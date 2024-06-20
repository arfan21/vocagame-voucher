// Code generated by mockery. DO NOT EDIT.

package transactionuc

import (
	context "context"

	entity "github.com/arfan21/vocagame/internal/entity"
	mock "github.com/stretchr/testify/mock"

	pgx "github.com/jackc/pgx/v5"

	transactionrepo "github.com/arfan21/vocagame/internal/transaction/repository"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

type Repository_Expecter struct {
	mock *mock.Mock
}

func (_m *Repository) EXPECT() *Repository_Expecter {
	return &Repository_Expecter{mock: &_m.Mock}
}

// Begin provides a mock function with given fields: ctx
func (_m *Repository) Begin(ctx context.Context) (pgx.Tx, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Begin")
	}

	var r0 pgx.Tx
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (pgx.Tx, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) pgx.Tx); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pgx.Tx)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_Begin_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Begin'
type Repository_Begin_Call struct {
	*mock.Call
}

// Begin is a helper method to define mock.On call
//   - ctx context.Context
func (_e *Repository_Expecter) Begin(ctx interface{}) *Repository_Begin_Call {
	return &Repository_Begin_Call{Call: _e.mock.On("Begin", ctx)}
}

func (_c *Repository_Begin_Call) Run(run func(ctx context.Context)) *Repository_Begin_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Repository_Begin_Call) Return(tx pgx.Tx, err error) *Repository_Begin_Call {
	_c.Call.Return(tx, err)
	return _c
}

func (_c *Repository_Begin_Call) RunAndReturn(run func(context.Context) (pgx.Tx, error)) *Repository_Begin_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: ctx, data
func (_m *Repository) Create(ctx context.Context, data entity.Transaction) error {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Transaction) error); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Repository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type Repository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - data entity.Transaction
func (_e *Repository_Expecter) Create(ctx interface{}, data interface{}) *Repository_Create_Call {
	return &Repository_Create_Call{Call: _e.mock.On("Create", ctx, data)}
}

func (_c *Repository_Create_Call) Run(run func(ctx context.Context, data entity.Transaction)) *Repository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(entity.Transaction))
	})
	return _c
}

func (_c *Repository_Create_Call) Return(err error) *Repository_Create_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *Repository_Create_Call) RunAndReturn(run func(context.Context, entity.Transaction) error) *Repository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// GetByPhoneNumber provides a mock function with given fields: ctx, phoneNumber
func (_m *Repository) GetByPhoneNumber(ctx context.Context, phoneNumber string) ([]entity.Transaction, error) {
	ret := _m.Called(ctx, phoneNumber)

	if len(ret) == 0 {
		panic("no return value specified for GetByPhoneNumber")
	}

	var r0 []entity.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]entity.Transaction, error)); ok {
		return rf(ctx, phoneNumber)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []entity.Transaction); ok {
		r0 = rf(ctx, phoneNumber)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, phoneNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_GetByPhoneNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByPhoneNumber'
type Repository_GetByPhoneNumber_Call struct {
	*mock.Call
}

// GetByPhoneNumber is a helper method to define mock.On call
//   - ctx context.Context
//   - phoneNumber string
func (_e *Repository_Expecter) GetByPhoneNumber(ctx interface{}, phoneNumber interface{}) *Repository_GetByPhoneNumber_Call {
	return &Repository_GetByPhoneNumber_Call{Call: _e.mock.On("GetByPhoneNumber", ctx, phoneNumber)}
}

func (_c *Repository_GetByPhoneNumber_Call) Run(run func(ctx context.Context, phoneNumber string)) *Repository_GetByPhoneNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Repository_GetByPhoneNumber_Call) Return(res []entity.Transaction, err error) *Repository_GetByPhoneNumber_Call {
	_c.Call.Return(res, err)
	return _c
}

func (_c *Repository_GetByPhoneNumber_Call) RunAndReturn(run func(context.Context, string) ([]entity.Transaction, error)) *Repository_GetByPhoneNumber_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateStatus provides a mock function with given fields: ctx, id, status
func (_m *Repository) UpdateStatus(ctx context.Context, id string, status entity.Status) error {
	ret := _m.Called(ctx, id, status)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStatus")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, entity.Status) error); ok {
		r0 = rf(ctx, id, status)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Repository_UpdateStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatus'
type Repository_UpdateStatus_Call struct {
	*mock.Call
}

// UpdateStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
//   - status entity.Status
func (_e *Repository_Expecter) UpdateStatus(ctx interface{}, id interface{}, status interface{}) *Repository_UpdateStatus_Call {
	return &Repository_UpdateStatus_Call{Call: _e.mock.On("UpdateStatus", ctx, id, status)}
}

func (_c *Repository_UpdateStatus_Call) Run(run func(ctx context.Context, id string, status entity.Status)) *Repository_UpdateStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(entity.Status))
	})
	return _c
}

func (_c *Repository_UpdateStatus_Call) Return(err error) *Repository_UpdateStatus_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *Repository_UpdateStatus_Call) RunAndReturn(run func(context.Context, string, entity.Status) error) *Repository_UpdateStatus_Call {
	_c.Call.Return(run)
	return _c
}

// WithTx provides a mock function with given fields: tx
func (_m *Repository) WithTx(tx pgx.Tx) *transactionrepo.Repository {
	ret := _m.Called(tx)

	if len(ret) == 0 {
		panic("no return value specified for WithTx")
	}

	var r0 *transactionrepo.Repository
	if rf, ok := ret.Get(0).(func(pgx.Tx) *transactionrepo.Repository); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transactionrepo.Repository)
		}
	}

	return r0
}

// Repository_WithTx_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithTx'
type Repository_WithTx_Call struct {
	*mock.Call
}

// WithTx is a helper method to define mock.On call
//   - tx pgx.Tx
func (_e *Repository_Expecter) WithTx(tx interface{}) *Repository_WithTx_Call {
	return &Repository_WithTx_Call{Call: _e.mock.On("WithTx", tx)}
}

func (_c *Repository_WithTx_Call) Run(run func(tx pgx.Tx)) *Repository_WithTx_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(pgx.Tx))
	})
	return _c
}

func (_c *Repository_WithTx_Call) Return(_a0 *transactionrepo.Repository) *Repository_WithTx_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Repository_WithTx_Call) RunAndReturn(run func(pgx.Tx) *transactionrepo.Repository) *Repository_WithTx_Call {
	_c.Call.Return(run)
	return _c
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}