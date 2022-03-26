package auth

import (
	"context"
	"time"

	pb "github.com/eduardojabes/CodeArena/proto/auth"
	pbU "github.com/eduardojabes/CodeArena/proto/user"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var JWTKey = []byte("code-arena-key")

type JWTClaims struct {
	UserID uuid.UUID `json:"id"`
	jwt.StandardClaims
}

type AuthService struct {
	pb.UnimplementedAuthServer
	userService UserService
}

type UserService interface {
	GetUserByUserNameAndPassword(ctx context.Context, in *pbU.GetUserByUserNameAndPasswordRequest) (*pbU.GetUserByUserNameAndPasswordResponse, error)
}

func (a *AuthService) GenerateToken(ctx context.Context, in *pb.GenerateTokenRequest) (*pb.GenerateTokenResponse, error) {

	user, err := a.userService.GetUserByUserNameAndPassword(ctx, &pbU.GetUserByUserNameAndPasswordRequest{
		Name:     in.GetUsername(),
		Password: in.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	userid, err := uuid.Parse(user.GetUserID())
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, JWTClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "code-arena",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		}, UserID: userid,
	})

	stoken, err := token.SignedString(JWTKey)
	if err != nil {
		return nil, err
	}
	return &pb.GenerateTokenResponse{Token: stoken}, nil
}

func (a *AuthService) VerifyToken(context.Context, *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	return nil, nil
}

func NewAuthService(us UserService) *AuthService {
	return &AuthService{userService: us}
}
