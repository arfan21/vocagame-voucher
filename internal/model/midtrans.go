package model

import (
	"crypto/sha512"
	"encoding/hex"

	"github.com/arfan21/vocagame/config"
)

// {
// 	"transaction_time": "2021-06-15 18:45:13",
// 	"transaction_status": "settlement",
// 	"transaction_id": "513f1f01-c9da-474c-9fc9-d5c64364b709",
// 	"status_message": "midtrans payment notification",
// 	"status_code": "200",
// 	"signature_key": "2496c78cac93a70ca08014bdaaff08eb7119ef79ef69c4833d4399cada077147febc1a231992eb8665a7e26d89b1dc323c13f721d21c7485f70bff06cca6eed3",
// 	"settlement_time": "2021-06-15 18:45:28",
// 	"payment_type": "gopay",
// 	"order_id": "Order-5100",
// 	"merchant_id": "G141532850",
// 	"gross_amount": "154600.00",
// 	"fraud_status": "accept",
// 	"currency": "IDR"
//   }

type MidtransNotification struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	SettlementTime    string `json:"settlement_time"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	Currency          string `json:"currency"`
}

// ValidateSignatureKey is a method to validate signature key
// SHA512(order_id+status_code+gross_amount+ServerKey)
func (m MidtransNotification) ValidateSignatureKey() bool {
	serverKey := config.Get().Midtrans.ServerKey

	signatureKey := m.OrderID + m.StatusCode + m.GrossAmount + serverKey
	hash := sha512.New()
	hash.Write([]byte(signatureKey))

	generatedKey := hex.EncodeToString(hash.Sum(nil))
	return m.SignatureKey == generatedKey
}
