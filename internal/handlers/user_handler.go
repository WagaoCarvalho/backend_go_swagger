package handlers

import (
	"fmt"
	"net/http"

	"github.com/WagaoCarvalho/backendgoswagger/internal/models"
	"github.com/WagaoCarvalho/backendgoswagger/internal/repositories"
	"github.com/WagaoCarvalho/backendgoswagger/internal/utils"
)

// GetUsersHandler lida com a requisição de busca de todos os usuários
func GetUsers(w http.ResponseWriter, r *http.Request) {
	// Chama a função GetUsers no repositório
	users, err := repositories.GetUsers()
	if err != nil {
		// Usando ErrorResponse para retornar um erro em formato JSON
		utils.ErrorResponse(w, fmt.Errorf("erro ao buscar usuários: %w", err), http.StatusInternalServerError)
		return
	}

	// Converte os usuários para JSON e envia na resposta com status
	response := utils.DefaultResponse{
		Data:   users,
		Status: http.StatusOK,
	}
	utils.ToJson(w, response)
}

// PostUser lida com a requisição de criação de um novo usuário
func PostUser(w http.ResponseWriter, r *http.Request) {
	// Verifica o método HTTP
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	// Cria uma instância de User para armazenar os dados
	var user models.User

	// Decodifica o corpo da requisição JSON para a struct User
	if err := utils.FromJson(r.Body, &user); err != nil {
		utils.ErrorResponse(w, fmt.Errorf("erro ao decodificar JSON: %w", err), http.StatusBadRequest)
		return
	}

	// Validação simples: garantir que os campos obrigatórios não estão vazios
	if user.Nickname == "" || user.Email == "" || user.Password == "" {
		utils.ErrorResponse(w, fmt.Errorf("nickname, email e password são obrigatórios"), http.StatusBadRequest)
		return
	}

	// Chama a função NewUser do repositório para inserir o novo usuário
	success, err := repositories.NewUser(user)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("erro ao criar usuário: %w", err), http.StatusInternalServerError)
		return
	}

	// Se a criação for bem-sucedida, retorna um status 201 Created
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
