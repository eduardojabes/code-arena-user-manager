package user

import (
	"context"
	"log"

	pb "github.com/eduardojabes/CodeArena/proto/user"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type UserRepository interface {
	GetUser(ctx context.Context, ID uuid.UUID) (*User, error)
	AddUser(ctx context.Context, user User) (*User, error)
	SearchUserByUsername(ctx context.Context, name string) (*User, error)
}

type UserService struct {
	pb.UnimplementedUserServer
	repository UserRepository
}

type User struct {
	Username string
	Password string
	ID       uuid.UUID
}

func (s *UserService) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Printf("Received: %v", in.GetName())
	log.Printf("Received: %v", in.GetPassword())

	//Preparar função para gerenciar a autenticação

	return &pb.CreateUserResponse{UserID: uuid.NewString()}, nil
}

func (s *UserService) Register(sr grpc.ServiceRegistrar) {
	pb.RegisterUserServer(sr, s)
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{repository: repository}
}
