package main

import (
	"log"
	"os"

	_ "github.com/danielalmeidafarias/go_stock_engine/docs"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, using system environment variables")
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
