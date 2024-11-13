// synchronization.go
package main

import (
	"fmt"
	"time"
)

// Producer function: generates messages and listens for a quit signal
func producer(ch chan<- string, quit <-chan bool) {
	for i := 1; ; i++ {
		select {
		case <-quit: // If quit signal is received, stop producing
			fmt.Println("Producer shutting down")
			return
		case ch <- fmt.Sprintf("Message %d", i): // Send message to channel
			fmt.Printf("Produced: Message %d\n", i)
			time.Sleep(1 * time.Second) // Sleep for 1 second
		}
	}
}

// Consumer function: consumes messages and sends a quit signal after consuming 10 messages
func consumer(ch <-chan string, quit chan<- bool) {
	for i := 0; i < 10; i++ {
		fmt.Println("Consumed:", <-ch) // Receive message from channel
	}
	quit <- true // Send shutdown signal to producer
}

func main() {
	ch := make(chan string)   // Create a string channel
	quit := make(chan bool)   // Create a quit signal channel

	go producer(ch, quit)     // Start producer in a new goroutine
	go consumer(ch, quit)     // Start consumer in a new goroutine

	<-quit                    // Wait for shutdown signal from the consumer
	fmt.Println("Main shutting down")
}
