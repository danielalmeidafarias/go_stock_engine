package usecases

import (
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/repository"
)

type GetAllProductStockUseCase struct {
	repo *repository.IProductStockRepository
}

func NewGetAllProductStockUseCase(repo *repository.IProductStockRepository) *GetAllProductStockUseCase {
	return &GetAllProductStockUseCase{
		repo: repo,
	}
}

func (uc *GetAllProductStockUseCase) Execute(pagination domain.Pagination) {

}
