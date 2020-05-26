package main

import (
	"context"
	"flag"
	"greet/greetpb"
	"io"
	"time"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

func main() {
	flag.Parse()
	flag.Set("logtostderr", "true")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		glog.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	//fmt.Printf("Created client: %f\n", c)
	//doUnary(c)
	//doServerStreaming(c)
	// doClientStreaming(c)
	doBiDirectionalStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	glog.Infof("Starting unary RPC client..")
	// invoking RPC
	greetRequest := &greetpb.GreetingRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Himanshu",
			LastName:  "Sahu",
		},
	}
	greetResponse, err := c.Greet(context.Background(), greetRequest)
	if err != nil {
		glog.Fatalf("error invoking greet RPC: %v", err)
	}
	glog.Infof("greet response %v", greetResponse)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	glog.Infof("Starting server stream RPC client..")
	greetRequest := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Himanshu",
			LastName:  "Sahu",
		},
	}
	greetStreamResponse, err := c.GreetManyTimes(context.Background(), greetRequest)
	if err != nil {
		glog.Infof("error invoking greet many times RPC: %v", err)
	}

	for {
		msg, err := greetStreamResponse.Recv()
		// when stream ends, client will get EOF
		if err == io.EOF {
			// when stream ends, client will get EOF
			break
		}
		if err != nil {
			glog.Fatalf("Error while reading stream: %v\n", err)
		}
		glog.Infoln(msg)
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	glog.Infof("Starting client stream RPC client..")
	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Himanshu",
			},
		}, &greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Foo",
			},
		}, &greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Baz",
			},
		},
	}
	clientStream, err := c.LongGreet(context.Background())
	if err != nil {
		glog.Fatalln("Error connecting to server")
	}
	for _, request := range requests {
		clientStream.Send(request)
		glog.Infof("Sending message: %v\n", request)
		time.Sleep(1000 * time.Millisecond)
	}
	res, err := clientStream.CloseAndRecv()
	if err != nil {
		glog.Fatalf("Error while receiveing response: %v", err)
	}
	glog.Infoln(res)
}

func doBiDirectionalStreaming(c greetpb.GreetServiceClient) {
	glog.Infof("Starting bidirectional stream RPC client..")

	stream, connectErr := c.GreetEveryone(context.Background())
	if connectErr != nil {
		glog.Fatalln("Error connecting to server")
	}

	// create wait channel
	waitc := make(chan struct{})

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Himanshu",
			},
		}, &greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Foo",
			},
		}, &greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Baz",
			},
		},
	}

	// send message in separate go routine
	go func() {
		for _, req := range requests {
			glog.Infof("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// receive message in separate go routine
	go func() {
		for {
			res, recvErr := stream.Recv()
			if recvErr == io.EOF {
				// close wait channel after receiving all response
				break
			}
			if recvErr != nil {
				glog.Fatalf("Error while receiving: %v\n", recvErr)
				break
			}
			glog.Infof("Received: %v\n", res)
		}
		close(waitc)
	}()

	// block all goroutine
	<-waitc
}
