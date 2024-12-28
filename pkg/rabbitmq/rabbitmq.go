package rabbitmq

import (
	"net/url"

	"go-currency-exchange/configs"

	amqp "github.com/rabbitmq/amqp091-go"
)

func OpenChannel() (*amqp.Channel, error) {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	url := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(config.RabbitMQUser, config.RabbitMQPassword),
		Host:   config.RabbitMQHost + ":" + config.RabbitMQPort,
	}
	println(url.String())

	conn, err := amqp.Dial(url.String())
	if err != nil {
		return nil, err
	}
	println("Connected to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func Consume(ch *amqp.Channel, out chan amqp.Delivery) error {
	queue := "transactions"
	consumer := "go-currency-exchange-consumer"

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
		consumer,
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
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
