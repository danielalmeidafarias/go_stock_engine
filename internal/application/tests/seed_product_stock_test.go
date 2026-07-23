package usecases_test

import (
	"testing"

	usecases "github.com/danielalmeidafarias/go_stock_engine/internal/application"
	"github.com/danielalmeidafarias/go_stock_engine/internal/application/tests/mocks"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
)

func TestSeedProductStock_CreatesProductsForEmptyRepository(t *testing.T) {
	var created []*domain.ProductStock
	repo := &mocks.MockProductStockRepository{
		GetAllFn: func(*domain.Pagination) ([]*domain.ProductStock, *domain.Error) {
			return nil, nil
		},
		CreateFn: func(product *domain.ProductStock) (string, *domain.Error) {
			created = append(created, product)
			return "id", nil
		},
	}

	err := usecases.NewSeedProductStockUseCase(repo).Execute()

	if err != nil {
		t.Fatalf("unexpected error: %s", err.Message)
	}
	if len(created) != 5 {
		t.Fatalf("created products: got %d, want 5", len(created))
	}
	if created[1].Category != "brakes" {
		t.Errorf("Category: got %q, want brakes", created[1].Category)
	}
}

func TestSeedProductStock_SkipsNonEmptyRepository(t *testing.T) {
	repo := &mocks.MockProductStockRepository{
		GetAllFn: func(*domain.Pagination) ([]*domain.ProductStock, *domain.Error) {
			return []*domain.ProductStock{{Name: "existing"}}, nil
		},
		CreateFn: func(*domain.ProductStock) (string, *domain.Error) {
			t.Fatal("seed must not create products for a non-empty repository")
			return "", nil
		},
	}

	if err := usecases.NewSeedProductStockUseCase(repo).Execute(); err != nil {
		t.Fatalf("unexpected error: %s", err.Message)
	}
}
