package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	for {
		// Receive task (number) from server
		task, _ := bufio.NewReader(conn).ReadString('\n')
		task = strings.TrimSpace(task)

		num, _ := strconv.Atoi(task)

		result := num * num
		fmt.Printf("Received task: %d, computed result: %d\n", num, result)

		// Send result back to server
		fmt.Fprintf(conn, "%d\n", result)
	}
}
