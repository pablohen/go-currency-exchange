package entity

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID          string  `json:"id" gorm:"primaryKey"`
	Description string  `json:"description"`
	Value       float64 `json:"value"`
	gorm.Model  `json:"-"`
}

var (
	DescriptionMaxLength           = 50
	ErrDescriptionExceedsMaxLength = errors.New(fmt.Sprintf("description can't be longer than %d characters", DescriptionMaxLength))
	ErrEmptyDescription            = errors.New("description can't be empty")
	ErrInvalidValue                = errors.New("value must be greater than 0")
)

func NewTransaction(description string, value float64, createdAt string) (*Transaction, error) {
	if description == "" {
		return nil, ErrEmptyDescription
	}

	if len(description) > DescriptionMaxLength {
		return nil, ErrDescriptionExceedsMaxLength
	}

	if value <= 0 {
		return nil, ErrInvalidValue
	}

	parsedCreatedAt := time.Now().UTC()

	if createdAt != "" {
		parse, err := time.Parse(time.RFC3339Nano, createdAt)
		if err != nil {
			return nil, err
		}

		parsedCreatedAt = parse
	}

	return &Transaction{
		ID:          uuid.New().String(),
		Description: description,
		Value:       value,
		Model: gorm.Model{
			CreatedAt: parsedCreatedAt,
			UpdatedAt: parsedCreatedAt,
		},
	}, nil
}
