package main

import (
	"log"
	"net"

	"github.com/XuananLe/Golang-For-Devops/Grpc/pb"
	"google.golang.org/grpc/reflection" // Import reflection package
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedSimpleServiceServer
}
func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer();
	pb.RegisterSimpleServiceServer(s, &Server{})
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}