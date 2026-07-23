package usecases_test

import (
	"testing"

	usecases "github.com/danielalmeidafarias/go_stock_engine/internal/application"
	"github.com/danielalmeidafarias/go_stock_engine/internal/application/tests/mocks"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
)

func TestGetByCategoryProductStock_Success(t *testing.T) {
	repo := &mocks.MockProductStockRepository{
		GetByCategoryFn: func(category string, pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error) {
			return []*domain.ProductStock{
				{Name: "Engine Part", Category: "engine"},
			}, nil
		},
	}
	config := domain.PaginationConfig{DefaultLimit: 20, MaxLimit: 100}
	uc := usecases.NewGetByCategoryProductStockUseCase(repo, config)

	products, err := uc.Execute(usecases.GetByCategoryDTO{
		Category:   "engine",
		Pagination: domain.Pagination{Page: 1, Limit: 10},
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err.Message)
	}
	if len(products) != 1 {
		t.Errorf("len: got %d, want 1", len(products))
	}
}

func TestGetByCategoryProductStock_EmptyCategory(t *testing.T) {
	repo := &mocks.MockProductStockRepository{}
	config := domain.PaginationConfig{DefaultLimit: 20, MaxLimit: 100}
	uc := usecases.NewGetByCategoryProductStockUseCase(repo, config)

	_, err := uc.Execute(usecases.GetByCategoryDTO{
		Category:   "",
		Pagination: domain.Pagination{Page: 1, Limit: 10},
	})

	if err == nil {
		t.Fatal("expected error for empty category")
	}
	if err.ErrCode != domain.ErrBadRequest {
		t.Errorf("ErrCode: got %d, want ErrBadRequest", err.ErrCode)
	}
}

func TestGetByCategoryProductStock_RepoError(t *testing.T) {
	repo := &mocks.MockProductStockRepository{
		GetByCategoryFn: func(category string, pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error) {
			return nil, domain.NewError("db error", domain.ErrInternal)
		},
	}
	config := domain.PaginationConfig{DefaultLimit: 20, MaxLimit: 100}
	uc := usecases.NewGetByCategoryProductStockUseCase(repo, config)

	_, err := uc.Execute(usecases.GetByCategoryDTO{
		Category:   "engine",
		Pagination: domain.Pagination{Page: 1, Limit: 10},
	})

	if err == nil {
		t.Fatal("expected error")
	}
}
