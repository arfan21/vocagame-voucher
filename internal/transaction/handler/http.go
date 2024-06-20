package transactionhandler

import (
	"context"
	"net/http"

	clientpayment "github.com/arfan21/vocagame/client/payment"
	"github.com/arfan21/vocagame/internal/model"
	"github.com/arfan21/vocagame/pkg/constant"
	"github.com/arfan21/vocagame/pkg/pkgutil"
	"github.com/gofiber/fiber/v2"
)

type UseCase interface {
	Create(ctx context.Context, data model.TransactionRequest, paymentClient clientpayment.Payment) (res model.TransactionCreatedResponse, err error)
	MidtransNotification(ctx context.Context, data model.MidtransNotification) (err error)
	GetByPhoneNumber(ctx context.Context, req model.TransactionTrackingRequest) (res []model.TransactionResponse, err error)
}

type HTTP struct {
	uc UseCase
}

func NewHTTP(uc UseCase) *HTTP {
	return &HTTP{uc}
}

// @Summary Create Transaction
// @Description Create Transaction
// @Tags Transaction
// @Accept json
// @Produce json
// @Param data body model.TransactionRequest true "Transaction Data"
// @Success 200 {object} pkgutil.HTTPResponse{data=model.TransactionResponse}
// @Router /api/v1/transactions [post]
func (h HTTP) Create(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var data model.TransactionRequest
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	paymentMethod := clientpayment.GetPaymentMethod(data.PaymentMethodID, data.PhoneNumber)
	if paymentMethod == nil {
		return constant.ErrPaymentMethoUnavailable
	}

	res, err := h.uc.Create(ctx, data, paymentMethod)
	if err != nil {
		return err
	}

	return c.JSON(pkgutil.HTTPResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    res,
	})
}

// @Summary Midtrans Notification
// @Description Midtrans Notification
// @Tags Transaction
// @Accept json
// @Produce json
// @Param data body model.MidtransNotification true "Midtrans Notification Data"
// @Success 200 {object} pkgutil.HTTPResponse
// @Router /api/v1/transactions/callback/midtrans [post]
func (h HTTP) MidtransNotification(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var data model.MidtransNotification
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	err := h.uc.MidtransNotification(ctx, data)
	if err != nil {
		return err
	}

	return c.JSON(pkgutil.HTTPResponse{
		Code:    http.StatusOK,
		Message: "Success",
	})
}

// @Summary Transaction Success Callback
// @Description Transaction Success Callback
// @Tags Transaction
// @Accept json
// @Produce json
// @Success 200 {object} pkgutil.HTTPResponse
// @Router /api/v1/transactions/callback/success [get]
func (h HTTP) SuccessCallback(c *fiber.Ctx) error {
	return c.JSON(pkgutil.HTTPResponse{
		Code:    http.StatusOK,
		Message: "Success",
	})
}

// @Summary Transaction Tracking Status
// @Description Transaction Tracking Status
// @Tags Transaction
// @Accept json
// @Produce json
// @Param phone_number query string true "Phone Number"
// @Success 200 {object} pkgutil.HTTPResponse{data=model.TransactionResponse}
// @Router /api/v1/transactions/tracking [get]
func (h HTTP) Tracking(c *fiber.Ctx) error {
	ctx := c.UserContext()
	req := model.TransactionTrackingRequest{}
	if err := c.QueryParser(&req); err != nil {
		return err
	}

	res, err := h.uc.GetByPhoneNumber(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(pkgutil.HTTPResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    res,
	})
}
