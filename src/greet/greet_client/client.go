package main

import (
	"fmt"
	"greet/greetpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("I'm client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	fmt.Printf("Created client: %f", c)
}
