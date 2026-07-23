package repository

import (
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
)

type IProductStockRepository interface {
	Create(in *domain.ProductStock) (string, *domain.Error)
	Update(in *domain.ProductStock) *domain.Error
	GetAll(pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error)
	GetOneByID(id string) (*domain.ProductStock, *domain.Error)
	GetByCategory(category string, pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error)
	DeleteProductStock(id string) *domain.Error
}
