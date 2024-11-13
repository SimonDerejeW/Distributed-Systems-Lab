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

	// Declare the same queue
	q, err := ch.QueueDeclare(
		"task_queue", // Queue name
		true,         // Durable
		false,        // Delete when unused
		false,        // Exclusive
		false,        // No-wait
		nil,          // Arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Start consuming messages
	msgs, err := ch.Consume(
		q.Name, // Queue name
		"",     // Consumer tag
		false,  // Auto-acknowledge set to false for manual acks
		false,  // Exclusive
		false,  // No-local
		false,  // No-wait
		nil,    // Arguments
	)
	failOnError(err, "Failed to register a consumer")

	// Goroutine to process incoming messages
	go func() {
		for d := range msgs {
			fmt.Printf("Received a message: %s\n", d.Body)
			d.Ack(false) // Manually acknowledge message after processing
		}
	}()

	fmt.Println("Waiting for messages...")
	select {} // Keep the program running
}
