package worker_test

import (
	"encoding/json"
	"testing"
	"time"

	"go-currency-exchange/internal/dto"
	"go-currency-exchange/internal/entity"
	"go-currency-exchange/internal/worker"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(description string, value float64, createdAt string) error {
	args := m.Called(description, value, createdAt)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetById(id string) (*entity.Transaction, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetAll() ([]entity.Transaction, error) {
	args := m.Called()
	return args.Get(0).([]entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetAllPaginated(page int, pageSize int) (dto.TransactionsPaginated, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).(dto.TransactionsPaginated), args.Error(1)
}

func TestCreateTransaction(t *testing.T) {
	messageChan := make(chan amqp.Delivery)
	mockRepo := new(MockTransactionRepository)

	go worker.CreateTransaction(messageChan, mockRepo)

	transactionMessage := dto.TransactionMessage{
		Description: "Test transaction",
		Value:       100.0,
		CreatedAt:   time.Now().Format(time.RFC3339Nano),
	}
	messageBody, _ := json.Marshal(transactionMessage)

	mockRepo.On("Create", transactionMessage.Description, transactionMessage.Value, transactionMessage.CreatedAt).Return(nil)

	messageChan <- amqp.Delivery{Body: messageBody}

	mockRepo.AssertExpectations(t)
}

func TestCreateTransaction_InvalidMessage(t *testing.T) {
	messageChan := make(chan amqp.Delivery)
	mockRepo := new(MockTransactionRepository)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				assert.NotNil(t, r)
			}
		}()
		worker.CreateTransaction(messageChan, mockRepo)
	}()

	invalidMessageBody := []byte("invalid message")

	messageChan <- amqp.Delivery{Body: invalidMessageBody}
}

func TestCreateTransaction_RepositoryError(t *testing.T) {
	messageChan := make(chan amqp.Delivery)
	mockRepo := new(MockTransactionRepository)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				assert.NotNil(t, r)
			}
		}()
		worker.CreateTransaction(messageChan, mockRepo)
	}()

	transactionMessage := dto.TransactionMessage{
		Description: "Test transaction",
		Value:       100.0,
		CreatedAt:   time.Now().Format(time.RFC3339Nano),
	}
	messageBody, _ := json.Marshal(transactionMessage)

	mockRepo.On("Create", transactionMessage.Description, transactionMessage.Value, transactionMessage.CreatedAt).Return(assert.AnError)

	messageChan <- amqp.Delivery{Body: messageBody}

	mockRepo.AssertExpectations(t)
}

func TestCreateTransaction_ValidMessage(t *testing.T) {
	messageChan := make(chan amqp.Delivery)
	mockRepo := new(MockTransactionRepository)

	go worker.CreateTransaction(messageChan, mockRepo)

	transactionMessage := dto.TransactionMessage{
		Description: "Valid transaction",
		Value:       200.0,
		CreatedAt:   time.Now().Format(time.RFC3339Nano),
	}
	messageBody, _ := json.Marshal(transactionMessage)

	mockRepo.On("Create", transactionMessage.Description, transactionMessage.Value, transactionMessage.CreatedAt).Return(nil)

	messageChan <- amqp.Delivery{Body: messageBody}

	mockRepo.AssertExpectations(t)
}
