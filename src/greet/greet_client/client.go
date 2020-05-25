package main

import (
	"context"
	"fmt"
	"greet/greetpb"
	"io"
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
	doServerStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting unary RPC client..")
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

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting server stream RPC client..")
	greetRequest := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Himanshu",
			LastName:  "Sahu",
		},
	}
	greetStreamResponse, err := c.GreetManyTimes(context.Background(), greetRequest)
	if err != nil {
		log.Fatalf("error invoking greet many times RPC: %v", err)
	}

	for {
		msg, err := greetStreamResponse.Recv()
		// when stream ends, client will get EOF
		if err == io.EOF {
			// when stream ends, client will get EOF
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v\n", err)
		}
		fmt.Println(msg)
	}
}
