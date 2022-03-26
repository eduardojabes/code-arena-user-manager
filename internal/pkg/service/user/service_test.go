package user

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/eduardojabes/CodeArena/internal/pkg/entity"
	"github.com/eduardojabes/CodeArena/internal/pkg/util"
	"github.com/eduardojabes/CodeArena/internal/test"
	pb "github.com/eduardojabes/CodeArena/proto/user"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func TestHashPassword(t *testing.T) {
	t.Run("Creating User", func(t *testing.T) {

	})
}

func TestCreateUser(t *testing.T) {
	t.Run("Creating User effectively", func(t *testing.T) {
		input := &pb.CreateUserRequest{Name: "name", Password: "password"}

		repository := &test.MockUserRepository{
			SearchUserByUsernameMock: func(ctx context.Context, name string) (*entity.User, error) {
				return nil, nil
			},
			AddUserMock: func(ctx context.Context, user entity.User) error {
				return nil
			},
		}
		service := NewUserService(repository, util.NewBCryptHasher())

		got, err := service.CreateUser(context.Background(), input)

		if err != nil {
			t.Errorf("error creating user")
		}
		if got.UserID == uuid.Nil.String() {
			t.Errorf("got %v want an ID", got)
		}
	})

	t.Run("Error to create User in SearchUser", func(t *testing.T) {
		want := errors.New("error")
		input := &pb.CreateUserRequest{Name: "name", Password: "password"}

		repository := &test.MockUserRepository{
			SearchUserByUsernameMock: func(ctx context.Context, name string) (*entity.User, error) {
				return nil, want
			},
		}
		service := NewUserService(repository, util.NewBCryptHasher())

		_, got := service.CreateUser(context.Background(), input)

		if !errors.Is(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("Error to create User in SearchUser", func(t *testing.T) {
		want := errors.New("error")
		input := &pb.CreateUserRequest{Name: "name", Password: "password"}

		repository := &test.MockUserRepository{
			SearchUserByUsernameMock: func(ctx context.Context, name string) (*entity.User, error) {
				return nil, nil
			},
			AddUserMock: func(ctx context.Context, user entity.User) error {
				return want
			},
		}
		service := NewUserService(repository, util.NewBCryptHasher())

		_, got := service.CreateUser(context.Background(), input)

		if !errors.Is(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("Error to create User that exists", func(t *testing.T) {
		input := &pb.CreateUserRequest{Name: "name", Password: "password"}

		user := &entity.User{
			ID:       uuid.New(),
			Username: "user",
			Password: "password",
		}

		repository := &test.MockUserRepository{
			SearchUserByUsernameMock: func(ctx context.Context, name string) (*entity.User, error) {
				return user, nil
			},
		}
		service := NewUserService(repository, util.NewBCryptHasher())

		_, got := service.CreateUser(context.Background(), input)

		if !errors.Is(got, ErrUserAlreadyExists) {
			t.Errorf("got %v want %v", got, ErrUserAlreadyExists)
		}
	})

	t.Run("Error while hashing password", func(t *testing.T) {
		want := errors.New("error")
		input := &pb.CreateUserRequest{Name: "name", Password: "password"}

		hasher := &test.MockHasher{
			GenerateFromPasswordMock: func(password string) (string, error) {
				return "", want
			},
		}

		repository := &test.MockUserRepository{
			SearchUserByUsernameMock: func(ctx context.Context, name string) (*entity.User, error) {
				return nil, nil
			},
		}
		service := NewUserService(repository, hasher)

		_, got := service.CreateUser(context.Background(), input)

		if !errors.Is(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestGetUser(t *testing.T) {
	t.Run("error on search user", func(t *testing.T) {
		want := errors.New("error")
		repository := &test.MockUserRepository{
			SearchUserByUsernameMock: func(ctx context.Context, name string) (*entity.User, error) {
				return nil, want
			},
		}

		service := NewUserService(repository, util.NewBCryptHasher())

		_, got := service.GetUserByUserName(context.Background(), &pb.GetUserByUsernameRequest{})
		if !errors.Is(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("error username not exists", func(t *testing.T) {
		want := ErrUserNotExists
		repository := &test.MockUserRepository{
			SearchUserByUsernameMock: func(ctx context.Context, name string) (*entity.User, error) {
				return nil, nil
			},
		}

		service := NewUserService(repository, util.NewBCryptHasher())

		_, got := service.GetUserByUserName(context.Background(), &pb.GetUserByUsernameRequest{})
		if !errors.Is(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("getting existent user", func(t *testing.T) {
		user := &entity.User{
			ID:       uuid.New(),
			Username: "user",
			Password: "password",
		}

		want := pb.GetUserByUsernameResponse{
			UserID:   user.ID.String(),
			Name:     user.Username,
			Password: user.Password,
		}

		repository := &test.MockUserRepository{
			SearchUserByUsernameMock: func(ctx context.Context, name string) (*entity.User, error) {
				return user, nil
			},
		}

		service := NewUserService(repository, util.NewBCryptHasher())

		got, err := service.GetUserByUserName(context.Background(), &pb.GetUserByUsernameRequest{Name: "user"})

		if err != nil {
			t.Errorf("got %v want nil", got)
		}

		if !reflect.DeepEqual(&want, got) {
			t.Errorf("got %v want %v", got, &want)
		}
	})
}

func TestRegisterUserService(t *testing.T) {
	s := grpc.NewServer()
	service := NewUserService(&test.MockUserRepository{}, util.NewBCryptHasher())

	service.Register(s)
}

func TestGetUserByUserNameAndPassword(t *testing.T) {
	t.Run("error getting user", func(t *testing.T) {
		want := errors.New("error")
		repository := &test.MockUserRepository{
			SearchUserByUsernameMock: func(ctx context.Context, name string) (*entity.User, error) {
				return nil, want
			},
		}

		service := NewUserService(repository, util.NewBCryptHasher())
		_, got := service.GetUserByUserNameAndPassword(context.Background(), &pb.GetUserByUserNameAndPasswordRequest{})
		if !errors.Is(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("error hashing password", func(t *testing.T) {
		user := &entity.User{
			ID:       uuid.New(),
			Username: "user",
			Password: "password",
		}

		want := errors.New("error")

		repository := &test.MockUserRepository{
			SearchUserByUsernameMock: func(ctx context.Context, name string) (*entity.User, error) {
				return user, nil
			},
		}

		hasher := &test.MockHasher{
			CompareHashAndPasswordMock: func(hasher string, password string) (bool, error) {
				return false, want
			},
		}

		service := NewUserService(repository, hasher)

		_, got := service.GetUserByUserNameAndPassword(context.Background(), &pb.GetUserByUserNameAndPasswordRequest{
			Name:     user.Username,
			Password: user.Password,
		})

		if !errors.Is(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("error not exists user", func(t *testing.T) {
		user := &entity.User{
			ID:       uuid.New(),
			Username: "user",
			Password: "password",
		}

		want := ErrInvalidPassword

		repository := &test.MockUserRepository{
			SearchUserByUsernameMock: func(ctx context.Context, name string) (*entity.User, error) {
				return nil, nil
			},
		}

		hasher := &test.MockHasher{
			CompareHashAndPasswordMock: func(hasher string, password string) (bool, error) {
				return false, want
			},
		}

		service := NewUserService(repository, hasher)

		_, got := service.GetUserByUserNameAndPassword(context.Background(), &pb.GetUserByUserNameAndPasswordRequest{
			Name:     user.Username,
			Password: user.Password,
		})

		if !errors.Is(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("error incorrect password", func(t *testing.T) {
		user := &entity.User{
			ID:       uuid.New(),
			Username: "user",
			Password: "password",
		}

		want := ErrInvalidPassword

		repository := &test.MockUserRepository{
			SearchUserByUsernameMock: func(ctx context.Context, name string) (*entity.User, error) {
				return user, nil
			},
		}

		hasher := &test.MockHasher{
			CompareHashAndPasswordMock: func(hasher string, password string) (bool, error) {
				return false, nil
			},
		}

		service := NewUserService(repository, hasher)

		_, got := service.GetUserByUserNameAndPassword(context.Background(), &pb.GetUserByUserNameAndPasswordRequest{
			Name:     user.Username,
			Password: user.Password,
		})

		if !errors.Is(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("correct user and password", func(t *testing.T) {
		user := &entity.User{
			ID:       uuid.New(),
			Username: "user",
			Password: "password",
		}

		want := &pb.GetUserByUserNameAndPasswordResponse{UserID: user.ID.String()}

		repository := &test.MockUserRepository{
			SearchUserByUsernameMock: func(ctx context.Context, name string) (*entity.User, error) {
				return user, nil
			},
		}

		hasher := &test.MockHasher{
			CompareHashAndPasswordMock: func(hasher string, password string) (bool, error) {
				return true, nil
			},
		}

		service := NewUserService(repository, hasher)

		got, _ := service.GetUserByUserNameAndPassword(context.Background(), &pb.GetUserByUserNameAndPasswordRequest{
			Name:     user.Username,
			Password: user.Password,
		})

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

}
