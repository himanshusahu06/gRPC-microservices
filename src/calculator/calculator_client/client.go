package main

import (
	"calculator/calculatorpb"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v\n", err)
	}
	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)
	doUnary(c)
}

func doUnary(client calculatorpb.CalculatorServiceClient) {
	var numbers []int64
	numbers = append(numbers, 4)
	numbers = append(numbers, 3)
	calculatorResponse, err := client.Sum(context.Background(), &calculatorpb.CalculatorRequest{
		Numbers: numbers,
	})
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	fmt.Printf("response: %v\n", calculatorResponse)
}
