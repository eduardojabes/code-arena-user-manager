package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/eduardojabes/CodeArena/internal/pkg/entity"
	"github.com/eduardojabes/CodeArena/proto/auth"
	pbU "github.com/eduardojabes/CodeArena/proto/user"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type userServiceMock struct {
	GetUserByUserNameAndPasswordMock func(ctx context.Context, in *pbU.GetUserByUserNameAndPasswordRequest) (*pbU.GetUserByUserNameAndPasswordResponse, error)
}

func (usm *userServiceMock) GetUserByUserNameAndPassword(ctx context.Context, in *pbU.GetUserByUserNameAndPasswordRequest) (*pbU.GetUserByUserNameAndPasswordResponse, error) {
	if usm.GetUserByUserNameAndPasswordMock != nil {
		return usm.GetUserByUserNameAndPasswordMock(ctx, in)
	}

	return nil, errors.New("GetUserByUserNameAndPasswordMock must be set")
}

func TestGenerateToken(t *testing.T) {
	t.Run("unable to get user by password", func(t *testing.T) {
		want := errors.New("error")
		userService := &userServiceMock{
			GetUserByUserNameAndPasswordMock: func(ctx context.Context, in *pbU.GetUserByUserNameAndPasswordRequest) (*pbU.GetUserByUserNameAndPasswordResponse, error) {
				return nil, want
			},
		}

		s := NewAuthService(userService, []byte("code-arena-key"))
		_, got := s.GenerateToken(context.Background(), nil)

		if !errors.Is(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("error in uuid parsing", func(t *testing.T) {
		userService := &userServiceMock{
			GetUserByUserNameAndPasswordMock: func(ctx context.Context, in *pbU.GetUserByUserNameAndPasswordRequest) (*pbU.GetUserByUserNameAndPasswordResponse, error) {
				return &pbU.GetUserByUserNameAndPasswordResponse{UserID: "ABC"}, nil
			},
		}

		s := NewAuthService(userService, []byte("code-arena-key"))
		_, got := s.GenerateToken(context.Background(), nil)

		if got == nil {
			t.Errorf("got %v want nil", got)
		}
	})
	t.Run("error to gerenate token", func(t *testing.T) {
		user := &entity.User{
			ID:       uuid.New(),
			Username: "user",
			Password: "password",
		}

		userService := &userServiceMock{
			GetUserByUserNameAndPasswordMock: func(ctx context.Context, in *pbU.GetUserByUserNameAndPasswordRequest) (*pbU.GetUserByUserNameAndPasswordResponse, error) {
				return &pbU.GetUserByUserNameAndPasswordResponse{UserID: user.ID.String()}, nil
			},
		}

		s := NewAuthService(userService, nil)
		token, got := s.GenerateToken(context.Background(), &auth.GenerateTokenRequest{Username: user.Username, Password: user.Password})

		if got == nil {
			t.Errorf("got %v want error", got)
		}
		if token != nil {
			t.Errorf("got %v want nil", got)
		}
	})
	t.Run("token correctly generated", func(t *testing.T) {
		user := &entity.User{
			ID:       uuid.New(),
			Username: "user",
			Password: "password",
		}

		userService := &userServiceMock{
			GetUserByUserNameAndPasswordMock: func(ctx context.Context, in *pbU.GetUserByUserNameAndPasswordRequest) (*pbU.GetUserByUserNameAndPasswordResponse, error) {
				return &pbU.GetUserByUserNameAndPasswordResponse{UserID: user.ID.String()}, nil
			},
		}

		s := NewAuthService(userService, []byte("code-arena-key"))
		_, got := s.GenerateToken(context.Background(), &auth.GenerateTokenRequest{Username: user.Username, Password: user.Password})

		if got != nil {
			t.Errorf("got %v want nil", got)
		}
	})
}

func TestVerifyToken(t *testing.T) {
	t.Run("Received wrong token", func(t *testing.T) {
		want := ErrInvalidToken
		userService := &userServiceMock{}

		s := NewAuthService(userService, []byte("code-arena-key"))

		_, got := s.VerifyToken(context.Background(), &auth.VerifyTokenRequest{Token: ""})
		if !errors.Is(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("Received invalid token", func(t *testing.T) {
		want := ErrInvalidToken
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
			StandardClaims: jwt.StandardClaims{
				Issuer:    "code-arena",
				IssuedAt:  time.Now().Add(time.Hour * -1).Unix(),
				ExpiresAt: time.Now().Add(time.Hour * -1).Unix(),
			}, UserID: uuid.New(),
		})

		tokenString, _ := token.SignedString([]byte("code-arena-key"))

		userService := &userServiceMock{}
		s := NewAuthService(userService, []byte("code-arena-ke"))
		_, got := s.VerifyToken(context.Background(), &auth.VerifyTokenRequest{Token: tokenString})
		if !errors.Is(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("verified token", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
			StandardClaims: jwt.StandardClaims{
				Issuer:    "code-arena",
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			}, UserID: uuid.New(),
		})

		tokenString, _ := token.SignedString([]byte("code-arena-key"))

		userService := &userServiceMock{}
		s := NewAuthService(userService, []byte("code-arena-key"))
		got, err := s.VerifyToken(context.Background(), &auth.VerifyTokenRequest{Token: tokenString})

		if err != nil {
			t.Errorf("got %v want nil", err)
		}
		if !got.GetValidToken() {
			t.Errorf("got %v want true", got.GetValidToken())
		}
	})
}

func TestRegisterUserService(t *testing.T) {
	s := grpc.NewServer()
	service := NewAuthService(&userServiceMock{}, []byte("code-arena-key"))

	service.Register(s)
}
