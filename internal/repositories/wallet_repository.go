package repositories

import (
	"context"
	"fmt"

	"github.com/WagaoCarvalho/backendgoswagger/internal/models"
	"github.com/jackc/pgx/v5"
)

// NewWallet insere uma nova carteira no banco de dados usando transação
func NewWallet(wallet models.Wallet) (bool, error) {
	db := Connect()
	defer db.Close()

	// Inicia uma transação
	tx, err := db.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return false, fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	// Query de inserção da carteira
	query := `
		INSERT INTO wallets (public_key, usr, balance, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		RETURNING public_key
	`

	// Executa a query dentro da transação
	err = tx.QueryRow(context.Background(), query,
		wallet.PublicKey,
		wallet.UserID,
		wallet.Balance,
	).Scan(&wallet.PublicKey)

	if err != nil {
		return false, fmt.Errorf("erro ao inserir carteira: %w", err)
	}

	// Confirma a transação
	if err = tx.Commit(context.Background()); err != nil {
		return false, fmt.Errorf("erro ao commitar transação: %w", err)
	}

	fmt.Printf("Carteira criada com sucesso! Chave pública: %s\n", wallet.PublicKey)
	return true, nil
}
