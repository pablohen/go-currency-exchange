package dto

type CreateTransactionInput struct {
	Description string  `json:"description"`
	Value       float64 `json:"value"`
}

type TransactionOutput struct {
	ID             string  `json:"id"`
	Description    string  `json:"description"`
	CreatedAt      string  `json:"created_at"`
	ConversionRate float64 `json:"conversion_rate"`
	OriginalValue  float64 `json:"original_value"`
	ConvertedValue float64 `json:"converted_value"`
}
