package paymentmidtrans

import (
	"context"
	"fmt"

	"github.com/arfan21/vocagame/config"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/shopspring/decimal"
)

type MidtransSnapCharge struct {
	PhoneNumber string
}

func NewMidtransSnapCharge(phoneNumber string) *MidtransSnapCharge {
	return &MidtransSnapCharge{
		PhoneNumber: phoneNumber,
	}
}

func (m MidtransSnapCharge) Pay(ctx context.Context, ID string, amount decimal.Decimal) (string, error) {
	s := snap.Client{}
	s.New(config.Get().Midtrans.ServerKey, midtrans.Sandbox)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  ID,
			GrossAmt: amount.IntPart(),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			Phone: m.PhoneNumber,
		},
	}

	res, errMidtrans := s.CreateTransaction(req)
	if errMidtrans != nil {
		var err error
		err = errMidtrans

		err = fmt.Errorf("midtrans snap charge: failed to create transaction: %w", err)
		return "", err
	}

	if res != nil {
		return res.RedirectURL, nil
	}

	return "", nil
}
