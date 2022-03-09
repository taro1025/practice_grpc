package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"greet/greet/greetpb"
	"io"
	"log"
	"time"
)

func main() {
	fmt.Println("Hello I'm a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	client := greetpb.NewGreetServiceClient(cc)
	fmt.Println("Created client: %f", client)

	//1:1 (server:client)
	doUnary(client)

	// n:1
	doServerStreaming(client)

	// 1:n
	doClientStreaming(client)

	// n:n
	doBiDiStreaming(client)
}

func doUnary(client greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Taro",
			LastName:  "Yaha",
		},
	}
	res, err := client.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf(" Response from Greet: %v", res.Result)
}

func doServerStreaming(client greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "taro",
			LastName:  "yaha",
		},
	}
	resStream, err := client.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}
	for {

		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

func doClientStreaming(client greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	reqests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Taro",
				LastName:  "Yaha",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "aaa",
				LastName:  "aaa",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "bbb",
				LastName:  "bbb",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "cccc",
				LastName:  "cccc",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "ddd",
				LastName:  "ddd",
			},
		},
	}

	stream, err := client.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while request %v: ", err)
	}

	for _, req := range reqests {
		stream.Send(req)
		log.Printf("send %v", req.GetGreeting().LastName)
		time.Sleep(1000 * time.Millisecond)
	}
	// RecieveとCloseはセット
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while recieving response %v: ", err)
	}
	log.Printf("response: %v", res)
}

func doBiDiStreaming(client greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC...")

	stream, err := client.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while createing stream: %v", err)
	}

	reqests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Taro",
				LastName:  "Yaha",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "aaa",
				LastName:  "aaa",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "bbb",
				LastName:  "bbb",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "cccc",
				LastName:  "cccc",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "ddd",
				LastName:  "ddd",
			},
		},
	}

	waitc := make(chan struct{})

	// we send a bunch of message to the client
	go func() {
		for _, req := range reqests {
			fmt.Printf("Sengding message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000* time.Millisecond)
		}
		stream.CloseSend()
	}()

	// we receive a bunch of message from the client
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}