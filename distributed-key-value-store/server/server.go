package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	store = make(map[string]string)
	mutex = &sync.Mutex{}
)

// File path for the JSON store
const jsonFile = "store.json"

// LoadStore reads the key-value store from the JSON file
func LoadStore() error {
	file, err := os.Open(jsonFile)
	if err != nil {
		// If the file does not exist, initialize an empty store
		if os.IsNotExist(err) {
			store = make(map[string]string)
			return nil
		}
		return fmt.Errorf("could not open JSON file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&store)
	if err != nil {
		return fmt.Errorf("could not decode JSON file: %v", err)
	}
	return nil
}

// SaveStore writes the key-value store to the JSON file
func SaveStore() error {
	file, err := os.Create(jsonFile)
	if err != nil {
		return fmt.Errorf("could not create JSON file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(store)
	if err != nil {
		return fmt.Errorf("could not encode JSON to file: %v", err)
	}
	return nil
}


// handleClient processes commands from the client
func handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		// Read input from the client
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected:", err)
			return
		}

		// Process the message
		command := strings.TrimSpace(message)
		response := processCommand(command)

		// Send the response back to the client
		conn.Write([]byte(response + "\n"))
	}
}

// processCommand handles the commands sent from the client
func processCommand(command string) string {
	parts := strings.Fields(command)
	if len(parts) < 1 {
		return "ERROR: Invalid command"
	}

	switch strings.ToUpper(parts[0]) {
	case "PUT":
		if len(parts) != 3 {
			return "ERROR: Usage PUT <key> <value>"
		}
		key := parts[1]
		value := parts[2]
		mutex.Lock()
		store[key] = value
		err := SaveStore()
		mutex.Unlock()
		if err != nil {
			return "ERROR: Could not save data"
		}
		return "OK: Key added/updated"

	case "GET":
		if len(parts) != 2 {
			return "ERROR: Usage GET <key>"
		}
		key := parts[1]
		mutex.Lock()
		value, exists := store[key]
		mutex.Unlock()
		if exists {
			return fmt.Sprintf("VALUE: %s", value)
		}
		return "ERROR: Key not found"

	case "DELETE":
		if len(parts) != 2 {
			return "ERROR: Usage DELETE <key>"
		}
		key := parts[1]
		mutex.Lock()
		_, exists := store[key]
		if exists {
			delete(store, key)
			err := SaveStore()
			mutex.Unlock()
			if err != nil {
				return "ERROR: Could not save data"
			}
			return "OK: Key deleted"
		}
		mutex.Unlock()
		return "ERROR: Key not found"

	case "LIST":
		mutex.Lock()
		if len(store) == 0 {
			mutex.Unlock()
			return "EMPTY: No key-value pairs stored"
		}
		var result []string
		for k, v := range store {
			result = append(result, fmt.Sprintf("%s: %s", k, v))
		}
		mutex.Unlock()
		return strings.Join(result, ", ")

	default:
		return "ERROR: Unknown command"
	}
}

// startServer initializes the TCP server
func startServer() {
	err := LoadStore()
	if err != nil {
		fmt.Println("Error loading store from JSON:", err)
		return
	}
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is running on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleClient(conn)
	}
}

func main() {
	startServer()
}
