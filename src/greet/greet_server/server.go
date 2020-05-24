package main

import (
	"context"
	"fmt"
	"greet/greetpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

// RPC implementation
func (*server) Greet(ctx context.Context, req *greetpb.GreetingRequest) (*greetpb.GreetingResponse, error) {
	fmt.Printf("Greet RPC was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	result := "Hello " + firstName + " " + lastName
	return &greetpb.GreetingResponse{
		Result: result,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	fmt.Println("gRPC server is running on 0.0.0.0:50051.")

	grpcServer := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(grpcServer, &server{})

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
