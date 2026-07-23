package usecases

import (
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/repository"
)

var seedProductStocks = []CreateProductStockDTO{
	{
		Name:              "Filtro de Óleo",
		Category:          "engine",
		CurrentStock:      15,
		MinimumStock:      20,
		AverageDailySales: 4,
		LeadTimeDays:      5,
		UnitCost:          18.50,
		CriticalityLevel:  3,
	},
	{
		Name:              "Pastilha de Freio",
		Category:          "brakes",
		CurrentStock:      8,
		MinimumStock:      10,
		AverageDailySales: 2,
		LeadTimeDays:      4,
		UnitCost:          42.00,
		CriticalityLevel:  5,
	},
	{
		Name:              "Amortecedor Dianteiro",
		Category:          "suspension",
		CurrentStock:      12,
		MinimumStock:      8,
		AverageDailySales: 1,
		LeadTimeDays:      7,
		UnitCost:          160.00,
		CriticalityLevel:  4,
	},
	{
		Name:              "Vela de Ignição",
		Category:          "electrical",
		CurrentStock:      40,
		MinimumStock:      30,
		AverageDailySales: 3,
		LeadTimeDays:      3,
		UnitCost:          12.00,
		CriticalityLevel:  2,
	},
	{
		Name:              "Kit de Embreagem",
		Category:          "transmission",
		CurrentStock:      5,
		MinimumStock:      6,
		AverageDailySales: 1,
		LeadTimeDays:      10,
		UnitCost:          520.00,
		CriticalityLevel:  5,
	},
}

type SeedProductStockUseCase struct {
	repo repository.IProductStockRepository
}

func NewSeedProductStockUseCase(repo repository.IProductStockRepository) *SeedProductStockUseCase {
	return &SeedProductStockUseCase{repo: repo}
}

func (uc *SeedProductStockUseCase) Execute() *domain.Error {
	products, err := uc.repo.GetAll(nil)
	if err != nil {
		return err
	}
	if len(products) > 0 {
		return nil
	}

	createProductStock := NewCreateProductStockUseCase(uc.repo)
	for _, product := range seedProductStocks {
		if _, err := createProductStock.Execute(product); err != nil {
			return err
		}
	}

	return nil
}
