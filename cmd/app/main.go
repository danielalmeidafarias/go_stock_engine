package main

import (
	"strconv"

	_ "github.com/danielalmeidafarias/go_stock_engine/docs"
	"github.com/danielalmeidafarias/go_stock_engine/internal/config"
)

func main() {
	configuration := config.Load()

	database := configuration.Database()
	handlerType := HandlerType(configuration.HandlerType())
	pagination := configuration.Pagination()
	policyConfig := configuration.PriorityPolicy()
	usePriorityPolicy := false
	if usePriorityPolicyStr := policyConfig.UsePolicy; usePriorityPolicyStr != "" {
		var err error
		usePriorityPolicy, err = strconv.ParseBool(usePriorityPolicyStr)
		if err != nil {
			panic("invalid PRIORITY_USE_POLICY")
		}
	}
	paginationConfig := NewPaginationConfig(pagination.DefaultLimit, pagination.MaxLimit)
	priorityPolicy := NewPriorityPolicy(
		usePriorityPolicy,
		policyConfig.NegativeStockFactor,
		policyConfig.LeadTimeFactor,
		policyConfig.ZeroSalesFactor,
	)

	productStockRepository := ProductStockRepositoryFactory(database)
	appHandler := AppHandlerFactory(handlerType, paginationConfig, productStockRepository, priorityPolicy)

	appHandler.Run()
}
