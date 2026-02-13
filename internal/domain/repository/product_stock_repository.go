package repository

import (
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/entities"
)

type IProductStockRepository interface {
	Create(in entities.ProductStock) (string, *domain.Error)
	Update(in entities.ProductStock) *domain.Error
	GetAll(pagination domain.Pagination) ([]*entities.ProductStock, *domain.Error)
	GetOneByID(id string) (*entities.ProductStock, *domain.Error)
	GetByCategory(category entities.ProductCategory, pagination domain.Pagination) ([]*entities.ProductStock, *domain.Error)
	DeleteProductStock(id string) *domain.Error
}
