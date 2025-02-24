package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	user    = "user"
	pass    = "pass"
	host    = "localhost"
	port    = 5432
	dbName  = "db_postgres"
	sslMode = "disable"
)

var dbPool *pgxpool.Pool

// Connect cria uma conexão com o banco de dados PostgreSQL usando pgxpool
func Connect() *pgxpool.Pool {
	url := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		user, pass, host, port, dbName, sslMode)

	// Configura o pool de conexões
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("Erro ao configurar a conexão: %v\n", err)
	}

	// Conecta ao banco de dados
	dbPool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v\n", err)
	}

	fmt.Println("Banco de dados conectado com sucesso!")
	return dbPool
}

// TestConnection verifica se a conexão está funcionando corretamente
func TestConnection() {
	con := Connect()
	defer con.Close()

	err := con.Ping(context.Background())
	if err != nil {
		log.Fatalf("Erro ao conectar: %v\n", err)
	}
	fmt.Println("Database connected!")
}
