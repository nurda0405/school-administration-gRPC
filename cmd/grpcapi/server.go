package main

import (
	"context"
	"fmt"
	"grpcapi/internals/api/handlers"
	"grpcapi/internals/repositories/mongodb"
	pb "grpcapi/proto/gen"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	mongodb.CreateMongoClient(context.Background())
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	s := grpc.NewServer()
	pb.RegisterExecsServiceServer(s, &handlers.Server{})
	pb.RegisterStudentsServiceServer(s, &handlers.Server{})
	pb.RegisterTeachersServiceServer(s, &handlers.Server{})

	reflection.Register(s)

	port := os.Getenv("SERVER_PORT")

	fmt.Println("gRPC server is running on port", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Could not listen to the port:", err)
	}
	err = s.Serve(lis)
	if err != nil {
		log.Fatal("Error serving the listener:", err)
	}
}
