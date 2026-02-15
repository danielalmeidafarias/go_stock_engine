package usecases_test

import (
	"testing"

	usecases "github.com/danielalmeidafarias/go_stock_engine/internal/application"
	"github.com/danielalmeidafarias/go_stock_engine/internal/application/tests/mocks"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
)

func TestCreateProductStock_Success(t *testing.T) {
	repo := &mocks.MockProductStockRepository{
		CreateFn: func(in *domain.ProductStock) (string, *domain.Error) {
			return "uuid-123", nil
		},
	}
	uc := usecases.NewCreateProductStockUseCase(repo)

	id, err := uc.Execute(usecases.CreateProductStockDTO{
		Name:              "Engine Oil",
		Category:          "engine",
		CurrentStock:      100,
		MinimumStock:      20,
		AverageDailySales: 5,
		LeadTimeDays:      3,
		UnitCost:          25.0,
		CriticalityLevel:  3,
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err.Message)
	}
	if id != "uuid-123" {
		t.Errorf("id: got %s, want uuid-123", id)
	}
}

func TestCreateProductStock_ValidationError(t *testing.T) {
	repo := &mocks.MockProductStockRepository{}
	uc := usecases.NewCreateProductStockUseCase(repo)

	_, err := uc.Execute(usecases.CreateProductStockDTO{
		Name:             "", // missing name
		Category:         "engine",
		UnitCost:         25.0,
		CriticalityLevel: 3,
	})

	if err == nil {
		t.Fatal("expected validation error")
	}
	if err.ErrCode != domain.ErrBadRequest {
		t.Errorf("ErrCode: got %d, want ErrBadRequest", err.ErrCode)
	}
}

func TestCreateProductStock_InvalidCategory(t *testing.T) {
	repo := &mocks.MockProductStockRepository{}
	uc := usecases.NewCreateProductStockUseCase(repo)

	_, err := uc.Execute(usecases.CreateProductStockDTO{
		Name:             "Test",
		Category:         "invalid_category",
		UnitCost:         25.0,
		CriticalityLevel: 3,
	})

	if err == nil {
		t.Fatal("expected error for invalid category")
	}
}

func TestCreateProductStock_RepoError(t *testing.T) {
	repo := &mocks.MockProductStockRepository{
		CreateFn: func(in *domain.ProductStock) (string, *domain.Error) {
			return "", domain.NewError("db error", domain.ErrInternal)
		},
	}
	uc := usecases.NewCreateProductStockUseCase(repo)

	_, err := uc.Execute(usecases.CreateProductStockDTO{
		Name:              "Engine Oil",
		Category:          "engine",
		CurrentStock:      100,
		MinimumStock:      20,
		AverageDailySales: 5,
		LeadTimeDays:      3,
		UnitCost:          25.0,
		CriticalityLevel:  3,
	})

	if err == nil {
		t.Fatal("expected repo error")
	}
	if err.ErrCode != domain.ErrInternal {
		t.Errorf("ErrCode: got %d, want ErrInternal", err.ErrCode)
	}
}
