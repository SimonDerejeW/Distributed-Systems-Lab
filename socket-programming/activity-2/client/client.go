package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	go receiveMessages(conn) // Receive messages in a separate goroutine

	for {
		message, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		conn.Write([]byte(message))
	}
}

func receiveMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Disconnected from server.")
			return // Exit the loop if there's an error (like a closed connection)
		}
		fmt.Print("Message from server: ", message)
	}
}
