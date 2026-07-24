package usecases

import (
	"context"

	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/repository"
)

type GetOneProductStockUseCase struct {
	repo repository.IProductStockRepository
}

func NewGetOneProductStockUseCase(repo repository.IProductStockRepository) *GetOneProductStockUseCase {
	return &GetOneProductStockUseCase{
		repo: repo,
	}
}

func (uc *GetOneProductStockUseCase) Execute(ctx context.Context, id string) (*domain.ProductStock, *domain.Error) {
	if id == "" {
		return nil, domain.NewError("id is required", domain.ErrBadRequest)
	}

	product, err := uc.repo.GetOneByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}
