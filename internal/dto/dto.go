package dto

import "go-currency-exchange/internal/entity"

type CreateTransactionInput struct {
	Description string  `json:"description"`
	Value       float64 `json:"value"`
	CreatedAt   string  `json:"created_at"`
}

type TransactionOutput struct {
	ID             string  `json:"id"`
	Description    string  `json:"description"`
	CreatedAt      string  `json:"created_at"`
	ConversionRate float64 `json:"conversion_rate"`
	OriginalValue  float64 `json:"original_value"`
	ConvertedValue float64 `json:"converted_value"`
}

type TransactionMessage struct {
	Description string  `json:"description"`
	Value       float64 `json:"value"`
	CreatedAt   string  `json:"created_at"`
}

type ItemsPaginated[T any] struct {
	Items    []T `json:"items"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

type TransactionsPaginated struct {
	Items    []entity.Transaction `json:"items"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"pageSize"`
	Total    int                  `json:"total"`
}
