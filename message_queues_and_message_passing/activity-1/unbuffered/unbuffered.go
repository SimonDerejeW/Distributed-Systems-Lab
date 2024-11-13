// unbuffered.go
package main

import (
	"fmt"
	"time"
)

// Producer function: sends messages to the channel
func producer(ch chan<- string) {
	for i := 1; i <= 5; i++ {
		ch <- fmt.Sprintf("Message %d", i) // Send message to channel
		time.Sleep(1 * time.Second)        // Sleep for 1 second
	}
	close(ch) // Close the channel after sending all messages
}

// Consumer function: receives and prints messages from the channel
func consumer(ch <-chan string) {
	for msg := range ch { // Read from channel until it's closed
		fmt.Println("Received:", msg)
	}
}

func main() {
	ch := make(chan string) // Create an unbuffered channel

	go producer(ch) // Start producer in a new goroutine
	consumer(ch)    // Run consumer in the main thread
}
