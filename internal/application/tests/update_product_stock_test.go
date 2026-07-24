package usecases_test

import (
	"context"
	"testing"

	usecases "github.com/danielalmeidafarias/go_stock_engine/internal/application"
	"github.com/danielalmeidafarias/go_stock_engine/internal/application/tests/mocks"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
)

func TestUpdateProductStock_Success(t *testing.T) {
	id := "uuid-1"
	existingProduct := &domain.ProductStock{
		ID:                &id,
		Name:              "Engine Oil",
		Category:          "engine",
		CurrentStock:      100,
		MinimumStock:      20,
		AverageDailySales: 5,
		LeadTimeDays:      3,
		UnitCost:          25.0,
		CriticalityLevel:  domain.High,
	}

	repo := &mocks.MockProductStockRepository{
		GetOneByIDFn: func(reqId string) (*domain.ProductStock, *domain.Error) {
			return existingProduct, nil
		},
		UpdateFn: func(in *domain.ProductStock) *domain.Error {
			return nil
		},
	}
	uc := usecases.NewUpdateProductStockUseCase(repo)

	newStock := 200
	err := uc.Execute(context.Background(), usecases.UpdateProductStockDTO{
		ID:           "uuid-1",
		CurrentStock: &newStock,
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err.Message)
	}
}

func TestUpdateProductStock_EmptyID(t *testing.T) {
	repo := &mocks.MockProductStockRepository{}
	uc := usecases.NewUpdateProductStockUseCase(repo)

	err := uc.Execute(context.Background(), usecases.UpdateProductStockDTO{ID: ""})
	if err == nil {
		t.Fatal("expected error for empty id")
	}
	if err.ErrCode != domain.ErrBadRequest {
		t.Errorf("ErrCode: got %d, want ErrBadRequest", err.ErrCode)
	}
}

func TestUpdateProductStock_NotFound(t *testing.T) {
	repo := &mocks.MockProductStockRepository{
		GetOneByIDFn: func(id string) (*domain.ProductStock, *domain.Error) {
			return nil, domain.NewError("not found", domain.ErrNotFound)
		},
	}
	uc := usecases.NewUpdateProductStockUseCase(repo)

	newStock := 50
	err := uc.Execute(context.Background(), usecases.UpdateProductStockDTO{
		ID:           "nonexistent",
		CurrentStock: &newStock,
	})

	if err == nil {
		t.Fatal("expected not found error")
	}
	if err.ErrCode != domain.ErrNotFound {
		t.Errorf("ErrCode: got %d, want ErrNotFound", err.ErrCode)
	}
}

func TestUpdateProductStock_InvalidFieldAfterUpdate(t *testing.T) {
	id := "uuid-1"
	existingProduct := &domain.ProductStock{
		ID:                &id,
		Name:              "Engine Oil",
		Category:          "engine",
		CurrentStock:      100,
		MinimumStock:      20,
		AverageDailySales: 5,
		LeadTimeDays:      3,
		UnitCost:          25.0,
		CriticalityLevel:  domain.High,
	}

	repo := &mocks.MockProductStockRepository{
		GetOneByIDFn: func(reqId string) (*domain.ProductStock, *domain.Error) {
			return existingProduct, nil
		},
	}
	uc := usecases.NewUpdateProductStockUseCase(repo)

	negativeMinStock := -10
	err := uc.Execute(context.Background(), usecases.UpdateProductStockDTO{
		ID:           "uuid-1",
		MinimumStock: &negativeMinStock,
	})

	if err == nil {
		t.Fatal("expected validation error for negative minimum stock")
	}
	if err.ErrCode != domain.ErrBadRequest {
		t.Errorf("ErrCode: got %d, want ErrBadRequest", err.ErrCode)
	}
}

func TestUpdateProductStock_AllFields(t *testing.T) {
	id := "uuid-1"
	existingProduct := &domain.ProductStock{
		ID:                &id,
		Name:              "Engine Oil",
		Category:          "engine",
		CurrentStock:      100,
		MinimumStock:      20,
		AverageDailySales: 5,
		LeadTimeDays:      3,
		UnitCost:          25.0,
		CriticalityLevel:  domain.High,
	}

	var updatedProduct *domain.ProductStock
	repo := &mocks.MockProductStockRepository{
		GetOneByIDFn: func(reqId string) (*domain.ProductStock, *domain.Error) {
			return existingProduct, nil
		},
		UpdateFn: func(in *domain.ProductStock) *domain.Error {
			updatedProduct = in
			return nil
		},
	}
	uc := usecases.NewUpdateProductStockUseCase(repo)

	newStock := 200
	newMinStock := 30
	newSales := 10
	newLeadTime := 7
	newCost := 50.0
	newCriticality := 5

	err := uc.Execute(context.Background(), usecases.UpdateProductStockDTO{
		ID:                "uuid-1",
		CurrentStock:      &newStock,
		MinimumStock:      &newMinStock,
		AverageDailySales: &newSales,
		LeadTimeDays:      &newLeadTime,
		UnitCost:          &newCost,
		CriticalityLevel:  &newCriticality,
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err.Message)
	}
	if updatedProduct.CurrentStock != 200 {
		t.Errorf("CurrentStock: got %d, want 200", updatedProduct.CurrentStock)
	}
	if updatedProduct.MinimumStock != 30 {
		t.Errorf("MinimumStock: got %d, want 30", updatedProduct.MinimumStock)
	}
	if updatedProduct.AverageDailySales != 10 {
		t.Errorf("AverageDailySales: got %d, want 10", updatedProduct.AverageDailySales)
	}
	if updatedProduct.LeadTimeDays != 7 {
		t.Errorf("LeadTimeDays: got %d, want 7", updatedProduct.LeadTimeDays)
	}
	if updatedProduct.UnitCost != 50.0 {
		t.Errorf("UnitCost: got %f, want 50.0", updatedProduct.UnitCost)
	}
	if updatedProduct.CriticalityLevel != domain.Critical {
		t.Errorf("CriticalityLevel: got %d, want Critical(%d)", updatedProduct.CriticalityLevel, domain.Critical)
	}
}
