package util

import (
	"testing"
)

func TestHasher(t *testing.T) {
	t.Run("Hashing Password", func(t *testing.T) {
		b := NewBCryptHasher()
		password := "password"
		hashedPassword, err := b.GenerateFromPassword(password)
		if err != nil {
			t.Errorf("Hashing Error: %v", err)
		}
		ok, err := b.CompareHashAndPassword(hashedPassword, password)
		if err != nil {
			t.Errorf("error comparing password hash: %v", err)
		}
		if !ok {
			t.Error("The passwords are different")
		}
	})
}
