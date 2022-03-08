package main

import (
	"context"
	"google.golang.org/grpc"
	"greet/calculator/calculatorpb"
	"log"
	"net"
)

type server struct {
	calculatorpb.CalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	result := &calculatorpb.SumResponse{Result: req.GetNum1() + req.GetNum2()}
	return result, nil
}

func main(){
	log.Printf("Go Server Started...")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
	log.Fatalf("Failed to Listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
	log.Fatalf("Failed to serve: %v", err)
	}
}
