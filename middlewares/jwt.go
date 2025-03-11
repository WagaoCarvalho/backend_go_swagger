package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/WagaoCarvalho/backendgoswagger/config"
	"github.com/golang-jwt/jwt/v5"
)

// Middleware para verificar autenticação via Bearer Token
func IsAuthByBearerToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Token ausente", http.StatusUnauthorized)
			return
		}

		// Verifica se o token começa com "Bearer "
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Formato de token inválido", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		secretKey := []byte(config.LoadConfig().Jwt.SecretKey)

		// Valida o token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		// Adiciona os claims ao contexto para uso posterior
		ctx := context.WithValue(r.Context(), "user", token.Claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
