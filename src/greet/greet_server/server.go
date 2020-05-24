package main

import (
	"greet/greetpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(grpcServer, &server{})

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
