package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"sync"
)

// Replica represents a node in the distributed system
type Replica struct {
	data    map[string]string
	mu      sync.Mutex
	peers   []string // List of peer addresses
	ackLock sync.Mutex
	acks    map[string]int // Track acknowledgments
}

// Args defines the key-value pair for updates
type Args struct {
	Key   string
	Value string
}

// Update is an RPC method to apply an update on the replica
func (r *Replica) Update(args *Args, reply *bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[args.Key] = args.Value
	*reply = true // Indicate success
	return nil
}

// propagateUpdates sends updates to all peers using RPC
func (r *Replica) propagateUpdates(key, value string) {
	r.ackLock.Lock()
	r.acks[key] = 0 // Initialize acknowledgment count for this key
	r.ackLock.Unlock()

	for _, peer := range r.peers {
		go func(peer string) {
			client, err := rpc.Dial("tcp", peer)
			if err != nil {
				fmt.Println("Error connecting to peer:", peer, err)
				return
			}
			defer client.Close()

			args := &Args{Key: key, Value: value}
			var reply bool
			err = client.Call("Replica.Update", args, &reply)
			if err == nil && reply {
				r.ackLock.Lock()
				r.acks[key]++ // Increment acknowledgment count on success
				r.ackLock.Unlock()
			}
		}(peer)
	}
}

// waitForAcknowledgments waits for the required number of acknowledgments
func (r *Replica) waitForAcknowledgments(key string, required int) {
	for {
		r.ackLock.Lock()
		if r.acks[key] >= required {
			r.ackLock.Unlock()
			break
		}
		r.ackLock.Unlock()
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run replica_strong.go <machine_ip:port> <peer1_ip:port> [<peer2_ip:port>...]")
		return
	}

	// Parse command-line arguments
	machineAddr := os.Args[1]
	peers := os.Args[2:]

	// Initialize the replica
	replica := &Replica{
		data:  make(map[string]string),
		peers: peers,
		acks:  make(map[string]int),
	}
	rpc.Register(replica)

	// Start the RPC server
	listener, err := net.Listen("tcp", machineAddr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Printf("Replica RPC Server listening on %s\n", machineAddr)

	// Accept connections and serve RPC requests
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go rpc.ServeConn(conn)
		}
	}()

	// Simulate a strong consistency update
	key, value := "key1", "value1"
	replica.Update(&Args{Key: key, Value: value}, nil)
	replica.propagateUpdates(key, value)
	replica.waitForAcknowledgments(key, len(replica.peers))
	fmt.Println("Update committed after receiving acknowledgments")
}
