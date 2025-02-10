package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pbOrder "grpc-orders/grpc-orders/proto/order"
	pbUser "grpc-orders/grpc-orders/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to the user service
	userConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to user service: %v", err)
	}
	defer userConn.Close()

	userClient := pbUser.NewUserServiceClient(userConn)

	// Create a new user
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createUserResp, err := userClient.CreateUser(ctx, &pbUser.CreateUserRequest{
		Name:  "Alice",
		Email: "alice@example.com",
	})
	if err != nil {
		log.Fatalf("could not create user: %v", err)
	}
	fmt.Printf("Created User: %v\n", createUserResp)

	// Connect to the order service
	orderConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to order service: %v", err)
	}
	defer orderConn.Close()

	orderClient := pbOrder.NewOrderServiceClient(orderConn)

	// Create a new order for the user
	createOrderResp, err := orderClient.CreateOrder(ctx, &pbOrder.CreateOrderRequest{
		UserId:  createUserResp.GetId(),
		Product: "Laptop",
		Price:   999.99,
	})
	if err != nil {
		log.Fatalf("could not create order: %v", err)
	}
	fmt.Printf("Created Order: %v\n", createOrderResp)
}
