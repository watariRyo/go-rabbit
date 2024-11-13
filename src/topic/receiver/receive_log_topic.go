package main

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", err, msg)
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
		"logs_topic", // name
		"topic", // type
		true,   // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil, // arguments
	)
	failOnError(err, "Failed to declare a exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [info] [warning] [error]", os.Args[0])
		os.Exit(0)
	}
	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
		   q.Name, "logs_direct", s)
		err = ch.QueueBind(
		  q.Name,        // queue name
		  s,             // routing key
		  "logs_topic", // exchange
		  false,
		  nil)
		failOnError(err, "Failed to bind a queue")
	}

	msgs, err := ch.Consume(
		q.Name, // queue
        "",     // consumer
        false,   // auto-ack
        false,   // exclusive
        false,   // no-local
        false,   // no-wait
        nil,   // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
        for d := range msgs {
            log.Printf("Received a message: %s", d.Body)
        }
    }()

	log.Printf("[*] Waiting for messages. To exit press CTRL+C\n")
	<-forever
}