package main

import (
	"context"
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
	//fmt.Printf("Created client: %f\n", c)
	doUnary(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	// invoking RPC
	greetRequest := &greetpb.GreetingRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Himanshu",
			LastName:  "Sahu",
		},
	}
	greetResponse, err := c.Greet(context.Background(), greetRequest)
	if err != nil {
		log.Fatalf("error invoking greet RPC: %v", err)
	}
	fmt.Printf("greet response %v", greetResponse)
}
