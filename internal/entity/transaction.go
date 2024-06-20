package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Status string

const (
	StatusWaitingPayment Status = "waiting_payment"
	StatusCompleted      Status = "completed"
	StatusFailed         Status = "failed"
)

type Transaction struct {
	ID              string          `json:"id"`
	ProductID       string          `json:"product_id"`
	PaymentMethodID int             `json:"payment_method_id"`
	PhoneNumber     string          `json:"phone_number"`
	Quantity        int             `json:"quantity"`
	TotalPrice      decimal.Decimal `json:"total_price"`
	Status          Status          `json:"status"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	Product         Product         `json:"product"`
	PaymentMethod   PaymentMethod   `json:"payment_method"`
}

func (Transaction) TableName() string {
	return "transactions"
}
