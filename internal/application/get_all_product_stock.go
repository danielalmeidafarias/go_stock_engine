package usecases

import (
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/entities"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/repository"
)

type GetAllProductStockUseCase struct {
	repo             repository.IProductStockRepository
	paginationConfig domain.PaginationConfig
}

func NewGetAllProductStockUseCase(repo repository.IProductStockRepository, paginationConfig domain.PaginationConfig) *GetAllProductStockUseCase {
	return &GetAllProductStockUseCase{
		repo:             repo,
		paginationConfig: paginationConfig,
	}
}

func (uc *GetAllProductStockUseCase) Execute(pagination domain.Pagination) ([]*entities.ProductStock, *domain.Error) {
	uc.paginationConfig.ApplyPaginationConfig(&pagination)

	products, err := uc.repo.GetAll(&pagination)
	if err != nil {
		return nil, err
	}

	return products, nil
}
