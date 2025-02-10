package main
import (
  "context"
  "fmt"
  "google.golang.org/grpc"
  "log"
  "net"
  pb "go-microservices/proto/user"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// In a real application, you would save the user to a database here
	user := &pb.CreateUserResponse{Id:    "12345",Name:  in.GetName(),Email: in.GetEmail(),} 
	return user, nil
}

func (s *server) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// In a real application, you would fetch the user from a database here
	user := &pb.GetUserResponse{Id:    in.GetId(),Name:  "John Doe",Email: "johndoe@example.com",}

	return user, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})
	fmt.Println("User Service is running on port 50051")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}