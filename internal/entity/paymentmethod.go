package entity

type PaymentMethod struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (PaymentMethod) TableName() string {
	return "payment_methods"
}
