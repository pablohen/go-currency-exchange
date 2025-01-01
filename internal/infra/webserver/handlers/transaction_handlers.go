package handlers

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"go-currency-exchange/internal/dto"
	"go-currency-exchange/internal/entity"
	"go-currency-exchange/internal/infra/database"
	"go-currency-exchange/pkg/rabbitmq"
	"go-currency-exchange/pkg/treasury"

	amqp "github.com/rabbitmq/amqp091-go"
)

type TransactionHandler struct {
	TransactionRepository database.TransactionInterface
	Channel               *amqp.Channel
}

func NewTransactionHandler(db database.TransactionInterface, channel *amqp.Channel) *TransactionHandler {
	return &TransactionHandler{
		TransactionRepository: db,
		Channel:               channel,
	}
}

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Create a new transaction with the given details
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body dto.CreateTransactionInput true "Transaction input"
// @Success 201 {object} dto.TransactionMessage
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var createTransactionInput dto.CreateTransactionInput
	err := json.NewDecoder(r.Body).Decode(&createTransactionInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}

	transaction, err := entity.NewTransaction(createTransactionInput.Description, createTransactionInput.Value, createTransactionInput.CreatedAt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	transactionMessage := dto.TransactionMessage{
		Description: transaction.Description,
		Value:       transaction.Value,
		CreatedAt:   transaction.CreatedAt.UTC().Format(time.RFC3339Nano),
	}
	transactionMessageJSON, err := json.Marshal(transactionMessage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = rabbitmq.Publish(h.Channel, "transactions", string(transactionMessageJSON))
	if err != nil {
		log.Printf("Failed to publish message, reconnecting to RabbitMQ: %v", err)
		h.Channel, _ = rabbitmq.Connect()
		err = rabbitmq.Publish(h.Channel, "transactions", string(transactionMessageJSON))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

// GetTransactionByIdWithExchangeRate godoc
// @Summary Get a transaction by ID with exchange rate
// @Description Get a transaction by ID and convert its value using the exchange rate at the time of creation
// @Tags transactions
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} dto.TransactionOutput
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Transaction not found"
// @Failure 500 {string} string "Internal server error"
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetTransactionByIdWithExchangeRate(w http.ResponseWriter, r *http.Request) {
	transactionID := r.PathValue("id")
	if transactionID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	transaction, err := h.TransactionRepository.GetById(transactionID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	rate, err := treasury.GetRatesByDate(transaction.CreatedAt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	floatExchangeRate, err := strconv.ParseFloat(rate.ExchangeRate, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	transactionOutput := dto.TransactionOutput{
		ID:             transaction.ID,
		Description:    transaction.Description,
		CreatedAt:      transaction.CreatedAt.UTC().Format(time.RFC3339Nano),
		OriginalValue:  transaction.Value,
		ConversionRate: math.Round(floatExchangeRate*100) / 100,
		ConvertedValue: math.Round(transaction.Value*floatExchangeRate*100) / 100,
	}

	json.NewEncoder(w).Encode(transactionOutput)
}

// GetAllTransactionsPaginated godoc
// @Summary Get all transactions with pagination
// @Description Get a paginated list of transactions
// @Tags transactions
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} dto.TransactionsPaginated
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /transactions [get]
func (h *TransactionHandler) GetAllTransactionsPaginated(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	pageSize := r.URL.Query().Get("pageSize")
	if pageSize == "" {
		pageSize = "10"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	items, err := h.TransactionRepository.GetAllPaginated(pageInt, pageSizeInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(items)
}
