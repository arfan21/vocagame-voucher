package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Product struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Price     decimal.Decimal `json:"price"`
	Stock     int             `json:"stock"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func (Product) TableName() string {
	return "products"
}
