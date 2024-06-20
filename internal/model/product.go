package model

import "github.com/shopspring/decimal"

type ProductResponse struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Price     decimal.Decimal `json:"price"`
	Stock     int             `json:"stock"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}
