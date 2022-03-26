package auth

import (
	"context"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	s := NewAuthService()
	_, _ = s.GenerateToken(context.Background(), nil)

}

func TestVerifyToken(t *testing.T) {
	s := NewAuthService()
	_, _ = s.VerifyToken(context.Background(), nil)
}
