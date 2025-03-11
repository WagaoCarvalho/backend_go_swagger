package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/WagaoCarvalho/backendgoswagger/internal/models"
	"github.com/jackc/pgx/v5"
)

func GetTransactions() ([]models.Transaction, error) {
	db := Connect()
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
		SELECT uid, origin, destination, amount, message, created_at 
		FROM transactions 
		ORDER BY created_at DESC
	`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar transações: %w", err)
	}
	defer rows.Close()

	var transactions []models.Transaction

	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(
			&transaction.UID,
			&transaction.Origin.PublicKey,
			&transaction.Target.PublicKey,
			&transaction.Cash,
			&transaction.Message,
			&transaction.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear transação: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre transações: %w", err)
	}

	return transactions, nil
}

func NewTransaction(transaction models.Transaction) (int64, error) {
	db := Connect()
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	debitQuery := `
		UPDATE wallets 
		SET balance = balance - $1
		WHERE public_key = $2
		RETURNING usr
	`
	var senderID int64
	err = tx.QueryRow(ctx, debitQuery, transaction.Cash, transaction.Origin.PublicKey).Scan(&senderID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, fmt.Errorf("saldo insuficiente ou carteira do remetente não encontrada")
		}
		return 0, fmt.Errorf("erro ao debitar carteira: %w", err)
	}

	creditQuery := `
		UPDATE wallets 
		SET balance = balance + $1
		WHERE public_key = $2
		RETURNING usr
	`
	var recipientID int64
	err = tx.QueryRow(ctx, creditQuery, transaction.Cash, transaction.Target.PublicKey).Scan(&recipientID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, fmt.Errorf("carteira do destinatário não encontrada")
		}
		return 0, fmt.Errorf("erro ao creditar carteira: %w", err)
	}

	insertQuery := `
		INSERT INTO transactions (origin, destination, amount, message, created_at) 
		VALUES ($1, $2, $3, $4, NOW()) 
		RETURNING uid
	`
	var transactionID int64
	err = tx.QueryRow(ctx, insertQuery, transaction.Origin.PublicKey, transaction.Target.PublicKey, transaction.Cash, transaction.Message).Scan(&transactionID)
	if err != nil {
		return 0, fmt.Errorf("erro ao registrar transação: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("erro ao confirmar transação: %w", err)
	}

	return transactionID, nil
}
