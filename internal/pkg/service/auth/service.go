package auth

import (
	"context"
	"errors"
	"time"

	pb "github.com/eduardojabes/CodeArena/proto/auth"
	pbU "github.com/eduardojabes/CodeArena/proto/user"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
	ErrEmptyKey     = errors.New("emptykey")
)

type JWTClaims struct {
	UserID uuid.UUID `json:"id"`
	jwt.StandardClaims
}

type AuthService struct {
	pb.UnimplementedAuthServer
	userService UserService
	JWTKey      []byte
}

type UserService interface {
	GetUserByUserNameAndPassword(ctx context.Context, in *pbU.GetUserByUserNameAndPasswordRequest) (*pbU.GetUserByUserNameAndPasswordResponse, error)
}

func (a *AuthService) generatejwt(userid uuid.UUID, secret []byte) (string, error) {
	if len(secret) == 0 {
		return "", ErrEmptyKey
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "code-arena",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		}, UserID: userid,
	})

	return token.SignedString(secret)
}

func (a *AuthService) GenerateToken(ctx context.Context, in *pb.GenerateTokenRequest) (*pb.GenerateTokenResponse, error) {

	user, err := a.userService.GetUserByUserNameAndPassword(ctx, &pbU.GetUserByUserNameAndPasswordRequest{
		Name:     in.GetUsername(),
		Password: in.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	userId, err := uuid.Parse(user.GetUserID())
	if err != nil {
		return nil, err
	}

	token, err := a.generatejwt(userId, a.JWTKey)
	if err != nil {
		return nil, err
	}
	return &pb.GenerateTokenResponse{Token: token}, nil
}

func (a *AuthService) VerifyToken(ctx context.Context, in *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	var claims JWTClaims

	token, err := jwt.ParseWithClaims(in.GetToken(), &claims, func(t *jwt.Token) (interface{}, error) {
		return a.JWTKey, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return &pb.VerifyTokenResponse{ValidToken: true}, nil
}

func (s *AuthService) Register(sr grpc.ServiceRegistrar) {
	pb.RegisterAuthServer(sr, s)
}

func NewAuthService(us UserService, jwtKey []byte) *AuthService {
	return &AuthService{
		userService: us,
		JWTKey:      jwtKey,
	}
}
