package postgres

import (
	"fmt"
	"log"
	"os"

	orm "github.com/danielalmeidafarias/go_stock_engine/internal/infraestructure/repository/gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

type PostgresConnection struct {
	db *gorm.DB
}

func (pgConn *PostgresConnection) GetORM() *gorm.DB {
	return pgConn.db
}

func NewPostgresConnection() *PostgresConnection {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o .env")
	}

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&orm.ProductStockModel{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	return &PostgresConnection{
		db: db,
	}
}
