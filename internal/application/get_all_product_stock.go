package usecases

import (
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/entities"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/repository"
)

type GetAllProductStockUseCase struct {
	repo   repository.IProductStockRepository
	config domain.AppConfig
}

func NewGetAllProductStockUseCase(repo repository.IProductStockRepository, appConfig domain.AppConfig) *GetAllProductStockUseCase {
	return &GetAllProductStockUseCase{
		repo:   repo,
		config: appConfig,
	}
}

func (uc *GetAllProductStockUseCase) Execute(pagination domain.Pagination) ([]*entities.ProductStock, *domain.Error) {
	uc.config.ApplyPaginationConfig(&pagination)

	products, err := uc.repo.GetAll(&pagination)
	if err != nil {
		return nil, err
	}

	return products, nil
}
