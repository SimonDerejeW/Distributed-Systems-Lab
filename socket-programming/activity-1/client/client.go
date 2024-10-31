package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
)

func main() {
    // Connect to the server on localhost port 8080
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println("Error connecting to server:", err)
        return
    }
    defer conn.Close()

    // Read a message from the command line
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter a message: ")
    message, _ := reader.ReadString('\n')

    // Send the message to the server
    fmt.Fprint(conn, message)

    // Wait for the response from the server
    response, _ := bufio.NewReader(conn).ReadString('\n')
    fmt.Println("Server response:", response)
}
