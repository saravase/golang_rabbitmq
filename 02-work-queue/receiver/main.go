package main

import (
	"log"
	"regexp"
	"strconv"
	"time"

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

	err = ch.Qos(
		1, // RabbitMQ not to give more than one message to a worker at a time
		0,
		false,
	)
	failOnError(err, "Failed to set Qos")

	q, err := ch.QueueDeclare(
		"wq",  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	p, err := regexp.Compile("[0-9]")
	failOnError(err, "Failed to compile regex")
	go func() {
		for d := range msgs {
			c, err := strconv.Atoi(p.FindString(string(d.Body)))
			failOnError(err, "Failed to string converstion")
			log.Printf("Received a message: %s", d.Body)
			t := time.Duration(c)
			if c == 3 {
				time.Sleep(time.Minute)
			} else {
				time.Sleep(t * time.Second)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
