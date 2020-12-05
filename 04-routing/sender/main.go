package main

import (
	"log"
	"os"
	"strings"

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
		"logs_direct",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare a exchange")

	body := getContent(os.Args)
	log.Println(body)
	err = ch.Publish(
		"logs_direct",          // exchange
		getRoutingKey(os.Args), // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			// queue won't be lost even if RabbitMQ restarts
			//If you need a stronger guarantee then you can use publisher confirms
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish a message")

}

func getContent(args []string) (content string) {

	if len(args) < 3 || args[2] == "" {
		content = "hello"
	} else {
		content = strings.Join(args[2:], " ")
	}
	return
}

func getRoutingKey(args []string) (key string) {

	if len(args) < 2 || args[1] == "" {
		key = "info"
	} else {
		key = args[1]
	}
	return
}
