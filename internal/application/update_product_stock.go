package usecases

import (
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/entities"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/repository"
)

type UpdateProductStockUseCase struct {
	repo repository.IProductStockRepository
}

func NewUpdateProductStockUseCase(repo repository.IProductStockRepository) *UpdateProductStockUseCase {
	return &UpdateProductStockUseCase{
		repo: repo,
	}
}

type UpdateProductStockDTO struct {
	ID                string
	CurrentStock      *int
	MinimumStock      *int
	AverageDailySales *int
	LeadTimeDays      *int
	UnitCost          *float64
	CriticalityLevel  *int
}

func (uc *UpdateProductStockUseCase) Execute(dto UpdateProductStockDTO) *domain.Error {
	if dto.ID == "" {
		return domain.NewError("id is required", domain.ErrBadRequest)
	}

	p, err := uc.repo.GetOneByID(dto.ID)
	if err != nil {
		return err
	}

	if dto.CurrentStock != nil {
		p.CurrentStock = *dto.CurrentStock
	}

	if dto.MinimumStock != nil {
		p.MinimumStock = *dto.MinimumStock
	}

	if dto.AverageDailySales != nil {
		p.AverageDailySales = *dto.AverageDailySales
	}

	if dto.LeadTimeDays != nil {
		p.LeadTimeDays = *dto.LeadTimeDays
	}

	if dto.UnitCost != nil {
		p.UnitCost = *dto.UnitCost
	}

	if dto.CriticalityLevel != nil {
		p.CriticalityLevel = entities.CriticalityLevel(*dto.CriticalityLevel)
	}

	p, err = entities.NewProductStock(
		&dto.ID,
		p.Name,
		p.Category,
		p.CurrentStock,
		p.MinimumStock,
		p.AverageDailySales,
		p.LeadTimeDays,
		p.UnitCost,
		p.CriticalityLevel,
	)

	if err != nil {
		return err
	}

	return uc.repo.Update(p)
}
