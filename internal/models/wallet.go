package models

import (
	"time"

	"github.com/WagaoCarvalho/backendgoswagger/internal/utils"
)

type Wallet struct {
	PublicKey string    `json:"public_key"`
	UserID    int64     `json:"user_id"` // Alterado para int64 para evitar consultas aninhadas desnecessárias
	Balance   float64   `json:"balance"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GeneratePublicKey gera uma chave pública única para a carteira
func (w *Wallet) GeneratePublicKey() error {
	publicKey, err := utils.GeneratePublicKey()
	if err != nil {
		return err
	}
	w.PublicKey = publicKey
	return nil
}
