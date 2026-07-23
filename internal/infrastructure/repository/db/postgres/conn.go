package postgres

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/danielalmeidafarias/go_stock_engine/internal/infrastructure/repository/db"
	"github.com/joho/godotenv"
)

const defaultSeedFile = "/usr/local/bin/seed.sql"

func NewPostgresConnection() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, using system environment variables")
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

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := conn.AutoMigrate(&db.ProductStockModel{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	if seedFile := seedFilePath(); seedFile != "" {
		var count int64
		conn.Model(&db.ProductStockModel{}).Count(&count)
		if count == 0 {
			sqlBytes, err := os.ReadFile(seedFile)
			if err != nil {
				log.Printf("seed file not found: %v", err)
			} else {
				if err := conn.Exec(string(sqlBytes)).Error; err != nil {
					log.Printf("seed error: %v", err)
				} else {
					log.Println("seed data loaded successfully")
				}
			}
		}
	}

	return conn
}

func seedFilePath() string {
	if os.Getenv("USE_SEED") != "true" {
		return ""
	}

	if seedFile := os.Getenv("SEED_FILE"); seedFile != "" {
		return seedFile
	}

	return defaultSeedFile
}
