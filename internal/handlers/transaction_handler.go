package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"

	"github.com/WagaoCarvalho/backendgoswagger/internal/models"
	"github.com/WagaoCarvalho/backendgoswagger/internal/repositories"
	"github.com/WagaoCarvalho/backendgoswagger/internal/utils"
	"github.com/WagaoCarvalho/backendgoswagger/internal/validations"
)

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	transactions, err := repositories.GetTransactions()
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.ToJson(w, transactions)
}

func PostTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	transaction, err := verifyTransaction(r)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	rowsAffected, err := repositories.NewTransaction(transaction)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"transaction":   transaction,
		"rows_affected": rowsAffected,
		"statusCode":    http.StatusCreated,
	}

	utils.ToJson(w, response)
}

func verifyTransaction(r *http.Request) (models.Transaction, error) {
	// Ler e validar o corpo da requisição
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return models.Transaction{}, errors.New("erro ao ler o corpo da requisição")
	}
	defer r.Body.Close()

	// Definir uma estrutura para os dados esperados no corpo da requisição
	var requestBody struct {
		OriginPublicKey      string  `json:"origin_public_key"`
		DestinationPublicKey string  `json:"destination_public_key"`
		Amount               float32 `json:"amount"`
	}

	if err := json.Unmarshal(body, &requestBody); err != nil {
		return models.Transaction{}, errors.New("JSON inválido")
	}

	// Validar se as chaves públicas foram fornecidas
	if validations.IsStringEmpty(requestBody.OriginPublicKey) {
		return models.Transaction{}, errors.New("chave pública de origem não fornecida")
	}
	if validations.IsStringEmpty(requestBody.DestinationPublicKey) {
		return models.Transaction{}, errors.New("chave pública de destino não fornecida")
	}

	// Buscar a carteira de origem
	originResp, err := repositories.GetWalletByPublicKey(requestBody.OriginPublicKey)
	if err != nil {
		return models.Transaction{}, err
	}
	origin := models.Wallet{
		PublicKey: originResp.Wallet.PublicKey,
		UserID:    originResp.Wallet.UserID,
		Balance:   originResp.Wallet.Balance,
		UpdatedAt: originResp.Wallet.UpdatedAt,
	}

	// Buscar a carteira de destino
	targetResp, err := repositories.GetWalletByPublicKey(requestBody.DestinationPublicKey)
	if err != nil {
		return models.Transaction{}, err
	}
	target := models.Wallet{
		PublicKey: targetResp.Wallet.PublicKey,
		UserID:    targetResp.Wallet.UserID,
		Balance:   targetResp.Wallet.Balance,
		UpdatedAt: targetResp.Wallet.UpdatedAt,
	}

	// Validar se o valor da transação é válido
	if requestBody.Amount <= 0 || math.IsNaN(float64(requestBody.Amount)) {
		return models.Transaction{}, errors.New("valor inválido para a transação")
	}

	const epsilon = 0.0001 // Margem de erro para precisão de ponto flutuante

	amount := float64(requestBody.Amount) // Convertendo para float64

	if math.Abs(origin.Balance-amount) < epsilon {
		requestBody.Amount = float32(origin.Balance) // Permite transferir o saldo exato
	} else if amount > origin.Balance {
		return models.Transaction{}, errors.New("saldo insuficiente na carteira de origem")
	}

	// Criar e retornar a transação
	transaction := models.Transaction{
		Origin:  origin,
		Target:  target,
		Cash:    requestBody.Amount,
		Message: fmt.Sprintf("%s transferiu %.2f $, para %s", originResp.User.Nickname, requestBody.Amount, targetResp.User.Nickname),
	}

	return transaction, nil
}
