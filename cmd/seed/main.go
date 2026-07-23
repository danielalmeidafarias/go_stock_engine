package main

import (
	"log"

	usecases "github.com/danielalmeidafarias/go_stock_engine/internal/application"
	"github.com/danielalmeidafarias/go_stock_engine/internal/config"
	"github.com/danielalmeidafarias/go_stock_engine/internal/infrastructure/repository/factory"
)

func main() {
	configuration := config.Load()
	repository := factory.NewProductStockRepository(configuration.Database())
	if err := usecases.NewSeedProductStockUseCase(repository).Execute(); err != nil {
		log.Fatal("failed to seed database: ", err)
	}
}
