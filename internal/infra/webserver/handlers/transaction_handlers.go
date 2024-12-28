package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"time"

	"go-currency-exchange/internal/dto"
	"go-currency-exchange/internal/entity"
	"go-currency-exchange/internal/infra/database"
	"go-currency-exchange/pkg/rabbitmq"
	"go-currency-exchange/pkg/treasury"
)

type TransactionHandler struct {
	TransactionDB database.TransactionInterface
}

func NewTransactionHandler(db database.TransactionInterface) *TransactionHandler {
	return &TransactionHandler{
		TransactionDB: db,
	}
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var createTransactionInput dto.CreateTransactionInput
	err := json.NewDecoder(r.Body).Decode(&createTransactionInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}

	transaction, err := entity.NewTransaction(createTransactionInput.Description, createTransactionInput.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = rabbitmq.Publish(ch, "transactions", string(transactionJSON))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *TransactionHandler) GetTransactionById(w http.ResponseWriter, r *http.Request) {
	transactionID := r.PathValue("id")
	if transactionID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	transaction, err := h.TransactionDB.GetById(transactionID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	rate, err := treasury.GetRatesByDate(transaction.CreatedAt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	floatExchangeRate, err := strconv.ParseFloat(rate.ExchangeRate, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	transactionOutput := dto.TransactionOutput{
		ID:             transaction.ID,
		Description:    transaction.Description,
		CreatedAt:      transaction.CreatedAt.UTC().Format(time.RFC3339),
		OriginalValue:  transaction.Value,
		ConversionRate: math.Round(floatExchangeRate*100) / 100,
		ConvertedValue: math.Round(transaction.Value*floatExchangeRate*100) / 100,
	}

	json.NewEncoder(w).Encode(transactionOutput)
}

func (h *TransactionHandler) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := h.TransactionDB.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transactions)
}

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

	items, err := h.TransactionDB.GetAllPaginated(pageInt, pageSizeInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(items)
}
