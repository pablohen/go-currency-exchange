package database_test

import (
	"testing"
	"time"

	"go-currency-exchange/internal/entity"
	"go-currency-exchange/internal/infra/database"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&entity.Transaction{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestTransactionRepository_Create(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := database.NewTransactionRepository(db)

	t.Run("valid transaction", func(t *testing.T) {
		description := "Valid description"
		value := 100.0
		createdAt := time.Now().UTC().Format(time.RFC3339Nano)

		err := repo.Create(description, value, createdAt)
		assert.NoError(t, err)
	})

	t.Run("invalid transaction", func(t *testing.T) {
		description := ""
		value := 100.0
		createdAt := time.Now().UTC().Format(time.RFC3339Nano)

		err := repo.Create(description, value, createdAt)
		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyDescription, err)
	})
}

func TestTransactionRepository_GetById(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	description := "Valid description"
	value := 100.0
	createdAt := time.Now().UTC().Format(time.RFC3339Nano)

	transaction, err := entity.NewTransaction(description, value, createdAt)
	assert.NoError(t, err)

	err = db.Create(transaction).Error
	assert.NoError(t, err)

	repo := database.NewTransactionRepository(db)
	err = repo.Create(description, value, createdAt)
	assert.NoError(t, err)

	// Retrieve the newly created transaction
	newTransaction, err := repo.GetById(transaction.ID)
	assert.NoError(t, err)
	assert.NotNil(t, newTransaction)
	assert.Equal(t, description, newTransaction.Description)
	assert.Equal(t, value, newTransaction.Value)
	assert.NoError(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, description, transaction.Description)
	assert.Equal(t, value, transaction.Value)
}

func TestTransactionRepository_GetAll(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := database.NewTransactionRepository(db)

	description := "Valid description"
	value := 100.0
	createdAt := time.Now().UTC().Format(time.RFC3339Nano)

	err = repo.Create(description, value, createdAt)
	assert.NoError(t, err)

	transactions, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, transactions, 1)
	assert.Equal(t, description, transactions[0].Description)
	assert.Equal(t, value, transactions[0].Value)
}

func TestTransactionRepository_GetAllPaginated(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := database.NewTransactionRepository(db)

	for i := 0; i < 10; i++ {
		description := "Valid description"
		value := float64((i + 1) * 10)
		createdAt := time.Now().UTC().Format(time.RFC3339Nano)

		err = repo.Create(description, value, createdAt)
		assert.NoError(t, err)
	}

	page := 1
	pageSize := 5
	paginatedItems, err := repo.GetAllPaginated(page, pageSize)
	assert.NoError(t, err)
	assert.Len(t, paginatedItems.Items, pageSize)
	assert.Equal(t, page, paginatedItems.Page)
	assert.Equal(t, pageSize, paginatedItems.PageSize)
	assert.Equal(t, 10, paginatedItems.Total)
}
