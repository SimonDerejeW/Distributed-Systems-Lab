// nats_publisher.go
package main

import (
	"fmt"
	"log"
	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to the NATS server using the default URL (nats://localhost:4222)
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Define the subject (topic) and message
	subject := "updates"
	message := "Hello, NATS!"

	// Publish the message to the specified subject
	if err := nc.Publish(subject, []byte(message)); err != nil {
		log.Fatal(err)
	}

	// Print confirmation of message sent
	fmt.Println("Sent:", message)
}
