package entity

import "time"

type CreateTransactionRequest struct {
	Products []TransactionProductRequest `json:"products" validate:"required,dive"`
}

type TransactionProductRequest struct {
	ProductID string `json:"productId" validate:"required,uuid"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

type CreateTransactionResponse struct {
	TransactionID string                       `json:"transaction_id"`
	TotalPrice    float64                      `json:"total_price"`
	Products      []TransactionProductResponse `json:"products"`
}

type TransactionProductResponse struct {
	ProductID   string  `json:"productId"`
	ProductName string  `json:"productName"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Subtotal    float64 `json:"subtotal"`
}

// Struktur response untuk detail transaksi
type TransactionResponse struct {
	ID          string                        `json:"id"`
	UserID      string                        `json:"userId"`
	TotalAmount float64                       `json:"totalAmount"`
	CreatedAt   time.Time                     `json:"createdAt"`
	Details     *[]TransactionProductResponse `json:"details,omitempty"`
}
