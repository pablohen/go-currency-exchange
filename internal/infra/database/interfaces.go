package database

import "go-currency-exchange/internal/entity"

type TransactionInterface interface {
	Create(description string, value float64) error
	GetById(id string) (*entity.Transaction, error)
	GetAll() ([]entity.Transaction, error)
	GetAllPaginated(page int, pageSize int) (ItemsPaginated, error)
}

type ItemsPaginated struct {
	Items    interface{} `json:"items"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
	Total    int         `json:"total"`
}
