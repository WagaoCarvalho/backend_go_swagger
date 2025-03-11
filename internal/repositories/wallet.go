package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/WagaoCarvalho/backendgoswagger/internal/models"
	"github.com/jackc/pgx/v5"
)

var (
	ErrWalletNotFound = errors.New("Carteira não existe")
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

type ResponseWallet struct {
	User   models.User   `json:"user"`
	Wallet models.Wallet `json:"wallet"`
}

func GetWallets() ([]ResponseWallet, error) {
	con := Connect()
	defer con.Close()

	query := `
		SELECT 
			u.uid, u.nickname, u.email, u.passwd, u.status, u.created_at, u.updated_at,
			w.public_key, w.balance, w.updated_at
		FROM wallets AS w
		INNER JOIN users AS u ON u.uid = w.usr
		ORDER BY w.usr ASC;
	`

	rows, err := con.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar carteiras: %w", err)
	}
	defer rows.Close()

	var response []ResponseWallet

	for rows.Next() {
		var wallet models.Wallet
		var user models.User

		if err := rows.Scan(
			&user.UID, &user.Nickname, &user.Email, &user.Password, &user.Status, &user.CreatedAt, &user.UpdatedAt,
			&wallet.PublicKey, &wallet.Balance, &wallet.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler carteira e usuário: %w", err)
		}

		response = append(response, ResponseWallet{
			User:   user,
			Wallet: wallet,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os resultados: %w", err)
	}

	return response, nil
}

func GetWalletByPublicKey(publicKey string) (ResponseWallet, error) {
	con := Connect()
	defer con.Close()

	query := `
		SELECT 
			u.uid, u.nickname, u.email, u.passwd, u.status, u.created_at, u.updated_at,
			w.public_key, w.balance, w.updated_at
		FROM wallets AS w
		INNER JOIN users AS u ON u.uid = w.usr
		WHERE w.public_key = $1;
	`

	var response ResponseWallet

	err := con.QueryRow(context.Background(), query, publicKey).Scan(
		&response.User.UID, &response.User.Nickname, &response.User.Email, &response.User.Password,
		&response.User.Status, &response.User.CreatedAt, &response.User.UpdatedAt,
		&response.Wallet.PublicKey, &response.Wallet.Balance, &response.Wallet.UpdatedAt,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return ResponseWallet{}, fmt.Errorf("carteira não encontrada")
		}
		return ResponseWallet{}, fmt.Errorf("erro ao buscar carteira: %w", err)
	}

	return response, nil
}

func UpdateWallet(w models.Wallet) (int64, int64, error) {
	db := Connect()
	defer db.Close()

	ctx := context.Background()

	// Inicia a transação
	tx, err := db.Begin(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("erro ao iniciar a transação: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		UPDATE wallets 
		SET balance = $1, updated_at = NOW() 
		WHERE public_key = $2
		RETURNING usr
	`

	var userID int64
	if err := tx.QueryRow(ctx, query, w.Balance, w.PublicKey).Scan(&userID); err != nil {
		if err == pgx.ErrNoRows {
			return 0, 0, fmt.Errorf("carteira não encontrada")
		}
		return 0, 0, fmt.Errorf("erro ao atualizar carteira: %w", err)
	}

	// Verifica o número de linhas afetadas
	cmd, err := tx.Exec(ctx, query, w.Balance, w.PublicKey)
	if err != nil {
		return 0, 0, fmt.Errorf("erro ao executar a atualização: %w", err)
	}

	rowsAffected := cmd.RowsAffected()

	if err := tx.Commit(ctx); err != nil {
		return 0, 0, fmt.Errorf("erro ao confirmar a transação: %w", err)
	}

	return userID, rowsAffected, nil
}

func AddBalanceWallet(w models.Wallet) (int64, int64, error) {
	db := Connect()
	defer db.Close()

	ctx := context.Background()

	// Inicia a transação
	tx, err := db.Begin(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("erro ao iniciar a transação: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		UPDATE wallets 
		SET balance = (balance + $1), updated_at = NOW() 
		WHERE public_key = $2
		RETURNING usr
	`

	var userID int64
	if err := tx.QueryRow(ctx, query, w.Balance, w.PublicKey).Scan(&userID); err != nil {
		if err == pgx.ErrNoRows {
			return 0, 0, fmt.Errorf("carteira não encontrada")
		}
		return 0, 0, fmt.Errorf("erro ao adicionar saldo à carteira: %w", err)
	}

	// Confirma a transação
	if err := tx.Commit(ctx); err != nil {
		return 0, 0, fmt.Errorf("erro ao confirmar a transação: %w", err)
	}

	// Retorna que uma linha foi afetada
	return userID, 1, nil
}
