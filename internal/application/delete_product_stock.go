package usecases

import (
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/repository"
)

type DeleteProductStockUseCase struct {
	repo repository.IProductStockRepository
}

func NewDeleteProductStockUseCase(repo repository.IProductStockRepository) *DeleteProductStockUseCase {
	return &DeleteProductStockUseCase{
		repo: repo,
	}
}

func (uc *DeleteProductStockUseCase) Execute(id string) *domain.Error {
	if id == "" {
		return domain.NewError("id is required", domain.ErrBadRequest)
	}

	_, err := uc.repo.GetOneByID(id)
	if err != nil {
		return err
	}

	return uc.repo.DeleteProductStock(id)
}
