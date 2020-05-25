package main

import (
	"context"
	"flag"
	"greet/greetpb"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

// struct that will implement gRPC interfaces
type server struct{}

// Greet is unary API
func (*server) Greet(ctx context.Context, req *greetpb.GreetingRequest) (*greetpb.GreetingResponse, error) {
	glog.Infof("Greet RPC was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	result := "Hello " + firstName + " " + lastName
	return &greetpb.GreetingResponse{
		Result: result,
	}, nil
}

// GreetManyTimes is streaming gRPC
func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	glog.Infof("Greet many times RPC was invoked with %v\n", req)
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
	glog.Infof("All data has been streamed successfully.")
	return nil
}

func (*server) LongGreet(serverStream greetpb.GreetService_LongGreetServer) error {
	glog.Info("Long greet RPC was invoked.")
	var recipient []string
	for {
		msg, err := serverStream.Recv()
		if err == io.EOF {
			glog.Info("all data recieved.")
			serverStream.SendAndClose(&greetpb.LongGreetResponse{
				Result: "Hello " + strings.Join(recipient[:], ","),
			})
			return nil
		}
		recipient = append(recipient, msg.GetGreeting().GetFirstName())
		if err != nil {
			glog.Fatalf("Error connecting stream: %v", err)
		}
	}
}

func main() {
	flag.Parse()
	flag.Set("logtostderr", "true")

	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		glog.Fatalf("Failed to listen: %v", err)
	}
	glog.Infof("gRPC server is running on 0.0.0.0:50051.")

	grpcServer := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(grpcServer, &server{})

	if err := grpcServer.Serve(listener); err != nil {
		glog.Fatalf("Failed to server: %v", err)
	}
}
