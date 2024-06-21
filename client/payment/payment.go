package clientpayment

import (
	"context"

	paymentmidtrans "github.com/arfan21/vocagame/client/payment/midtrans"
	"github.com/shopspring/decimal"
)

type Payment interface {
	Pay(ctx context.Context, ID string, amount decimal.Decimal) (string, error)
}

func GetPaymentMethod(paymentMethodId int, email string) Payment {

	// Check if payment method is available
	if _, ok := paymentmidtrans.MidtransPaymentMethod[paymentMethodId]; ok {
		return paymentmidtrans.NewMidtransSnapCharge(email)
	}

	return nil
}
