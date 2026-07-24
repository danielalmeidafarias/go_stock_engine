package usecases_test

import (
	"context"
	"testing"

	usecases "github.com/danielalmeidafarias/go_stock_engine/internal/application"
	"github.com/danielalmeidafarias/go_stock_engine/internal/application/tests/mocks"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
)

func TestGetAllProductStock_Success(t *testing.T) {
	repo := &mocks.MockProductStockRepository{
		GetAllFn: func(pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error) {
			return []*domain.ProductStock{
				{Name: "P1"},
				{Name: "P2"},
			}, nil
		},
	}
	config := domain.PaginationConfig{DefaultLimit: 20, MaxLimit: 100}
	uc := usecases.NewGetAllProductStockUseCase(repo, config)

	products, err := uc.Execute(context.Background(), domain.Pagination{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Message)
	}
	if len(products) != 2 {
		t.Errorf("len: got %d, want 2", len(products))
	}
}

func TestGetAllProductStock_RepoError(t *testing.T) {
	repo := &mocks.MockProductStockRepository{
		GetAllFn: func(pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error) {
			return nil, domain.NewError("db error", domain.ErrInternal)
		},
	}
	config := domain.PaginationConfig{DefaultLimit: 20, MaxLimit: 100}
	uc := usecases.NewGetAllProductStockUseCase(repo, config)

	_, err := uc.Execute(context.Background(), domain.Pagination{Page: 1, Limit: 10})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetAllProductStock_PaginationApplied(t *testing.T) {
	var capturedPagination *domain.Pagination
	repo := &mocks.MockProductStockRepository{
		GetAllFn: func(pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error) {
			capturedPagination = pagination
			return []*domain.ProductStock{}, nil
		},
	}
	config := domain.PaginationConfig{DefaultLimit: 20, MaxLimit: 50}
	uc := usecases.NewGetAllProductStockUseCase(repo, config)

	// Limit exceeds max, should be capped to 50
	uc.Execute(context.Background(), domain.Pagination{Page: 1, Limit: 999})

	if capturedPagination == nil {
		t.Fatal("pagination should be passed to repo")
	}
	if capturedPagination.Limit != 50 {
		t.Errorf("Limit should be capped to 50, got %d", capturedPagination.Limit)
	}
}
