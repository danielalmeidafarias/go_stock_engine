package usecases

import (
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/entities"
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

func (uc *GetOneProductStockUseCase) Execute(id string) (*entities.ProductStock, *domain.Error) {
	if id == "" {
		return nil, domain.NewError("id is required", domain.ErrBadRequest)
	}

	product, err := uc.repo.GetOneByID(id)
	if err != nil {
		return nil, err
	}

	return product, nil
}
