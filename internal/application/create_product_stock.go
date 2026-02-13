package usecases

import (
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/repository"
)

type CreateProductStockUseCase struct {
	repo *repository.IProductStockRepository
}

func NewCreateProductStockUseCase(repo *repository.IProductStockRepository) *CreateProductStockUseCase {
	return &CreateProductStockUseCase{
		repo: repo,
	}
}

// Validar se isso faz sentido, mas na minha cabeça, os enums devem ser recebidos como strings e int aqui
type CreateProductStockDTO struct {
	Name              string  `json:"name"`
	Category          string  `json:"category"`
	CurrentStock      int     `json:"currentStock"`
	MinimumStock      int     `json:"minimumStock"`
	AverageDailySales int     `json:"averageDailySales"`
	LeadTimeDays      int     `json:"leadTimeDays"`
	UnitCost          float64 `json:"unitCost"`
	CriticalityLevel  int     `json:"criticalityLevel"`
}

func (uc *CreateProductStockUseCase) Execute(pagination domain.Pagination) {

}
