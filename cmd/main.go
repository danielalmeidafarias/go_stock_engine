package main

import (
	"log"
	"os"
	"strconv"

	_ "github.com/danielalmeidafarias/go_stock_engine/docs"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println(".env not found, using system environment variables")
	}

	repositoryType := RepositoryType(os.Getenv("REPOSITORY_TYPE"))
	handlerType := HandlerType(os.Getenv("HANDLER_TYPE"))
	paginationDefaultLimit := os.Getenv("PAGINATION_DEFAULT_LIMIT")
	paginationMaxLimit := os.Getenv("PAGINATION_MAX_LIMIT")
	usePriorityPolicy := false
	if usePriorityPolicyStr := os.Getenv("PRIORITY_USE_POLICY"); usePriorityPolicyStr != "" {
		var err error
		usePriorityPolicy, err = strconv.ParseBool(usePriorityPolicyStr)
		if err != nil {
			panic("invalid PRIORITY_USE_POLICY")
		}
	}
	negativeStockFactor := os.Getenv("PRIORITY_NEGATIVE_STOCK_FACTOR")
	leadTimeFactor := os.Getenv("PRIORITY_LEAD_TIME_FACTOR")
	zeroSalesFactor := os.Getenv("PRIORITY_ZERO_SALES_FACTOR")

	paginationConfig := NewPaginationConfig(paginationDefaultLimit, paginationMaxLimit)
	priorityPolicy := NewPriorityPolicy(usePriorityPolicy, negativeStockFactor, leadTimeFactor, zeroSalesFactor)

	productStockRepository := ProductStockRepositoryFactory(repositoryType)
	appHadler := AppHandlerFactory(handlerType, paginationConfig, productStockRepository, priorityPolicy)

	appHadler.Run()
}
