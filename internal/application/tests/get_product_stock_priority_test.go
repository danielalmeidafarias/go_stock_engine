package usecases_test

import (
	"context"
	"testing"

	usecases "github.com/danielalmeidafarias/go_stock_engine/internal/application"
	"github.com/danielalmeidafarias/go_stock_engine/internal/application/tests/mocks"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
)

func TestGetProductPriority_Success(t *testing.T) {
	repo := &mocks.MockProductStockRepository{
		GetAllFn: func(pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error) {
			return []*domain.ProductStock{
				{Name: "A", CurrentStock: 10, MinimumStock: 50, AverageDailySales: 5, LeadTimeDays: 3, UnitCost: 10, CriticalityLevel: domain.Critical},
				{Name: "B", CurrentStock: 200, MinimumStock: 50, AverageDailySales: 5, LeadTimeDays: 3, UnitCost: 10, CriticalityLevel: domain.Low},
				{Name: "C", CurrentStock: 5, MinimumStock: 100, AverageDailySales: 10, LeadTimeDays: 5, UnitCost: 10, CriticalityLevel: domain.High},
			}, nil
		},
	}
	config := domain.PaginationConfig{DefaultLimit: 20, MaxLimit: 100}
	policy := domain.PriorityPolicy{
		NegativeStockFactor: 1.5,
		LeadTimeFactor:      1.2,
		ZeroSalesFactor:     0.5,
	}
	uc := usecases.NewGetProductPriorityUseCase(repo, config, policy)

	priorities, err := uc.Execute(context.Background(), domain.Pagination{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Message)
	}

	// Product B has projected stock = 200-15=185 > 50, so restockNeeded=false, excluded
	// A and C should be included
	if len(priorities) != 2 {
		t.Fatalf("len: got %d, want 2", len(priorities))
	}

	// Should be sorted by urgency score descending
	if priorities[0].UrgencyScore < priorities[1].UrgencyScore {
		t.Error("priorities should be sorted by urgency score descending")
	}
}

func TestGetProductPriority_EmptyList(t *testing.T) {
	repo := &mocks.MockProductStockRepository{
		GetAllFn: func(pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error) {
			return []*domain.ProductStock{}, nil
		},
	}
	config := domain.PaginationConfig{DefaultLimit: 20, MaxLimit: 100}
	policy := domain.PriorityPolicy{}
	uc := usecases.NewGetProductPriorityUseCase(repo, config, policy)

	priorities, err := uc.Execute(context.Background(), domain.Pagination{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Message)
	}
	if len(priorities) != 0 {
		t.Errorf("len: got %d, want 0", len(priorities))
	}
}

func TestGetProductPriority_RepoError(t *testing.T) {
	repo := &mocks.MockProductStockRepository{
		GetAllFn: func(pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error) {
			return nil, domain.NewError("db error", domain.ErrInternal)
		},
	}
	config := domain.PaginationConfig{DefaultLimit: 20, MaxLimit: 100}
	policy := domain.PriorityPolicy{}
	uc := usecases.NewGetProductPriorityUseCase(repo, config, policy)

	_, err := uc.Execute(context.Background(), domain.Pagination{Page: 1, Limit: 10})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetProductPriority_NoRestockNeeded(t *testing.T) {
	repo := &mocks.MockProductStockRepository{
		GetAllFn: func(pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error) {
			return []*domain.ProductStock{
				{Name: "Well Stocked", CurrentStock: 1000, MinimumStock: 10, AverageDailySales: 1, LeadTimeDays: 1, UnitCost: 10, CriticalityLevel: domain.Low},
			}, nil
		},
	}
	config := domain.PaginationConfig{DefaultLimit: 20, MaxLimit: 100}
	policy := domain.PriorityPolicy{}
	uc := usecases.NewGetProductPriorityUseCase(repo, config, policy)

	priorities, err := uc.Execute(context.Background(), domain.Pagination{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Message)
	}
	if len(priorities) != 0 {
		t.Errorf("no products should need restock, got %d", len(priorities))
	}
}

func TestGetProductPriority_SortOrder(t *testing.T) {
	repo := &mocks.MockProductStockRepository{
		GetAllFn: func(pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error) {
			return []*domain.ProductStock{
				{Name: "Low Priority", CurrentStock: 40, MinimumStock: 50, AverageDailySales: 1, LeadTimeDays: 1, UnitCost: 10, CriticalityLevel: domain.Low},
				{Name: "High Priority", CurrentStock: 0, MinimumStock: 100, AverageDailySales: 20, LeadTimeDays: 10, UnitCost: 10, CriticalityLevel: domain.Critical},
				{Name: "Mid Priority", CurrentStock: 10, MinimumStock: 80, AverageDailySales: 5, LeadTimeDays: 5, UnitCost: 10, CriticalityLevel: domain.High},
			}, nil
		},
	}
	config := domain.PaginationConfig{DefaultLimit: 20, MaxLimit: 100}
	policy := domain.PriorityPolicy{
		NegativeStockFactor: 1.5,
		LeadTimeFactor:      1.2,
		ZeroSalesFactor:     0.5,
	}
	uc := usecases.NewGetProductPriorityUseCase(repo, config, policy)

	priorities, err := uc.Execute(context.Background(), domain.Pagination{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Message)
	}

	for i := 0; i < len(priorities)-1; i++ {
		if priorities[i].HasHigherPriorityThan(priorities[i+1]) == false {
			t.Errorf("priority[%d] should have higher priority than priority[%d]", i, i+1)
		}
	}
}
