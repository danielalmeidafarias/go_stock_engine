package main

import (
	"log"
	"os"

	usecases "github.com/danielalmeidafarias/go_stock_engine/internal/application"
	"github.com/danielalmeidafarias/go_stock_engine/internal/infrastructure/repository/factory"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println(".env not found, using system environment variables")
	}

	repository := factory.NewProductStockRepository(factory.Type(os.Getenv("REPOSITORY_TYPE")))
	if err := usecases.NewSeedProductStockUseCase(repository).Execute(); err != nil {
		log.Fatal("failed to seed database: ", err)
	}
}
