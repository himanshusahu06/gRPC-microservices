package main

import (
	"calculator/calculatorpb"
	"context"
	"flag"
	"io"
	"log"
	"math"
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

func (*server) ComputeAverage(serverStream calculatorpb.CalculatorService_ComputeAverageServer) error {
	glog.Infof("Compute Average RPC was recieved.")
	var sum int64 = 0
	var count int64 = 0
	for {
		msg, err := serverStream.Recv()
		if err == io.EOF {
			serverStream.SendAndClose(&calculatorpb.AverageNumberResponse{
				Result: float64(sum) / float64(count),
			})
			glog.Infoln("All request stream recieved.")
			break
		}
		glog.Infof("Stream request recieved: %v\n", msg)
		sum = sum + msg.GetNumber()
		count++
	}
	return nil
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	glog.Infof("Find maxmimum RPC was recieved.")
	var runningMax int64 = math.MinInt64
	for {
		msg, recvErr := stream.Recv()
		if recvErr == io.EOF {
			return nil
		}
		if recvErr != nil {
			glog.Fatalf("Failed to receive messages from client: %v", recvErr)
			return recvErr
		}
		if runningMax < msg.GetNumber() {
			runningMax = msg.GetNumber()
			sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
				Result: runningMax,
			})
			if sendErr != nil {
				log.Fatalf("Error sending response: %v", sendErr)
				return sendErr
			}
		}
	}
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
