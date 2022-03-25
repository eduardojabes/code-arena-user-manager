package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"

	userRepository "github.com/eduardojabes/CodeArena/internal/pkg/repository/user"
	userService "github.com/eduardojabes/CodeArena/internal/pkg/service/user"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
)

var (
	port        = int(50051)
	DatabaseUrl = "postgres://postgres:postgres@localhost:5432/user-manager"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	// Set up connection to DB
	conn, err := pgx.Connect(context.Background(), DatabaseUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Set up a connection to the server.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	repository := userRepository.NewPostgreUserRepository(conn)
	service := userService.NewUserService(repository)
	service.Register(s)

	go func() {
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx := context.Background()

	s.GracefulStop()
	err = conn.Close(shutdownCtx)
	if err != nil {
		panic(err)
	}
}
