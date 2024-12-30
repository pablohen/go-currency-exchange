package database

import (
	"go-currency-exchange/internal/entity"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		DB: db,
	}
}

func (t *TransactionRepository) Create(description string, value float64, createdAt string) error {
	transaction, err := entity.NewTransaction(description, value, createdAt)
	if err != nil {
		return err
	}

	err = t.DB.Create(&transaction).Error
	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionRepository) GetById(id string) (*entity.Transaction, error) {
	var transaction entity.Transaction
	err := t.DB.First(&transaction, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (t *TransactionRepository) GetAll() ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := t.DB.Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (t *TransactionRepository) GetAllPaginated(page int, pageSize int) (ItemsPaginated, error) {
	var transactions []entity.Transaction
	err := t.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&transactions).Error
	if err != nil {
		return ItemsPaginated{}, err
	}

	var total int64
	err = t.DB.Model(&entity.Transaction{}).Count(&total).Error
	if err != nil {
		return ItemsPaginated{}, err
	}

	return ItemsPaginated{
		Items:    transactions,
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	}, nil
}
