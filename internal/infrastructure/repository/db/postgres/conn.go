package postgres

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/danielalmeidafarias/go_stock_engine/internal/infrastructure/repository/db"
)

func NewPostgresConnection(connectionString string) *gorm.DB {
	conn, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := conn.AutoMigrate(&db.ProductStockModel{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	return conn
}
