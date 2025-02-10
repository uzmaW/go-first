package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "grpc-orders/grpc-orders/proto/order"
)

type server struct {
	pb.UnimplementedOrderServiceServer
}

func (s *server) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	order := &pb.CreateOrderResponse{
		OrderId: in.GetUserId() + "_12345",
		UserId:  in.GetUserId(),
		Product: in.GetProduct(),
		Price:   in.GetPrice()}
	return order, nil
}

func (s *server) GetOrder(ctx context.Context, in *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order := &pb.GetOrderResponse{
		OrderId: in.GetOrderId(),
		UserId:  "12345",
		Product: "Laptop",
		Price:   999.99,
	}
	return order, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, &server{})
	fmt.Println("Order Service is running on port 50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
