package auth

import (
	"errors"

	"github.com/WagaoCarvalho/backendgoswagger/internal/models"
	"github.com/WagaoCarvalho/backendgoswagger/internal/repositories"
	"github.com/WagaoCarvalho/backendgoswagger/internal/utils"
)

var (
	ErrInvalidPassword = errors.New("Senha inválida")
)

func SignIn(user models.User) (models.User, error) {
	inputPassword := user.Password

	user, err := repositories.GetUserByEmail(user.Email)
	if err != nil {
		return user, err
	}

	hasher := utils.BcryptHasher{}

	err = hasher.CheckPasswordErr(inputPassword, user.Password)
	if err != nil {
		return models.User{}, ErrInvalidPassword
	}

	return user, nil
}
