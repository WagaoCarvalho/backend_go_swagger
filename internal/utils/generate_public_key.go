package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GeneratePublicKey gera uma chave pública aleatória de 16 bytes (32 caracteres hexadecimais)
func GeneratePublicKey() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar chave pública: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
