package auth

import (
	"context"
	"errors"
	"time"

	pb "github.com/eduardojabes/CodeArena/proto/auth"
	pbU "github.com/eduardojabes/CodeArena/proto/user"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
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

func (a *AuthService) VerifyToken(ctx context.Context, in *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	var claims JWTClaims

	token, err := jwt.ParseWithClaims(in.GetToken(), claims, func(t *jwt.Token) (interface{}, error) {
		return JWTKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}
	return &pb.VerifyTokenResponse{ValidToken: true}, nil
}

func NewAuthService(us UserService) *AuthService {
	return &AuthService{userService: us}
}
