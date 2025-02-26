package repositories

import (
	"context"
	"fmt"

	"github.com/WagaoCarvalho/backendgoswagger/internal/models"
	"github.com/WagaoCarvalho/backendgoswagger/internal/utils"
	"github.com/jackc/pgx/v5"
)

func GetUsers() ([]models.User, error) {
	db := Connect()
	defer db.Close()

	ctx := context.Background()
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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os resultados: %w", err)
	}

	return users, nil
}

func GetUser(uid int64) (models.User, error) {
	db := Connect()
	defer db.Close()

	var user models.User
	query := `SELECT uid, nickname, email, passwd, status, created_at, updated_at FROM users WHERE uid = $1`

	err := db.QueryRow(context.Background(), query, uid).Scan(
		&user.UID,
		&user.Nickname,
		&user.Email,
		&user.Password,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return user, fmt.Errorf("usuário não encontrado")
		}
		return user, fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	return user, nil
}

func NewUser(user models.User) (bool, error) {
	db := Connect()
	defer db.Close()

	ctx := context.Background()

	tx, err := db.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("erro ao iniciar a transação: %w", err)
	}

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

	if err := tx.Commit(ctx); err != nil {
		return false, fmt.Errorf("erro ao confirmar a transação: %w", err)
	}

	fmt.Printf("Usuário e carteira criados com sucesso! UID: %d\n", user.UID)
	return true, nil
}

func UpdateUser(user *models.User) error {
	db := Connect()
	defer db.Close()

	query := `
		UPDATE users 
		SET nickname = $1, email = $2, status = $3, updated_at = NOW() 
		WHERE uid = $4
		RETURNING uid, nickname, email, status, created_at, updated_at
	`

	err := db.QueryRow(context.Background(), query,
		user.Nickname,
		user.Email,
		user.Status,
		user.UID,
	).Scan(
		&user.UID,
		&user.Nickname,
		&user.Email,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("usuário não encontrado")
		}
		return fmt.Errorf("erro ao atualizar usuário: %w", err)
	}

	return nil
}

func DeleteUser(uid int64) error {
	db := Connect()
	defer db.Close()

	query := `DELETE FROM users WHERE uid = $1`

	cmdTag, err := db.Exec(context.Background(), query, uid)
	if err != nil {
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("usuário não encontrado")
	}

	fmt.Printf("Usuário com UID %d deletado com sucesso!\n", uid)
	return nil
}
