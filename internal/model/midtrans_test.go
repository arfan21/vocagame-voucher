package model

import (
	"testing"

	"github.com/arfan21/vocagame/config"
)

const (
	signature = "e78e2223638cb60dbdbc88d23deb9b927ac41be7263ab38758605bac834dc25425705543707504bfef0802914cfa3f5f538fa308d1f9086211c420e7892ba2ba"
)

func TestMidtransValidateSignature(t *testing.T) {
	orderId := "Postman-1578568851"
	statusCode := "200"
	grossAmount := "10000.00"
	serverKey := "VT-server-HJMpl9HLr_ntOKt5mRONdmKj"

	payload := MidtransNotification{
		SignatureKey: signature,
		OrderID:      orderId,
		StatusCode:   statusCode,
		GrossAmount:  grossAmount,
	}

	config.Get().Midtrans.ServerKey = serverKey

	if !payload.ValidateSignatureKey() {
		t.Error("Failed to validate signature key")
	}

	config.Get().Midtrans.ServerKey = ""
}
