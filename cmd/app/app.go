package main

import (
	"strconv"

	usecases "github.com/danielalmeidafarias/go_stock_engine/internal/application"
	"github.com/danielalmeidafarias/go_stock_engine/internal/config"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/repository"
	"github.com/danielalmeidafarias/go_stock_engine/internal/infrastructure/repository/factory"
	"github.com/danielalmeidafarias/go_stock_engine/internal/presentation/http"
)

func ProductStockRepositoryFactory(database config.DatabaseConfig) repository.IProductStockRepository {
	return factory.NewProductStockRepository(database)
}

type HandlerType string

const (
	HTTP HandlerType = "HTTP"
)

func AppHandlerFactory(
	handlerType HandlerType,
	paginationConfig domain.PaginationConfig,
	repo repository.IProductStockRepository,
	priorityPolicy domain.PriorityPolicy,
) domain.App {
	createUC := usecases.NewCreateProductStockUseCase(repo)
	getAllUC := usecases.NewGetAllProductStockUseCase(repo, paginationConfig)
	getOneUC := usecases.NewGetOneProductStockUseCase(repo)
	updateUC := usecases.NewUpdateProductStockUseCase(repo)
	deleteUC := usecases.NewDeleteProductStockUseCase(repo)
	getByCategoryUC := usecases.NewGetByCategoryProductStockUseCase(repo, paginationConfig)
	getPriorityUC := usecases.NewGetProductPriorityUseCase(repo, paginationConfig, priorityPolicy)

	switch handlerType {
	case HTTP:
		productStockHandler := http.NewProductStockHandler(
			createUC,
			getAllUC,
			getOneUC,
			updateUC,
			deleteUC,
			getByCategoryUC,
			getPriorityUC,
		)

		return http.NewGinApp(productStockHandler)
	default:
		panic("invalid handler type")
	}
}

func NewPaginationConfig(paginationDefaultLimitStr, paginationMaxLimitStr string) domain.PaginationConfig {
	paginationDefaultLimit, err := strconv.Atoi(paginationDefaultLimitStr)
	if err != nil {
		panic("bad pagination default limit configuration")
	}

	paginationMaxLimit, err := strconv.Atoi(paginationMaxLimitStr)
	if err != nil {
		panic("bad pagination max limit configuration")

	}

	return domain.PaginationConfig{
		DefaultLimit: paginationDefaultLimit,
		MaxLimit:     paginationMaxLimit,
	}
}

func NewPriorityPolicy(
	usePolicy bool,
	negativeStockFactorStr,
	leadTimeFactorStr,
	zeroSalesFactorStr string,
) domain.PriorityPolicy {
	if !usePolicy {
		return domain.PriorityPolicy{}
	}

	negativeStockFactor, err := strconv.ParseFloat(negativeStockFactorStr, 64)
	if err != nil {
		panic("invalid PRIORITY_NEGATIVE_STOCK_FACTOR")
	}

	leadTimeFactor, err := strconv.ParseFloat(leadTimeFactorStr, 64)
	if err != nil {
		panic("invalid PRIORITY_LEAD_TIME_FACTOR")
	}

	zeroSalesFactor, err := strconv.ParseFloat(zeroSalesFactorStr, 64)
	if err != nil {
		panic("invalid PRIORITY_ZERO_SALES_FACTOR")
	}

	return domain.PriorityPolicy{
		UsePolicy:           true,
		NegativeStockFactor: negativeStockFactor,
		LeadTimeFactor:      leadTimeFactor,
		ZeroSalesFactor:     zeroSalesFactor,
	}
}
