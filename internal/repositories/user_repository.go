package repositories

import (
	"context"
	"fmt"

	"github.com/WagaoCarvalho/backendgoswagger/internal/models"
	"github.com/WagaoCarvalho/backendgoswagger/internal/utils"
)

// NewUser insere um novo usuário no banco de dados
func NewUser(user models.User) (bool, error) {
	db := Connect()
	defer db.Close()

	ctx := context.Background()

	tx, err := db.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("erro ao iniciar a transação: %w", err)
	}

	// Inserção do usuário
	queryUser := `
		INSERT INTO users (nickname, email, passwd, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING uid
	`

	hasher := utils.BcryptHasher{}
	hashedPass, err := hasher.HashPassword(user.Password)

	err = tx.QueryRow(ctx, queryUser,
		user.Nickname,
		user.Email,
		hashedPass,
		user.Status,
	).Scan(&user.UID)

	if err != nil {
		tx.Rollback(ctx)
		return false, fmt.Errorf("erro ao inserir usuário: %w", err)
	}

	// Inserção da wallet
	wallet := models.Wallet{UserID: user.UID}
	if err := wallet.GeneratePublicKey(); err != nil {
		tx.Rollback(ctx)
		return false, fmt.Errorf("erro ao gerar chave pública: %w", err)
	}

	queryWallet := `
		INSERT INTO wallets (public_key, usr) 
		VALUES ($1, $2)
	`

	_, err = tx.Exec(ctx, queryWallet, wallet.PublicKey, wallet.UserID)
	if err != nil {
		tx.Rollback(ctx)
		return false, fmt.Errorf("erro ao inserir wallet: %w", err)
	}

	// Confirma a transação
	if err := tx.Commit(ctx); err != nil {
		return false, fmt.Errorf("erro ao confirmar a transação: %w", err)
	}

	fmt.Printf("Usuário e carteira criados com sucesso! UID: %d\n", user.UID)
	return true, nil
}

func GetUsers() ([]models.User, error) {
	db := Connect()
	defer db.Close()

	ctx := context.Background()

	// Query para selecionar todos os usuários
	query := `SELECT uid, nickname, email, passwd, status, created_at, updated_at FROM users`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuários: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UID, &user.Nickname, &user.Email, &user.Password, &user.Status, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, fmt.Errorf("erro ao ler os dados do usuário: %w", err)
		}
		users = append(users, user)
	}

	// Verifica se houve algum erro ao iterar sobre as linhas
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os resultados: %w", err)
	}

	return users, nil
}
