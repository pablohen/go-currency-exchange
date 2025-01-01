package rabbitmq

import (
	"log"
	"net/url"
	"time"

	"go-currency-exchange/configs"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect() (*amqp.Channel, error) {
	var rabbitmqChannel *amqp.Channel
	var err error
	for {
		rabbitmqChannel, err = getChannel()
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ, retrying in 2 seconds: %v", err)
		time.Sleep(2 * time.Second)
	}
	return rabbitmqChannel, err
}

func getChannel() (*amqp.Channel, error) {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	url := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(config.RabbitMQUser, config.RabbitMQPassword),
		Host:   config.RabbitMQHost + ":" + config.RabbitMQPort,
	}

	conn, err := amqp.Dial(url.String())
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func Consume(ch *amqp.Channel, queue string, out chan amqp.Delivery) error {
	_, err := ch.QueueDeclare(queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	messages, err := ch.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for message := range messages {
		out <- message
	}

	return nil
}

func Publish(ch *amqp.Channel, queue string, message string) error {
	err := ch.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			Body: []byte(message),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
