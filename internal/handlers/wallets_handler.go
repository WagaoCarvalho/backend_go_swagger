package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/WagaoCarvalho/backendgoswagger/internal/models"
	"github.com/WagaoCarvalho/backendgoswagger/internal/repositories"
	"github.com/WagaoCarvalho/backendgoswagger/internal/utils"
)

func GetWallets(w http.ResponseWriter, r *http.Request) {
	wallets, err := repositories.GetWallets()
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	utils.ToJson(w, wallets)
}

func GetWallet(w http.ResponseWriter, r *http.Request) {
	// Verifica se o corpo da requisição está vazio antes de tentar decodificá-lo
	if r.Body == nil {
		utils.ErrorResponse(w, fmt.Errorf("corpo da requisição vazio"), http.StatusBadRequest)
		return
	}

	// Decodifica o JSON do corpo da requisição
	var requestBody map[string]string
	if err := utils.FromJson(r.Body, &requestBody); err != nil {
		utils.ErrorResponse(w, fmt.Errorf("JSON inválido"), http.StatusBadRequest)
		return
	}

	// Verifica se a chave pública foi fornecida no corpo da requisição
	publicKey, exists := requestBody["public_key"]
	if !exists || publicKey == "" {
		utils.ErrorResponse(w, fmt.Errorf("chave pública não fornecida"), http.StatusBadRequest)
		return
	}

	// Busca a carteira usando a chave pública
	wallet, err := repositories.GetWalletByPublicKey(publicKey)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	utils.ToJson(w, wallet)
}

func PutWallet(w http.ResponseWriter, r *http.Request) {
	// Valida o método HTTP
	if r.Method != http.MethodPut {
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	// Verifica se o corpo da requisição está vazio antes de tentar decodificá-lo
	if r.Body == nil {
		utils.ErrorResponse(w, fmt.Errorf("corpo da requisição vazio"), http.StatusBadRequest)
		return
	}

	// Decodifica o JSON do corpo da requisição
	var wallet models.Wallet
	if err := utils.FromJson(r.Body, &wallet); err != nil {
		utils.ErrorResponse(w, fmt.Errorf("JSON inválido"), http.StatusBadRequest)
		return
	}

	// Verifica se a chave pública foi fornecida no corpo da requisição
	if wallet.PublicKey == "" {
		utils.ErrorResponse(w, fmt.Errorf("chave pública não fornecida"), http.StatusBadRequest)
		return
	}

	// Captura o parâmetro de adição no query string (?add=true)
	add, err := strconv.ParseBool(r.URL.Query().Get("add"))
	if err != nil {
		add = false // Define como false em caso de erro na conversão
	}

	// Variáveis para armazenar o resultado da operação no repositório
	var userID, rowsAffected int64
	var errRepo error

	// Seleciona a função de atualização de saldo ou adição de saldo
	if add {
		userID, rowsAffected, errRepo = repositories.AddBalanceWallet(wallet)
	} else {
		userID, rowsAffected, errRepo = repositories.UpdateWallet(wallet)
	}

	// Tratamento de erro do repositório
	if errRepo != nil {
		statusCode := http.StatusInternalServerError
		if errRepo.Error() == "carteira não encontrada" {
			statusCode = http.StatusNotFound
		}
		utils.ErrorResponse(w, errRepo, statusCode)
		return
	}

	// Verifica se alguma linha foi afetada
	if rowsAffected == 0 {
		utils.ErrorResponse(w, fmt.Errorf("nenhuma carteira atualizada"), http.StatusNotFound)
		return
	}

	// Atualiza o ID do usuário na carteira para o retorno
	wallet.UserID = userID

	// Prepara o objeto de resposta
	response := map[string]interface{}{
		"wallet":        wallet,
		"rows_affected": rowsAffected,
	}

	utils.ToJson(w, response)
}
