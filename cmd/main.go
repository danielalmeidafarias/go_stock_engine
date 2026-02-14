package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o .env")
	}

	repositoryType := RepositoryType(os.Getenv("REPOSITORY_TYPE"))
	handlerType := HandlerType(os.Getenv("HANDLER_TYPE"))
	paginationDefaultLimit := os.Getenv("PAGINATION_DEFAULT_LIMIT")
	paginationMaxLimit := os.Getenv("PAGINATION_MAX_LIMIT")

	paginationConfig := NewPaginationConfig(paginationDefaultLimit, paginationMaxLimit)

	productStockRepository := ProductStockRepositoryFactory(repositoryType)
	appHadler := AppHandlerFactory(handlerType, paginationConfig, productStockRepository)

	appHadler.Run()
}
