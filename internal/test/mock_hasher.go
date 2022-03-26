package test

import "errors"

type MockHasher struct {
	GenerateFromPasswordMock   func(password string) (string, error)
	CompareHashAndPasswordMock func(hashedPassword string, password string) (bool, error)
}

func (m *MockHasher) GenerateFromPassword(password string) (string, error) {
	if m.GenerateFromPasswordMock != nil {
		return m.GenerateFromPasswordMock(password)
	}
	return "", errors.New("GenerateFromPasswordMock must be set")
}

func (m *MockHasher) CompareHashAndPassword(hashedPassword string, password string) (bool, error) {
	if m.CompareHashAndPasswordMock != nil {
		return m.CompareHashAndPasswordMock(hashedPassword, password)
	}
	return false, errors.New("CompareHashAndMock must be set")
}
