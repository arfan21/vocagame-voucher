package model

import "github.com/shopspring/decimal"

type TransactionRequest struct {
	ProductID       string `json:"product_id" validate:"required"`
	PaymentMethodID int    `json:"payment_method_id" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Quantity        int    `json:"quantity" validate:"required,min=1"`
}

type TransactionCreatedResponse struct {
	ID          string `json:"id"`
	RedirectURL string `json:"redirect_url"`
}

type TransactionTrackingRequest struct {
	Email string `json:"email" query:"email" validate:"required,email"`
}

type TransactionResponse struct {
	ID            string          `json:"id"`
	ProductName   string          `json:"product_name"`
	PaymentMethod string          `json:"payment_method"`
	Email         string          `json:"email"`
	Quantity      int             `json:"quantity"`
	TotalPrice    decimal.Decimal `json:"total_price"`
	Status        string          `json:"status"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
}
