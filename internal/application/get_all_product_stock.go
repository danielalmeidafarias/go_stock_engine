package usecases

import (
	"context"

	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
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

func (uc *GetAllProductStockUseCase) Execute(ctx context.Context, pagination domain.Pagination) ([]*domain.ProductStock, *domain.Error) {
	domain.ApplyPaginationRules(&pagination, uc.paginationConfig)

	products, err := uc.repo.GetAll(ctx, &pagination)
	if err != nil {
		return nil, err
	}

	return products, nil
}
