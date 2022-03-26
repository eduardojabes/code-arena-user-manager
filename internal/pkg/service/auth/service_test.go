package auth

import (
	"context"
	"errors"
	"testing"

	pbU "github.com/eduardojabes/CodeArena/proto/user"
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
	userService := &userServiceMock{
		GetUserByUserNameAndPasswordMock: func(ctx context.Context, in *pbU.GetUserByUserNameAndPasswordRequest) (*pbU.GetUserByUserNameAndPasswordResponse, error) {
			return nil, nil
		},
	}

	s := NewAuthService(userService)
	_, _ = s.GenerateToken(context.Background(), nil)

}

func TestVerifyToken(t *testing.T) {
	userService := &userServiceMock{
		GetUserByUserNameAndPasswordMock: func(ctx context.Context, in *pbU.GetUserByUserNameAndPasswordRequest) (*pbU.GetUserByUserNameAndPasswordResponse, error) {
			return nil, nil
		},
	}

	s := NewAuthService(userService)
	_, _ = s.VerifyToken(context.Background(), nil)
}
