package auth

import (
	"context"

	pb "github.com/eduardojabes/CodeArena/proto/auth"
)

type AuthService struct {
	pb.UnimplementedAuthServer
}

func (a *AuthService) GenerateToken(context.Context, *pb.GenerateTokenRequest) (*pb.GenerateTokenResponse, error) {
	return nil, nil
}

func (a *AuthService) VerifyToken(context.Context, *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	return nil, nil
}

func NewAuthService() *AuthService {
	return &AuthService{}
}
