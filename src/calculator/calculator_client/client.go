package main

import (
	"calculator/calculatorpb"
	"context"
	"flag"
	"io"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

func main() {
	flag.Parse()
	flag.Set("logtostderr", "true")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		glog.Fatalf("Failed to connect: %v\n", err)
	}
	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)
	//doUnary(c)
	//doPrimeDecomposeStreaming(c)
	doComputeAvarege(c)
}

func doUnary(client calculatorpb.CalculatorServiceClient) {
	var numbers []int64
	numbers = append(numbers, 4)
	numbers = append(numbers, 3)
	calculatorResponse, err := client.Sum(context.Background(), &calculatorpb.CalculatorRequest{
		Numbers: numbers,
	})
	if err != nil {
		glog.Fatalf("Failed to connect to server: %v", err)
	}
	glog.Infof("response: %v\n", calculatorResponse)
}

func doPrimeDecomposeStreaming(client calculatorpb.CalculatorServiceClient) {
	var number int64 = 120
	stream, err := client.DecomposePrime(context.Background(), &calculatorpb.PrimeNumberDecompositionRequest{
		Number: number,
	})
	if err != nil {
		glog.Fatalf("error invoking prime decomposition times RPC: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			glog.Fatalf("error connecting stream: %v", err)
		}
		glog.Infof("One of the prime decomposition of %d is %d\n", number, msg.Result)
	}
}

func doComputeAvarege(client calculatorpb.CalculatorServiceClient) {
	numbers := []*calculatorpb.AverageNumberRequest{
		&calculatorpb.AverageNumberRequest{
			Number: 1,
		},
		&calculatorpb.AverageNumberRequest{
			Number: 2,
		},
		&calculatorpb.AverageNumberRequest{
			Number: 3,
		},
		&calculatorpb.AverageNumberRequest{
			Number: 4,
		},
	}

	stream, err := client.ComputeAverage(context.Background())
	if err != nil {
		glog.Fatalf("Error connecting server: %v", err)
	}
	for _, number := range numbers {
		glog.Infof("Sending number: %v\n", number)
		stream.Send(number)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		glog.Fatalf("Error receiveing response from server: %v\n", err)
	}
	glog.Infof("Average of %v is %v\n", numbers, res.GetResult())
}
