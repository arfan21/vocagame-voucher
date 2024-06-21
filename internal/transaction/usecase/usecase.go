package transactionuc

import (
	"context"
	"fmt"
	"time"

	clientpayment "github.com/arfan21/vocagame/client/payment"
	"github.com/arfan21/vocagame/internal/entity"
	"github.com/arfan21/vocagame/internal/model"
	productusecase "github.com/arfan21/vocagame/internal/product/usecase"
	transactionrepo "github.com/arfan21/vocagame/internal/transaction/repository"
	"github.com/arfan21/vocagame/pkg/constant"
	"github.com/arfan21/vocagame/pkg/logger"
	"github.com/arfan21/vocagame/pkg/validation"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/shopspring/decimal"
)

type Repository interface {
	Begin(ctx context.Context) (tx pgx.Tx, err error)
	WithTx(tx pgx.Tx) *transactionrepo.Repository
	Create(ctx context.Context, data entity.Transaction) (err error)
	UpdateStatus(ctx context.Context, id string, status entity.Status) (err error)
	GetByEmail(ctx context.Context, email string) (res []entity.Transaction, err error)
	GetByID(ctx context.Context, id string) (res entity.Transaction, err error)
}

type ProductUsecase interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	WithTx(tx pgx.Tx) *productusecase.UseCase
	GetByID(ctx context.Context, id string, isForUpdate bool) (model.ProductResponse, error)
	ReduceStock(ctx context.Context, id string, qty int) error
}

type PaymentMethodUsecase interface {
	GetByID(ctx context.Context, id int) (model.PaymentMethodResponse, error)
}

type MidtransCoreAPI interface {
	CheckTransaction(ctx context.Context, orderID string) (coreapi.TransactionStatusResponse, error)
}

type NotificationProducer interface {
	Produce(ctx context.Context, event model.Event) error
}

type UseCase struct {
	repo            Repository
	productUc       ProductUsecase
	paymentMethodUc PaymentMethodUsecase
	midtransCoreAPI MidtransCoreAPI
	notifProducer   NotificationProducer
}

func New(
	repo Repository,
	productUc ProductUsecase,
	paymentMethodUc PaymentMethodUsecase,
	midtransCoreAPI MidtransCoreAPI,
	notifProducer NotificationProducer,
) *UseCase {
	return &UseCase{repo: repo, productUc: productUc, paymentMethodUc: paymentMethodUc, midtransCoreAPI: midtransCoreAPI, notifProducer: notifProducer}
}

func (uc UseCase) Create(ctx context.Context, data model.TransactionRequest, paymentClient clientpayment.Payment) (res model.TransactionCreatedResponse, err error) {
	err = validation.Validate(data)
	if err != nil {
		err = fmt.Errorf("transaction.uc.Create: failed to validate data: %w", err)
		return
	}

	product, err := uc.productUc.GetByID(ctx, data.ProductID, true)
	if err != nil {
		err = fmt.Errorf("transaction.uc.Create: failed to get product: %w", err)
		return
	}

	if product.Stock < data.Quantity {
		err = fmt.Errorf("transaction.uc.Create: failed to create transaction: stock is not enough, %w", constant.ErrOutOfStock)
		return
	}

	_, err = uc.paymentMethodUc.GetByID(ctx, data.PaymentMethodID)
	if err != nil {
		err = fmt.Errorf("transaction.uc.Create: failed to get payment method: %w", err)
		return
	}

	totalPrice := product.Price.Mul(decimal.NewFromInt(int64(data.Quantity)))

	tx, err := uc.repo.Begin(ctx)
	if err != nil {
		err = fmt.Errorf("transaction.uc.Create: failed to begin transaction: %w", err)
		return
	}

	defer func() {
		if err != nil {
			logger.Log(ctx).Error().Err(err).Msg("transaction.uc.Create: failed to create transaction, rollback transaction")
			errRb := tx.Rollback(ctx)
			if errRb != nil {
				err = fmt.Errorf("transaction.uc.Create: failed to rollback transaction: %w", errRb)
			}
			return
		}

		err = tx.Commit(ctx)
	}()

	id, err := uuid.NewV7()
	if err != nil {
		err = fmt.Errorf("transaction.uc.Create: failed to generate transaction id: %w", err)
		return
	}

	err = uc.repo.WithTx(tx).Create(ctx, entity.Transaction{
		ID:              id.String(),
		ProductID:       data.ProductID,
		PaymentMethodID: data.PaymentMethodID,
		Email:           data.Email,
		Quantity:        data.Quantity,
		TotalPrice:      totalPrice,
	})
	if err != nil {
		err = fmt.Errorf("transaction.uc.Create: failed to create transaction: %w", err)
		return
	}

	err = uc.productUc.WithTx(tx).ReduceStock(ctx, data.ProductID, data.Quantity)
	if err != nil {
		err = fmt.Errorf("transaction.uc.Create: failed to reduce stock: %w", err)
		return
	}

	resPayment, err := paymentClient.Pay(ctx, id.String(), totalPrice)
	if err != nil {
		err = fmt.Errorf("transaction.uc.Create: failed to pay: %w", err)
		return
	}

	res = model.TransactionCreatedResponse{
		ID:          id.String(),
		RedirectURL: resPayment,
	}

	err = uc.notifProducer.Produce(ctx, model.Event{
		Name:          model.EventTransactionNotification,
		TransactionID: id.String(),
		Email:         data.Email,
		Url:           resPayment,
	})
	if err != nil {
		err = fmt.Errorf("transaction.uc.Create: failed to produce notification: %w", err)
		return
	}

	return res, nil
}

func (uc UseCase) MidtransNotification(ctx context.Context, data model.MidtransNotification) (err error) {
	if !data.ValidateSignatureKey() {
		err = fmt.Errorf("transaction.uc.MidtransNotification: failed to validate signature key: %w", constant.ErrInvalidSignatureKey)
		return
	}

	transactionStatusResp, e := uc.midtransCoreAPI.CheckTransaction(ctx, data.OrderID)
	if e != nil {
		err = fmt.Errorf("transaction.uc.MidtransNotification: failed to check transaction: %w", e)
		return
	}

	statusForUpdate := entity.StatusWaitingPayment

	switch transactionStatusResp.TransactionStatus {
	case "capture":
		if transactionStatusResp.FraudStatus == "challenge" {
			// update status on database to 'challenge'
			// e.g: 'Payment status challenged. Please take action on Merchant Administration Portal
		} else if transactionStatusResp.FraudStatus == "accept" {
			statusForUpdate = entity.StatusCompleted
		}
	case "settlement":
		// update status on database to 'success'
		statusForUpdate = entity.StatusCompleted
	case "deny":
		// TODO you can ignore 'deny', because most of the time it allows payment retries
		// and later can become success
	case "cancel", "expire":
		// update status on database to 'failure'
		statusForUpdate = entity.StatusFailed
	case "pending":
		// update status on database to 'pending' / waiting payment
	}

	transactionDataDB, err := uc.repo.GetByID(ctx, data.OrderID)
	if err != nil {
		err = fmt.Errorf("transaction.uc.MidtransNotification: failed to get transaction data: %w", err)
		return
	}

	err = uc.repo.UpdateStatus(ctx, data.OrderID, statusForUpdate)
	if err != nil {
		err = fmt.Errorf("transaction.uc.MidtransNotification: failed to update status: %w", err)
		return
	}

	errNotif := uc.notifProducer.Produce(ctx, model.Event{
		Name:          model.EventTransactionNotification,
		TransactionID: transactionDataDB.ID,
		Email:         transactionDataDB.Email,
	})
	// if failed to produce notification, just log it
	if errNotif != nil {
		logger.Log(ctx).Warn().Err(errNotif).Msg("transaction.uc.MidtransNotification: failed to produce notification")
		return
	}

	return nil
}

func (uc UseCase) GetByEmail(ctx context.Context, req model.TransactionTrackingRequest) (res []model.TransactionResponse, err error) {
	// if phone number not start with '+', add '+'
	// maybe + is missing because url query param will be encoded
	// if len(req.PhoneNumber) > 0 && req.PhoneNumber[0] != '+' {
	// 	// if phone number start with empty space, remove it
	// 	if req.PhoneNumber[0] == 32 {
	// 		req.PhoneNumber = req.PhoneNumber[1:]
	// 	}
	// 	req.PhoneNumber = "+" + req.PhoneNumber
	// }

	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("transaction.uc.GetByPhoneNumber: failed to validate request: %w", err)
		return
	}

	transactions, err := uc.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		err = fmt.Errorf("transaction.uc.GetByPhoneNumber: failed to get transaction by phone number: %w", err)
		return
	}

	res = make([]model.TransactionResponse, len(transactions))

	for i, transaction := range transactions {
		res[i] = model.TransactionResponse{
			ID:            transaction.ID,
			ProductName:   transaction.Product.Name,
			PaymentMethod: transaction.PaymentMethod.Name,
			Email:         transaction.Email,
			Quantity:      transaction.Quantity,
			TotalPrice:    transaction.TotalPrice,
			Status:        string(transaction.Status),
			CreatedAt:     transaction.CreatedAt.Format(time.DateTime),
			UpdatedAt:     transaction.UpdatedAt.Format(time.DateTime),
		}
	}

	return res, nil
}
