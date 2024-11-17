package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
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

	
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("Enter command: ")
			scanner.Scan()
			text := scanner.Text()
			
			if strings.ToUpper(text) == "EXIT" {
				fmt.Println("Exiting client...")
				return
			}

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

}

func main() {
	startClient()
}
