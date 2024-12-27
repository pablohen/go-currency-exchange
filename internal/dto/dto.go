package dto

type CreateTransactionInput struct {
	Description string  `json:"description"`
	Value       float64 `json:"value"`
}
