package main

import (
	"net/http"

	"go-currency-exchange/configs"
	"go-currency-exchange/internal/entity"
	"go-currency-exchange/internal/infra/database"
	"go-currency-exchange/internal/infra/webserver/handlers"

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

	println("Running at port: " + config.WebServerPort)
	transactionDB := database.NewTransaction(db)
	transactionHandler := handlers.NewTransactionHandler(transactionDB)

	http.HandleFunc("GET /transactions", transactionHandler.GetAllTransactionsPaginated)
	http.HandleFunc("POST /transactions", transactionHandler.CreateTransaction)
	http.HandleFunc("GET /transactions/{id}", transactionHandler.GetTransactionById)

	err = http.ListenAndServe(":"+config.WebServerPort, nil)
	if err != nil {
		panic(err)
	}
}
