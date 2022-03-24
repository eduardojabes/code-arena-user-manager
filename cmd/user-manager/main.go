package main

import (
	"fmt"
	"log"
	"net"

	"github.com/eduardojabes/CodeArena/internal/pkg/service/user"
	"google.golang.org/grpc"
)

var (
	port = int(50051)
)

func main() {
	// Set up a connection to the server.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	service := user.NewUserService()
	service.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
