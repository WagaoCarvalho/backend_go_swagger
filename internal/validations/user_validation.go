package validations

import (
	"errors"

	"github.com/WagaoCarvalho/backendgoswagger/internal/models"
)

var (
	ErrEmptyFields  = errors.New("Um ou mais campos estão vazios")
	ErrInvalidEmail = errors.New("Email inválido")
)

func ValidateNewUser(user models.User) (models.User, error) {
	if IsStringEmpty(user.Nickname) || IsStringEmpty(user.Email) || IsStringEmpty(user.Password) {
		return models.User{}, ErrEmptyFields
	}
	if !IsValidEmail(user.Email) {
		return models.User{}, ErrInvalidEmail
	}
	return user, nil
}
