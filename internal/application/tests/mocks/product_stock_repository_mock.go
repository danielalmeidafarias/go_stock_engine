package mocks

import (
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
)

// MockProductStockRepository is a manual mock for repository.IProductStockRepository.
type MockProductStockRepository struct {
	CreateFn             func(in *domain.ProductStock) (string, *domain.Error)
	UpdateFn             func(in *domain.ProductStock) *domain.Error
	GetAllFn             func(pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error)
	GetOneByIDFn         func(id string) (*domain.ProductStock, *domain.Error)
	GetByCategoryFn      func(category domain.ProductCategory, pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error)
	DeleteProductStockFn func(id string) *domain.Error
}

func (m *MockProductStockRepository) Create(in *domain.ProductStock) (string, *domain.Error) {
	return m.CreateFn(in)
}

func (m *MockProductStockRepository) Update(in *domain.ProductStock) *domain.Error {
	return m.UpdateFn(in)
}

func (m *MockProductStockRepository) GetAll(pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error) {
	return m.GetAllFn(pagination)
}

func (m *MockProductStockRepository) GetOneByID(id string) (*domain.ProductStock, *domain.Error) {
	return m.GetOneByIDFn(id)
}

func (m *MockProductStockRepository) GetByCategory(category domain.ProductCategory, pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error) {
	return m.GetByCategoryFn(category, pagination)
}

func (m *MockProductStockRepository) DeleteProductStock(id string) *domain.Error {
	return m.DeleteProductStockFn(id)
}
