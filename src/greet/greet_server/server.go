package main

import (
	"context"
	"fmt"
	"greet/greetpb"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

// struct that will implement gRPC interfaces
type server struct{}

// Greet is unary API
func (*server) Greet(ctx context.Context, req *greetpb.GreetingRequest) (*greetpb.GreetingResponse, error) {
	fmt.Printf("Greet RPC was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	result := "Hello " + firstName + " " + lastName
	return &greetpb.GreetingResponse{
		Result: result,
	}, nil
}

// GreetManyTimes is streaming gRPC
func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("Greet many times RPC was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	// just greet for 10 times
	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimesResponse{
			Result: "Hello " + firstName + " " + lastName + " " + strconv.Itoa(i),
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	fmt.Println("All data has been streamed successfully.")
	return nil
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
