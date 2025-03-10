package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/WagaoCarvalho/backendgoswagger/config"
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
var configs = config.LoadConfig()

func Connect() *pgxpool.Pool {
	url := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		configs.Database.User, configs.Database.Pass, host, port, configs.Database.DbName, sslMode)

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

func TestConnection() {
	con := Connect()
	defer con.Close()

	err := con.Ping(context.Background())
	if err != nil {
		log.Fatalf("Erro ao conectar: %v\n", err)
	}
	fmt.Println("Database connected!")
}
