package usecases

import (
	"context"

	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/repository"
)

type CreateProductStockUseCase struct {
	repo repository.IProductStockRepository
}

func NewCreateProductStockUseCase(repo repository.IProductStockRepository) *CreateProductStockUseCase {
	return &CreateProductStockUseCase{
		repo: repo,
	}
}

type CreateProductStockDTO struct {
	Name              string
	Category          string
	CurrentStock      int
	MinimumStock      int
	AverageDailySales int
	LeadTimeDays      int
	UnitCost          float64
	CriticalityLevel  int
}

func (uc *CreateProductStockUseCase) Execute(ctx context.Context, dto CreateProductStockDTO) (string, *domain.Error) {
	productStock, err := domain.NewProductStock(
		nil,
		dto.Name,
		dto.Category,
		dto.CurrentStock,
		dto.MinimumStock,
		dto.AverageDailySales,
		dto.LeadTimeDays,
		dto.UnitCost,
		domain.CriticalityLevel(dto.CriticalityLevel),
	)
	if err != nil {
		return "", err
	}

	return uc.repo.Create(ctx, productStock)
}
