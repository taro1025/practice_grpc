package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"greet/calculator/calculatorpb"
	"log"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()
	client := calculatorpb.NewCalculatorServiceClient(cc)
	fmt.Println("Created client: %f", client)

	req := calculatorpb.SumRequest{Num1: 10, Num2: 3}
	result, _ := client.Sum(context.Background(), &req)
	log.Println("Addition Result: %v", result)
}
