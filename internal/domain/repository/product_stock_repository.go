package repository

import (
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/entities"
)

type CreateProductStockInput struct {
	Name              string
	Category          entities.ProductCategory
	CurrentStock      int
	MinimumStock      int
	AverageDailySales int
	LeadTimeDays      int
	UnitCost          float64
	CriticalityLevel  int
}

type UpdateProductStockInput struct {
	MinimumStock      int
	AverageDailySales int
	LeadTimeDays      int
	UnitCost          float64
	CriticalityLevel  int
}

type IProductStockRepository interface {
	Create(in CreateProductStockInput) (*entities.ProductStock, error)
	Update(in UpdateProductStockInput) (*entities.ProductStock, error)
	GetAll() ([]*entities.ProductStock, error)
	GetOneByID(id string) (*entities.ProductStock, error)
	GetByCategory(category entities.ProductCategory) ([]*entities.ProductStock, error)
	DeleteProductStock(id string) error
}
