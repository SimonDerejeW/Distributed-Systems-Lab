package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type Args struct {
	A, B int
}

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Error connecting to RPC server: ",err)
	}
	args := Args{A: 3, B: 5}


	var muliplyReply int
	err = client.Call("Calculator.Multiply", args, &muliplyReply)
	if err != nil {
		log.Fatal("Error calling RPC: ",err)
	}
	fmt.Printf("Result of %d * %d = %d\n", args.A, args.B, muliplyReply)

	var addReply int
	err = client.Call("Calculator.Add", args, &addReply)
	if err != nil {
		log.Fatal("Error calling RPC: ",err)
	}
	fmt.Printf("Result of %d + %d = %d\n", args.A, args.B, addReply)

	var subtractReply int
	err = client.Call("Calculator.Subtract", args, &subtractReply)
	if err != nil {
		log.Fatal("Error calling RPC: ",err)
	}
	fmt.Printf("Result of %d - %d = %d\n", args.A, args.B, subtractReply)

	var divideReply float32
	err = client.Call("Calculator.Divide", args, &divideReply)
	if err != nil {
		log.Fatal("Error calling RPC: ",err)
	}
	fmt.Printf("Result of %d / %d = %.2f\n", args.A, args.B, divideReply)
}