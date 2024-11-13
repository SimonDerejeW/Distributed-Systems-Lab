// nats_subscriber.go
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

	// Define the subject (topic) to subscribe to
	subject := "updates"

	// Subscribe to the subject and define the message handler
	_, err = nc.Subscribe(subject, func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Subscribed to 'updates'. Waiting for messages...")

	// Keep the program running to listen for messages
	select {}
}
