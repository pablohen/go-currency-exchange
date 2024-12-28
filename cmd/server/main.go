package main

import (
	"encoding/json"
	"net/http"

	"go-currency-exchange/configs"
	"go-currency-exchange/internal/entity"
	"go-currency-exchange/internal/infra/database"
	"go-currency-exchange/internal/infra/webserver/handlers"
	"go-currency-exchange/pkg/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	dsn := "host=" + config.DBHost + " user=" + config.DBUser + " password=" + config.DBPassword + " dbname=" + config.DBName + " port=" + config.DBPort + " sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.Transaction{})
	if err != nil {
		panic(err)
	}
	println("Database migrated")

	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	transactionDB := database.NewTransaction(db)
	transactionHandler := handlers.NewTransactionHandler(transactionDB)

	http.HandleFunc("GET /transactions", transactionHandler.GetAllTransactionsPaginated)
	http.HandleFunc("POST /transactions", transactionHandler.CreateTransaction)
	http.HandleFunc("GET /transactions/{id}", transactionHandler.GetTransactionById)

	go http.ListenAndServe(":"+config.WebServerPort, nil)
	println("Running at port: " + config.WebServerPort)

	rabbitmqMessagesChannel := make(chan amqp.Delivery)
	go rabbitmq.Consume(ch, rabbitmqMessagesChannel)
	rabbitmqWorker(rabbitmqMessagesChannel, transactionDB)
}

func rabbitmqWorker(messageChan chan amqp.Delivery, transactionDB *database.Transaction) {
	for message := range messageChan {
		println(string(message.Body))

		var transaction entity.Transaction
		err := json.Unmarshal(message.Body, &transaction)
		if err != nil {
			panic(err)
		}

		err = transactionDB.Create(transaction.Description, transaction.Value)
		if err != nil {
			panic(err)
		}
		message.Ack(false)
	}
}
