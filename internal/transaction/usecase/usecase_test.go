package transactionuc

import (
	"context"
	"fmt"
	"testing"

	"github.com/arfan21/vocagame/internal/model"
	productrepo "github.com/arfan21/vocagame/internal/product/repository"
	productusecase "github.com/arfan21/vocagame/internal/product/usecase"
	transactionrepo "github.com/arfan21/vocagame/internal/transaction/repository"
	mockPayment "github.com/arfan21/vocagame/mocks/client/payment"
	mockUc "github.com/arfan21/vocagame/mocks/internal_/transaction/usecase"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	uc              *UseCase
	repo            Repository
	productUc       *productusecase.UseCase
	paymentMethodUc *mockUc.PaymentMethodUsecase
	midtransCoreAPI *mockUc.MidtransCoreAPI
	clientPayment   *mockPayment.Payment
)

func initPgMock(t *testing.T) pgxmock.PgxPoolIface {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	return mock
}

func setup(t *testing.T, dbMock pgxmock.PgxPoolIface) {
	repo = transactionrepo.New(dbMock)

	productRepo := productrepo.New(dbMock)
	productUc = productusecase.New(productRepo)
	paymentMethodUc = mockUc.NewPaymentMethodUsecase(t)
	midtransCoreAPI = mockUc.NewMidtransCoreAPI(t)

	clientPayment = mockPayment.NewPayment(t)

	uc = New(repo, productUc, paymentMethodUc, midtransCoreAPI)
}

func TestCreate(t *testing.T) {
	dbMock := initPgMock(t)
	setup(t, dbMock)

	paramSucces := model.TransactionRequest{
		ProductID:       "1",
		Quantity:        1,
		PaymentMethodID: 1,
		PhoneNumber:     "+628123456789",
	}

	test := []struct {
		name       string
		init       func()
		excpectErr bool
		param      model.TransactionRequest
	}{
		{
			name: "success",
			init: func() {
				paymentMethodUc.On("GetByID", mock.Anything, paramSucces.PaymentMethodID).Return(model.PaymentMethodResponse{}, nil).Once()

				clientPayment.EXPECT().Pay(mock.Anything, mock.Anything, decimal.Zero).Return("1", nil).Once()

				// handle pgxmock
				dbMock.ExpectQuery("SELECT (.+) FROM products WHERE id = (.+) FOR UPDATE").
					WithArgs(paramSucces.ProductID).
					WillReturnRows(dbMock.NewRows([]string{"id", "name", "price", "stock", "created_at", "updated_at"}).
						AddRow(paramSucces.ProductID, "product-1", decimal.Zero, 10, "2021-01-01", "2021-01-01"))

				dbMock.ExpectBegin()
				dbMock.ExpectExec("INSERT INTO transactions (.+) VALUES (.+)").
					WithArgs(pgxmock.AnyArg(), paramSucces.ProductID, paramSucces.PaymentMethodID, paramSucces.PhoneNumber, paramSucces.Quantity, pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				dbMock.ExpectExec("UPDATE products SET stock = (.+) WHERE id = (.+)").
					WithArgs(paramSucces.Quantity, paramSucces.ProductID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				dbMock.ExpectCommit()

			},
			excpectErr: false,
			param:      paramSucces,
		},
		{
			name:       "failed, validation error",
			init:       func() {},
			excpectErr: true,
			param: model.TransactionRequest{
				ProductID: "",
			},
		},
		{
			name: "failed, get product by id",
			init: func() {

				// handle pgxmock
				dbMock.ExpectQuery("SELECT (.+) FROM products WHERE id = (.+) FOR UPDATE").
					WithArgs(paramSucces.ProductID).
					WillReturnError(fmt.Errorf("unexpected error"))

			},
			excpectErr: true,
			param:      paramSucces,
		},
		{
			name: "failed, stock is not enough",
			init: func() {

				// handle pgxmock
				dbMock.ExpectQuery("SELECT (.+) FROM products WHERE id = (.+) FOR UPDATE").
					WithArgs(paramSucces.ProductID).
					WillReturnRows(dbMock.NewRows([]string{"id", "name", "price", "stock", "created_at", "updated_at"}).
						AddRow(paramSucces.ProductID, "product-1", decimal.Zero, 0, "2021-01-01", "2021-01-01"))
			},
			excpectErr: true,
			param:      paramSucces,
		},
		{
			name: "failed, get payment method by id",
			init: func() {
				paymentMethodUc.On("GetByID", mock.Anything, paramSucces.PaymentMethodID).Return(model.PaymentMethodResponse{}, fmt.Errorf("unexpected error")).Once()

				dbMock.ExpectQuery("SELECT (.+) FROM products WHERE id = (.+) FOR UPDATE").
					WithArgs(paramSucces.ProductID).
					WillReturnRows(dbMock.NewRows([]string{"id", "name", "price", "stock", "created_at", "updated_at"}).
						AddRow(paramSucces.ProductID, "product-1", decimal.Zero, 10, "2021-01-01", "2021-01-01"))
			},
			excpectErr: true,
			param:      paramSucces,
		},
		{
			name: "failed, create transaction",
			init: func() {
				paymentMethodUc.On("GetByID", mock.Anything, paramSucces.PaymentMethodID).Return(model.PaymentMethodResponse{}, nil).Once()

				// handle pgxmock
				dbMock.ExpectQuery("SELECT (.+) FROM products WHERE id = (.+) FOR UPDATE").
					WithArgs(paramSucces.ProductID).
					WillReturnRows(dbMock.NewRows([]string{"id", "name", "price", "stock", "created_at", "updated_at"}).
						AddRow(paramSucces.ProductID, "product-1", decimal.Zero, 10, "2021-01-01", "2021-01-01"))

				dbMock.ExpectBegin()
				dbMock.ExpectExec("INSERT INTO transactions (.+) VALUES (.+)").
					WithArgs(pgxmock.AnyArg(), paramSucces.ProductID, paramSucces.PaymentMethodID, paramSucces.PhoneNumber, paramSucces.Quantity, pgxmock.AnyArg()).
					WillReturnError(fmt.Errorf("unexpected error"))

				dbMock.ExpectRollback()
			},
			excpectErr: true,
			param:      paramSucces,
		},
		{
			name: "failed, update stock",
			init: func() {
				paymentMethodUc.On("GetByID", mock.Anything, paramSucces.PaymentMethodID).Return(model.PaymentMethodResponse{}, nil).Once()

				// handle pgxmock
				dbMock.ExpectQuery("SELECT (.+) FROM products WHERE id = (.+) FOR UPDATE").
					WithArgs(paramSucces.ProductID).
					WillReturnRows(dbMock.NewRows([]string{"id", "name", "price", "stock", "created_at", "updated_at"}).
						AddRow(paramSucces.ProductID, "product-1", decimal.Zero, 10, "2021-01-01", "2021-01-01"))

				dbMock.ExpectBegin()
				dbMock.ExpectExec("INSERT INTO transactions (.+) VALUES (.+)").
					WithArgs(pgxmock.AnyArg(), paramSucces.ProductID, paramSucces.PaymentMethodID, paramSucces.PhoneNumber, paramSucces.Quantity, pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				dbMock.ExpectExec("UPDATE products SET stock = (.+) WHERE id = (.+)").
					WithArgs(paramSucces.Quantity, paramSucces.ProductID).
					WillReturnError(fmt.Errorf("unexpected error"))

				dbMock.ExpectRollback()
			},
			excpectErr: true,
			param:      paramSucces,
		},
		{
			name: "failed, create payment",
			init: func() {
				paymentMethodUc.On("GetByID", mock.Anything, paramSucces.PaymentMethodID).Return(model.PaymentMethodResponse{}, nil).Once()

				clientPayment.EXPECT().Pay(mock.Anything, mock.Anything, decimal.Zero).Return("", fmt.Errorf("unexpected error")).Once()

				// handle pgxmock
				dbMock.ExpectQuery("SELECT (.+) FROM products WHERE id = (.+) FOR UPDATE").
					WithArgs(paramSucces.ProductID).
					WillReturnRows(dbMock.NewRows([]string{"id", "name", "price", "stock", "created_at", "updated_at"}).
						AddRow(paramSucces.ProductID, "product-1", decimal.Zero, 10, "2021-01-01", "2021-01-01"))

				dbMock.ExpectBegin()
				dbMock.ExpectExec("INSERT INTO transactions (.+) VALUES (.+)").
					WithArgs(pgxmock.AnyArg(), paramSucces.ProductID, paramSucces.PaymentMethodID, paramSucces.PhoneNumber, paramSucces.Quantity, pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				dbMock.ExpectExec("UPDATE products SET stock = (.+) WHERE id = (.+)").
					WithArgs(paramSucces.Quantity, paramSucces.ProductID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				dbMock.ExpectRollback()
			},
			excpectErr: true,
			param:      paramSucces,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			tt.init()
			res, err := uc.Create(context.Background(), tt.param, clientPayment)
			fmt.Println("err", err)
			if tt.excpectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, res)
			assert.NotEmpty(t, res.ID)

			paymentMethodUc.AssertExpectations(t)
			clientPayment.AssertExpectations(t)
		})
	}

}
