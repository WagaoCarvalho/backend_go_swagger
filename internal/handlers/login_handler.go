package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/WagaoCarvalho/backendgoswagger/auth"
	"github.com/WagaoCarvalho/backendgoswagger/internal/models"
	"github.com/WagaoCarvalho/backendgoswagger/internal/utils"
)

func Login(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var user models.User
	err := json.Unmarshal(body, &user)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	authResponse, err := auth.SignIn(user) // Captura corretamente a estrutura `Auth`
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	utils.ToJson(w, authResponse) // Retorna `auth.Auth`, que inclui o token JWT
}
