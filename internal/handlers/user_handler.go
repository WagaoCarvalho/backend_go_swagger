package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/WagaoCarvalho/backendgoswagger/internal/models"
	"github.com/WagaoCarvalho/backendgoswagger/internal/repositories"
	"github.com/WagaoCarvalho/backendgoswagger/internal/utils"
	"github.com/WagaoCarvalho/backendgoswagger/internal/validations"
	"github.com/gorilla/mux"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := repositories.GetUsers()
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("erro ao buscar usuários: %w", err), http.StatusInternalServerError)
		return
	}

	response := utils.DefaultResponse{
		Data:   users,
		Status: http.StatusOK,
	}
	utils.ToJson(w, response)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	uidParam := vars["id"]

	uid, err := strconv.ParseInt(uidParam, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	// Recupera o usuário com o ID
	user, err := repositories.GetUserById(uid)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	// Prepara a resposta
	response := utils.DefaultResponse{
		Data:   user,
		Status: http.StatusOK,
	}
	utils.ToJson(w, response)
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}
	var user models.User

	if err := utils.FromJson(r.Body, &user); err != nil {
		utils.ErrorResponse(w, fmt.Errorf("erro ao decodificar JSON: %w", err), http.StatusBadRequest)
		return
	}

	user, err := validations.ValidateNewUser(user)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusUnprocessableEntity)
		return
	}

	success, err := repositories.NewUser(user)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("erro ao criar usuário: %w", err), http.StatusInternalServerError)
		return
	}

	if success {
		response := utils.DefaultResponse{
			Data:   user,
			Status: http.StatusCreated,
		}
		utils.ToJson(w, response)
	} else {
		utils.ErrorResponse(w, fmt.Errorf("erro desconhecido ao criar usuário"), http.StatusInternalServerError)
	}
}

func PutUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	uidParam := vars["id"]

	uid, err := strconv.ParseInt(uidParam, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	var user models.User
	if err := utils.FromJson(r.Body, &user); err != nil {
		utils.ErrorResponse(w, fmt.Errorf("JSON inválido"), http.StatusBadRequest)
		return
	}

	user.UID = uid
	if err := repositories.UpdateUser(&user); err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.ToJson(w, user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	uidParam := vars["id"]

	uid, err := strconv.ParseInt(uidParam, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("ID inválido"), http.StatusBadRequest)
		return
	}

	if err := repositories.DeleteUser(uid); err != nil {
		if err.Error() == "usuário não encontrado" {
			utils.ErrorResponse(w, err, http.StatusNotFound)
		} else {
			utils.ErrorResponse(w, err, http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
