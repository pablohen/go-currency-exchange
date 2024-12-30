package worker

import (
	"encoding/json"
	"log"

	"go-currency-exchange/internal/entity"
	"go-currency-exchange/internal/infra/database"

	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateTransaction(messageChan chan amqp.Delivery, transactionRepository *database.TransactionRepository) {
	for message := range messageChan {
		log.Printf("Received message: %s", message.Body)

		var transaction entity.Transaction
		err := json.Unmarshal(message.Body, &transaction)
		if err != nil {
			// TODO: log and push to error queue
			panic(err)
		}

		err = transactionRepository.Create(transaction.Description, transaction.Value)
		if err != nil {
			// TODO: log and push to error queue
			panic(err)
		}
		message.Ack(false)
	}
}
