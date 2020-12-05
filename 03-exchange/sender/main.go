package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare a exchange")

	for i := 0; i < 10; i++ {

		body := fmt.Sprint("Data ", i)
		log.Println(body)
		err = ch.Publish(
			"logs", // exchange
			"",     // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				// queue won't be lost even if RabbitMQ restarts
				//If you need a stronger guarantee then you can use publisher confirms
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(body),
			})
		failOnError(err, "Failed to publish a message")
	}
}
