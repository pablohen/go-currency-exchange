package database

import (
	"go-currency-exchange/internal/dto"
	"go-currency-exchange/internal/entity"
)

type TransactionInterface interface {
	Create(description string, value float64, createdAt string) error
	GetById(id string) (*entity.Transaction, error)
	GetAll() ([]entity.Transaction, error)
	GetAllPaginated(page int, pageSize int) (dto.TransactionsPaginated, error)
}
