package main

import (
	"encoding/json"
	"log"
	"net/http"

	"go-currency-exchange/configs"
	"go-currency-exchange/internal/entity"
	"go-currency-exchange/internal/infra/database"
	"go-currency-exchange/internal/infra/webserver/handlers"
	"go-currency-exchange/pkg/rabbitmq"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.SetFlags(log.LstdFlags | log.LUTC)
	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(log.Writer(), "", log.LstdFlags|log.LUTC)})

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
	log.Printf("Database migrated")

	rabbitmqChannel, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer rabbitmqChannel.Close()

	transactionRepository := database.NewTransactionRepository(db)
	transactionHandler := handlers.NewTransactionHandler(transactionRepository, rabbitmqChannel)

	transactionRouter := chi.NewRouter()
	transactionRouter.Get("/", transactionHandler.GetAllTransactionsPaginated)
	transactionRouter.Post("/", transactionHandler.CreateTransaction)
	transactionRouter.Get("/{id}", transactionHandler.GetTransactionById)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Mount("/transactions", transactionRouter)
	go http.ListenAndServe(":"+config.WebServerPort, r)
	log.Printf("Running at port: %s", config.WebServerPort)

	transactionMessagesChannel := make(chan amqp.Delivery)
	go rabbitmq.Consume(rabbitmqChannel, "transactions", "", transactionMessagesChannel)

	createTransactionWorker(transactionMessagesChannel, transactionRepository)
}

func createTransactionWorker(messageChan chan amqp.Delivery, transactionRepository *database.TransactionRepository) {
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
