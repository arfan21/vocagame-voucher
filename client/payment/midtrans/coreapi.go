package paymentmidtrans

import (
	"context"

	"github.com/arfan21/vocagame/config"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type MidtransCoreAPI struct {
	client coreapi.Client
}

func NewMidtransCoreAPI() *MidtransCoreAPI {
	s := coreapi.Client{}
	s.New(config.Get().Midtrans.ServerKey, midtrans.Sandbox)
	return &MidtransCoreAPI{
		client: s,
	}
}

func (m MidtransCoreAPI) CheckTransaction(ctx context.Context, orderID string) (coreapi.TransactionStatusResponse, error) {
	res, errCore := m.client.CheckTransaction(orderID)
	if errCore != nil {
		return coreapi.TransactionStatusResponse{}, errCore
	}

	if res != nil {
		return *res, nil
	}

	return coreapi.TransactionStatusResponse{}, nil
}
