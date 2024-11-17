package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// startClient connects to the server and handles user input
func startClient() {
	conn, err := net.Dial("tcp", "localhost:8080")
	
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to server. You can start sending commands (PUT, GET, DELETE, LIST)")
	
	timer := time.NewTimer(30 * time.Second)
	// Channel to signal user input activity
	reset := make(chan bool)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("Enter command: ")
			scanner.Scan()
			text := scanner.Text()
			
			if strings.ToUpper(text) == "EXIT" {
				fmt.Println("Exiting client...")
				return
			}

			reset <- true
	
			_, err := conn.Write([]byte(text + "\n"))
			if err != nil {
				fmt.Println("Error sending message:", err)
				return
			}
	
			// Read response from server
			response, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println("Error reading response:", err)
				return
			}
			fmt.Println("Server response:", strings.TrimSpace(response))
		}

	}()

	for {
		select {
		case <-timer.C:
			// Timer expired, exit the program
			fmt.Println("No activity for 30 seconds. Exiting...")
			return
		case <-reset:
			// User input received, reset the timer
			if !timer.Stop() {
				<-timer.C // Drain the channel to prevent blocking
			}
			// Restart the timer
			timer.Reset(30 * time.Second)
		}
	}
}

func main() {
	startClient()
}
