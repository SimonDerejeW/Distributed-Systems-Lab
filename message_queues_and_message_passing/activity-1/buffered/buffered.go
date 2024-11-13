// buffered.go
package main

import (
	"fmt"
	"time"
)

// Producer function: sends messages to the channel
func producer(ch chan<- string) {
	for i := 1; i <= 5; i++ {
		fmt.Printf("Sending: Message %d\n", i)
		ch <- fmt.Sprintf("Message %d", i) // Send message to channel
	}
	close(ch) // Close the channel after sending all messages
}

// Consumer function: receives and prints messages from the channel
func consumer(ch <-chan string) {
	for msg := range ch { // Read from channel until it's closed
		fmt.Println("Received:", msg)
		time.Sleep(2 * time.Second) // Simulate processing time
	}
}

func main() {
	ch := make(chan string, 3) // Buffered channel with a capacity of 3

	go producer(ch) // Start producer in a new goroutine
	consumer(ch)    // Run consumer in the main thread
}
