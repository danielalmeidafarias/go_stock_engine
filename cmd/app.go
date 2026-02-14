package main

import (
	"strconv"

	usecases "github.com/danielalmeidafarias/go_stock_engine/internal/application"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/repository"
	"github.com/danielalmeidafarias/go_stock_engine/internal/infraestructure/repository/db"
	"github.com/danielalmeidafarias/go_stock_engine/internal/infraestructure/repository/db/postgres"
	"github.com/danielalmeidafarias/go_stock_engine/internal/presentation/http"
)

type RepositoryType string

const (
	Postgres RepositoryType = "POSTGRES"
)

func ProductStockRepositoryFactory(repoType RepositoryType) repository.IProductStockRepository {
	switch repoType {
	case Postgres:
		conn := postgres.NewPostgresConnection()
		errMapper := postgres.NewPostgresErrMapper()
		return db.NewProductStockRepository(conn, errMapper)
	default:
		panic("invalid database type")
	}
}

type HandlerType string

const (
	HTTP HandlerType = "HTTP"
)

func AppHandlerFactory(handlerType HandlerType, paginationConfig domain.PaginationConfig, repo repository.IProductStockRepository) domain.App {
	createUC := usecases.NewCreateProductStockUseCase(repo)
	getAllUC := usecases.NewGetAllProductStockUseCase(repo, paginationConfig)
	getOneUC := usecases.NewGetOneProductStockUseCase(repo)
	updateUC := usecases.NewUpdateProductStockUseCase(repo)
	deleteUC := usecases.NewDeleteProductStockUseCase(repo)
	getByCategoryUC := usecases.NewGetByCategoryProductStockUseCase(repo, paginationConfig)
	getPriorityUC := usecases.NewGetProductPriorityUseCase(repo, paginationConfig)

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
