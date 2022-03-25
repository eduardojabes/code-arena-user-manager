package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	pb "github.com/eduardojabes/CodeArena/proto/user"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ErrCreateUserMessage        = "error creation failed"
	ErrUserAlreadyExistsMessage = "user already exists"
)

type UserRepository interface {
	GetUser(ctx context.Context, ID uuid.UUID) (*User, error)
	AddUser(ctx context.Context, user User) error
	SearchUserByUsername(ctx context.Context, name string) (*User, error)
}

type UserService struct {
	pb.UnimplementedUserServer
	repository UserRepository
}

type User struct {
	ID       uuid.UUID
	Username string
	Password string
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	return reflect.DeepEqual([]byte(hash), []byte(password))
}

func (s *UserService) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user, err := s.repository.SearchUserByUsername(ctx, in.GetName())
	if err != nil {
		err = fmt.Errorf("error while creating user: %v", err)
		log.Println(err)

		return nil, status.Error(codes.Internal, ErrCreateUserMessage)
	}

	if user != nil {
		return nil, status.Error(codes.AlreadyExists, ErrUserAlreadyExistsMessage)
	}

	password, err := HashPassword(in.GetPassword())

	if err != nil {
		err = fmt.Errorf("error while encrypting the password: %v", err)
		log.Println(err)

		return nil, status.Error(codes.Internal, ErrCreateUserMessage)
	}

	newUser := User{
		ID:       uuid.New(),
		Username: in.GetName(),
		Password: password,
	}

	err = s.repository.AddUser(ctx, newUser)

	if err != nil {
		err = fmt.Errorf("error while attempt to access user DB: %v", err)
		log.Println(err)

		return nil, status.Error(codes.Internal, ErrCreateUserMessage)
	}

	return &pb.CreateUserResponse{UserID: uuid.NewString()}, nil
}

func (s *UserService) LogUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	password, err := HashPassword(in.GetPassword())
	if err != nil {
		err = fmt.Errorf("error hashing password: %v", err)
		log.Println(err)

		return nil, status.Error(codes.Internal, ErrCreateUserMessage)
	}
	user, err := s.repository.SearchUserByUsername(ctx, in.GetName())
	if err != nil {
		err = fmt.Errorf("error getting user data from server: %v", err)
		log.Println(err)

		return nil, status.Error(codes.Internal, ErrCreateUserMessage)
	}

	if CheckPasswordHash(password, user.Password) {
		return &pb.CreateUserResponse{UserID: uuid.NewString()}, nil
	}

	err = errors.New("wrong password")
	return nil, err
}

func (s *UserService) Register(sr grpc.ServiceRegistrar) {
	pb.RegisterUserServer(sr, s)
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{repository: repository}
}
