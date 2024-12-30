package dto

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
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Value       float64 `json:"value"`
	CreatedAt   string  `json:"created_at"`
}
