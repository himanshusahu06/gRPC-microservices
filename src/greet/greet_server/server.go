package main

import (
	"context"
	"flag"
	"greet/greetpb"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

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

// GreetManyTimes is server streaming gRPC
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

// LongGreet is client streaming gRPC
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

// GreetEveryone is bidirectional streaming gRPC
func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	glog.Infoln("Greet everyone RPC was invoked.")
	for {
		req, recvErr := stream.Recv()
		if recvErr == io.EOF {
			return nil
		}
		if recvErr != nil {
			glog.Fatalf("Error receiving client stream: %v", recvErr)
			return recvErr
		}
		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: "Hello " + req.GetGreeting().GetFirstName() + "!",
		})
		if sendErr != nil {
			log.Fatalf("Error sending response: %v", sendErr)
			return sendErr
		}
	}
}

// GreetWithDeadlines is unary API with deadlines
func (*server) GreetWithDeadlines(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	glog.Infof("Greet with deadline RPC was invoked with %v\n", req)
	// keep checking if client has canceled the request
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			// client canceled the request
			glog.Info("Client has canceled the request!")
			return nil, status.Errorf(codes.Canceled, "Client has cancelled the request")
		}
		time.Sleep(time.Second)
	}
	firstName := req.GetGreeting().GetFirstName()
	return &greetpb.GreetWithDeadlineResponse{
		Result: "Hello " + firstName + "!",
	}, nil
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
	// Register reflection service on the given gRPC server.
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listener); err != nil {
		glog.Fatalf("Failed to server: %v", err)
	}
}
