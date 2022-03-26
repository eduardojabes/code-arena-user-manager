package util

import "golang.org/x/crypto/bcrypt"

type BCryptHasher struct {
	cost int
}

func (b *BCryptHasher) GenerateFromPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	return string(hash), err
}

func (b *BCryptHasher) CompareHashAndPassword(hashedPassword string, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil, err
}
func NewBCryptHasherWithCost(cost int) *BCryptHasher {
	return &BCryptHasher{cost: cost}
}

func NewBCryptHasher() *BCryptHasher {
	return NewBCryptHasherWithCost(12)
}
