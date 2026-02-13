package usecases

import (
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/entities"
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

func (uc *CreateProductStockUseCase) Execute(dto CreateProductStockDTO) (string, *domain.Error) {
	productStock, err := entities.NewProductStock(
		nil,
		dto.Name,
		entities.ProductCategory(dto.Category),
		dto.CurrentStock,
		dto.MinimumStock,
		dto.AverageDailySales,
		dto.LeadTimeDays,
		dto.UnitCost,
		entities.CriticalityLevel(dto.CriticalityLevel),
	)
	if err != nil {
		return "", err
	}

	return uc.repo.Create(productStock)
}
