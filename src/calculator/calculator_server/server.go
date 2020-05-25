package main

import (
	"calculator/calculatorpb"
	"context"
	"flag"
	"net"

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
