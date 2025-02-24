package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Interface para hashing de senha (permite mocks nos testes)
type PasswordHasher interface {
	HashPassword(password string) (string, error)
	CheckPassword(password, hash string) bool
}

// Implementação real usando bcrypt
type BcryptHasher struct{}

// Gera um hash da senha
func (b BcryptHasher) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Falha ao gerar o hash da senha")
	}
	return string(hash), nil
}

// Compara uma senha com um hash armazenado
func (b BcryptHasher) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
