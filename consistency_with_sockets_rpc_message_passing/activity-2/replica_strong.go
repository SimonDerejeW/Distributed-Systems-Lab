package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"sync"
	"time" // Importing "time"
)

type Replica struct {
	data    map[string]string
	mu      sync.Mutex
	peers   []string
	ackLock sync.Mutex
	acks    map[string]int // track acknowledgements
}

type Args struct {
	Key    string
	Value  string
	Source string // Adding "Source" field
}

func (r *Replica) Update(args *Args, reply *bool) error {
	fmt.Println("> [UPDATE] Update request committed with args: " + args.Key + "<->" + args.Value + " sourced from: " + args.Source) // Printing source
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[args.Key] = args.Value
	if reply != nil {
		*reply = true
	}
	return nil
}

func (r *Replica) propagateUpdates(key, value, machineAddress string) { // Adding "machineAddress" parameter
	time.Sleep(4 * time.Second) // Adding delay
	r.ackLock.Lock()
	r.acks[key] = 0
	r.ackLock.Unlock()

	for _, peer := range r.peers {
		go func(peer string) {
			client, err := rpc.Dial("tcp", peer)
			if err != nil {
				fmt.Println("Error connecting to peer:", peer, err)
				return
			}

			defer client.Close()
			args := &Args{Key: key, Value: value, Source: machineAddress} // Adding "Source" field
			var reply bool = false
			err = client.Call("Replica.Update", args, &reply)
			fmt.Println("> [PROPAGATION] Update request propagated to peer " + peer) // Printing propagation
			if err == nil && reply {
				fmt.Println("> [ACK] Peer " + peer + " acknowledged update request") // Printing acknowledgment
				r.ackLock.Lock()
				r.acks[key]++
				r.ackLock.Unlock()
			}
		}(peer)
	}
}

func (r *Replica) waitForAcknowledgements(key string, quorum int) { // Adding quorum parameter
	for {
		r.ackLock.Lock()
		if r.acks[key] >= quorum { // Using quorum
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

	machineAddr := os.Args[1]
	peers := os.Args[2:]

	replica := &Replica{
		data:  make(map[string]string),
		peers: peers,
		acks:  make(map[string]int),
	}

	rpc.Register(replica)

	listener, err := net.Listen("tcp", machineAddr)
	if err != nil {
		panic(err)
	}

	defer listener.Close()
	fmt.Printf("Replica RPC server listening on %s\n", machineAddr)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}

			go rpc.ServeConn(conn)
		}
	}()

	key, value := "key1", "value1"
	fmt.Println("> [INIT] Update initialized locally") // Printing initialization
	replica.Update(&Args{Key: key, Value: value, Source: "local"}, nil) // Adding "Source" field
	replica.propagateUpdates(key, value, machineAddr) // Adding "machineAddress" parameter

	var Q int = (len(replica.peers) / 2) // Calculating quorum
	replica.waitForAcknowledgements(key, Q) // Using quorum
	fmt.Println("> [COMMITED] Update committed after receiving sufficient acknowledgments") // Printing commit
	time.Sleep(5 * time.Second) // Adding delay
}
