package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// Helper function to handle errors
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(
		"task_queue", // Queue name
		true,         // Durable
		false,        // Delete when unused
		false,        // Exclusive
		false,        // No-wait
		nil,          // Arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Send multiple messages
	for i := 1; i <= 10; i++ {
		body := fmt.Sprintf("Task %d", i)
		err = ch.Publish(
			"",     // Exchange
			q.Name, // Routing key (queue name)
			false,  // Mandatory
			false,  // Immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent, // Mark message as persistent
				ContentType:  "text/plain",
				Body:         []byte(body),
			})
		failOnError(err, "Failed to publish a message")
		fmt.Println("Sent:", body)
	}
}
