package main
import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pbUser "github.com/uzmaW/go-first/grpc/proto/user"
	pbOrder "github.com/uzmaW/go-first/grpc/proto/order"
	"log"
)
func main() {
	// Connect to the user service
	userConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials())
	if err != nil {
		log.Fatalf("did not connect to user service: %v", err)}
		defer userConn.Close()
		userClient := pbUser.NewUserServiceClient(userConn)
		// Create a new user
		createUserResp, err := userClient.CreateUser(context.Background(), &pbUser.CreateUserRequest{
			Name: "Alice", Email: "alice@example.com",},
			)
			if err != nil {
				log.Fatalf("could not create user: %v", err)
				
			}
		fmt.Printf("Created User: %v\n", createUserResp)
// Connect to the order serviceorderConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())if err != nil {log.Fatalf("did not connect to order service: %v", err)}defer orderConn.Close()orderClient := pbOrder.NewOrderServiceClient(orderConn)
// Create a new order for the user
createOrderResp, err := orderClient.CreateOrder(context.Background(), &pbOrder.CreateOrderRequest{
	UserId: createUserResp.GetId(),
	Product: "Laptop",Price:   999.99,
	})
	if err != nil {
		log.Fatalf("could not create order: %v", err)
	}
	fmt.Printf("Created Order: %v\n", createOrderResp)
}