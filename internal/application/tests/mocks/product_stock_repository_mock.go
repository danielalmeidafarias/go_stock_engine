package mocks

import (
	"context"

	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
)

// MockProductStockRepository is a manual mock for repository.IProductStockRepository.
type MockProductStockRepository struct {
	CreateFn             func(in *domain.ProductStock) (string, *domain.Error)
	UpdateFn             func(in *domain.ProductStock) *domain.Error
	GetAllFn             func(pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error)
	GetOneByIDFn         func(id string) (*domain.ProductStock, *domain.Error)
	GetByCategoryFn      func(category string, pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error)
	DeleteProductStockFn func(id string) *domain.Error
}

func (m *MockProductStockRepository) Create(_ context.Context, in *domain.ProductStock) (string, *domain.Error) {
	return m.CreateFn(in)
}

func (m *MockProductStockRepository) Update(_ context.Context, in *domain.ProductStock) *domain.Error {
	return m.UpdateFn(in)
}

func (m *MockProductStockRepository) GetAll(_ context.Context, pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error) {
	return m.GetAllFn(pagination)
}

func (m *MockProductStockRepository) GetOneByID(_ context.Context, id string) (*domain.ProductStock, *domain.Error) {
	return m.GetOneByIDFn(id)
}

func (m *MockProductStockRepository) GetByCategory(_ context.Context, category string, pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error) {
	return m.GetByCategoryFn(category, pagination)
}

func (m *MockProductStockRepository) DeleteProductStock(_ context.Context, id string) *domain.Error {
	return m.DeleteProductStockFn(id)
}
