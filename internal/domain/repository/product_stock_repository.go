package repository

import (
	"context"

	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
)

type IProductStockRepository interface {
	Create(ctx context.Context, in *domain.ProductStock) (string, *domain.Error)
	Update(ctx context.Context, in *domain.ProductStock) *domain.Error
	GetAll(ctx context.Context, pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error)
	GetOneByID(ctx context.Context, id string) (*domain.ProductStock, *domain.Error)
	GetByCategory(ctx context.Context, category string, pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error)
	DeleteProductStock(ctx context.Context, id string) *domain.Error
}
