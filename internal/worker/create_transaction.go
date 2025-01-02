package worker

import (
	"encoding/json"
	"log"

	"go-currency-exchange/internal/dto"
	"go-currency-exchange/internal/infra/database"

	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateTransaction(messageChan chan amqp.Delivery, transactionRepository database.TransactionInterface) {
	for message := range messageChan {
		log.Printf("Received message: %s", message.Body)

		var transactionMessage dto.TransactionMessage
		err := json.Unmarshal(message.Body, &transactionMessage)
		if err != nil {
			// TODO: push to error queue
			log.Printf("Error unmarshalling message: %s", message.Body)
		}

		err = transactionRepository.Create(transactionMessage.Description, transactionMessage.Value, transactionMessage.CreatedAt)
		if err != nil {
			// TODO: push to error queue
			log.Printf("Error creating transaction: %s", message.Body)
		}
		message.Ack(false)
	}
}
