package entity_test

import (
	"testing"
	"time"

	"go-currency-exchange/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	t.Run("Should return error when description is empty", func(t *testing.T) {
		_, err := entity.NewTransaction("", 10, "")
		assert.Equal(t, entity.ErrEmptyDescription, err)
	})

	t.Run("Should return error when description is longer than 50 characters", func(t *testing.T) {
		_, err := entity.NewTransaction("This is a description that is longer than 50 characters", 10, "")
		assert.Equal(t, entity.ErrDescriptionExceedsMaxLength, err)
	})

	t.Run("Should return error when value is less than or equal to 0", func(t *testing.T) {
		_, err := entity.NewTransaction("Description", 0, "")
		assert.Equal(t, entity.ErrInvalidValue, err)
	})

	t.Run("Should return a new transaction with current time when createdAt is empty", func(t *testing.T) {
		transaction, err := entity.NewTransaction("Description", 10, "")
		assert.Nil(t, err)
		assert.Equal(t, "Description", transaction.Description)
		assert.Equal(t, 10.0, transaction.Value)
		assert.NotEmpty(t, transaction.ID)
		assert.NotEmpty(t, transaction.CreatedAt)
		assert.NotEmpty(t, transaction.UpdatedAt)
	})

	t.Run("Should return a new transaction with the provided createdAt", func(t *testing.T) {
		createdAt := time.Now().Add(-time.Hour).Format(time.RFC3339Nano)
		transaction, err := entity.NewTransaction("Description", 10, createdAt)
		assert.Nil(t, err)
		assert.Equal(t, "Description", transaction.Description)
		assert.Equal(t, 10.0, transaction.Value)
		assert.NotEmpty(t, transaction.ID)
		assert.NotEmpty(t, transaction.CreatedAt)
		assert.NotEmpty(t, transaction.UpdatedAt)
		assert.Equal(t, createdAt, transaction.CreatedAt.Format(time.RFC3339Nano))
		assert.Equal(t, createdAt, transaction.UpdatedAt.Format(time.RFC3339Nano))
	})
}
