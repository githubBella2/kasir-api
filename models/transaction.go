package models

import "time"

type Transaction struct {
	ID                 int                `json:"id"`
	TotalAmount        int                `json:"total_amount"`
	Details []TransactionDetails `json:"details"`
	CreatedAt   time.Time            `json:"created_at,omitempty"`
}

type TransactionDetails struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name"`
	Quantity      int    `json:"quantity"`
	Subtotal      int    `json:"subtotal"`
}

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
