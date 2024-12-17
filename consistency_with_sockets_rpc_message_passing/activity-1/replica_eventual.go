package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Replica struct {
	data  map[string]string
	mu    sync.Mutex
	peers []string
}

// Update modifies the local data of the replica
func (r *Replica) Update(key, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[key] = value
}

// propagateUpdates sends key-value updates to all peers with simulated delay
func (r *Replica) propagateUpdates(key, value string) {
	startTime := time.Now()
	for _, peer := range r.peers {
		go func(peer string) {
			// Simulate network delay
			time.Sleep(time.Millisecond * 100) // Adjust delay as needed
			conn, err := net.Dial("tcp", peer)
			if err != nil {
				fmt.Println("Error connecting to peer:", peer, err)
				return
			}
			defer conn.Close()
			fmt.Fprintf(conn, "%s:%s\n", key, value)
		}(peer)
	}
	// Log the time it takes to start sending updates
	fmt.Printf("Propagation initiated: %v elapsed since start\n", time.Since(startTime))
}

// handleConnection handles incoming messages from other replicas
func handleConnection(conn net.Conn, replica *Replica) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		parts := strings.Split(strings.TrimSpace(message), ":")
		if len(parts) == 2 {
			replica.Update(parts[0], parts[1])
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run replica_eventual.go <machine_ip:port> <peer1_ip:port> [<peer2_ip:port>...]")
		return
	}

	// Parse command-line arguments
	machineAddr := os.Args[1]
	peers := os.Args[2:]

	// Initialize the replica
	replica := &Replica{
		data:  make(map[string]string),
		peers: peers,
	}

	// Start the server
	listener, err := net.Listen("tcp", machineAddr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Printf("Replica listening on %s\n", machineAddr)

	// Handle incoming connections
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go handleConnection(conn, replica)
		}
	}()

	// Simulate an update
	replica.Update("key1", "value1")
	startPropagation := time.Now()
	replica.propagateUpdates("key1", "value1")

	// Wait for updates to propagate
	time.Sleep(5 * time.Second)

	// Check final state and log convergence time
	replica.mu.Lock()
	fmt.Println("Replica Data:", replica.data)
	convergenceTime := time.Since(startPropagation)
	fmt.Printf("Convergence Time: %v\n", convergenceTime)
	replica.mu.Unlock()
}
