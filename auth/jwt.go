package auth

import (
	"errors"
	"time"

	"github.com/WagaoCarvalho/backendgoswagger/config"
	"github.com/WagaoCarvalho/backendgoswagger/internal/models"
	"github.com/WagaoCarvalho/backendgoswagger/internal/repositories"
	"github.com/WagaoCarvalho/backendgoswagger/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

// Estrutura de resposta para autenticação
type Auth struct {
	IsValid bool   `json:"is_valid"`
	Token   string `json:"token,omitempty"`
}

// Erro de senha inválida
var ErrInvalidPassword = errors.New("senha inválida")

// Carrega configurações
var configs = config.LoadConfig()

// SignIn verifica o usuário e retorna um token JWT se válido
func SignIn(user models.User) (Auth, error) {
	inputPassword := user.Password

	// Buscar usuário no banco de dados
	user, err := repositories.GetUserByEmail(user.Email)
	if err != nil {
		return Auth{IsValid: false}, err
	}

	// Verificar senha
	hasher := utils.BcryptHasher{}
	if err := hasher.CheckPasswordErr(inputPassword, user.Password); err != nil {
		return Auth{IsValid: false}, ErrInvalidPassword
	}

	// Gerar token JWT
	token, err := GenerateJWT(user)
	if err != nil {
		return Auth{IsValid: false}, err
	}

	return Auth{IsValid: true, Token: token}, nil
}

// GenerateJWT cria um token JWT assinado
func GenerateJWT(user models.User) (string, error) {
	// Criar claims do token
	claims := jwt.MapClaims{
		"authorized": true,
		"user_id":    user.UID,
		"exp":        time.Now().Add(1 * time.Hour).Unix(),
	}

	// Criar token com método de assinatura HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Assinar o token com a chave secreta do config
	return token.SignedString([]byte(configs.Jwt.SecretKey))
}
