package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Database struct {
	ConnURL string
}

type Jwt struct {
	SecretKey string
}

type Config struct {
	Database   Database
	Jwt        Jwt
	ServerPort string
}

func LoadConfig() Config {
	// Carrega o .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar .env")
	}

	return Config{
		Database:   Database{os.Getenv("DB_CONN_URL")},
		Jwt:        Jwt{os.Getenv("JWT_SECRET_KEY")},
		ServerPort: os.Getenv("SERVER_PORT"),
	}
}
