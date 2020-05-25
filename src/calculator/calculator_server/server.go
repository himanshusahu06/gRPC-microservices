package main

import (
	"calculator/calculatorpb"
	"context"
	"flag"
	"net"
	"time"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	glog.Infof("Sum RPC was recieved: %v\n", req)
	var result int64 = 0
	for _, element := range req.GetNumbers() {
		result = result + element
	}
	return &calculatorpb.CalculatorResponse{
		Result: result,
	}, nil
}

func (*server) DecomposePrime(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_DecomposePrimeServer) error {
	glog.Infof("Decompose RPC was recieved: %v\n", req)
	number := req.GetNumber()
	var divisor int64 = 2
	for {
		if number <= 1 {
			break
		}
		if number%divisor == int64(0) {
			// prime number found, send it to client
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				Result: divisor,
			})
			number = number / divisor
			time.Sleep(1000 * time.Millisecond)
		} else {
			divisor = divisor + 1
		}
	}
	return nil
}

func main() {
	flag.Parse()
	flag.Set("logtostderr", "true")

	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		glog.Fatalf("Failed to listen: %v\n", err)
	}
	glog.Info("Listening on 0.0.0.0:500051")

	grpcServer := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(grpcServer, &server{})

	if err := grpcServer.Serve(listener); err != nil {
		glog.Fatalf("Failed to server: %v", err)
	}
}
