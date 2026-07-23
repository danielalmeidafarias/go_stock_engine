package usecases_test

import (
	"context"
	"testing"

	usecases "github.com/danielalmeidafarias/go_stock_engine/internal/application"
	"github.com/danielalmeidafarias/go_stock_engine/internal/application/tests/mocks"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
)

func TestGetOneProductStock_Success(t *testing.T) {
	id := "uuid-1"
	repo := &mocks.MockProductStockRepository{
		GetOneByIDFn: func(reqId string) (*domain.ProductStock, *domain.Error) {
			return &domain.ProductStock{ID: &id, Name: "Test"}, nil
		},
	}
	uc := usecases.NewGetOneProductStockUseCase(repo)

	product, err := uc.Execute(context.Background(), "uuid-1")
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Message)
	}
	if product.Name != "Test" {
		t.Errorf("Name: got %s, want Test", product.Name)
	}
}

func TestGetOneProductStock_EmptyID(t *testing.T) {
	repo := &mocks.MockProductStockRepository{}
	uc := usecases.NewGetOneProductStockUseCase(repo)

	_, err := uc.Execute(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty id")
	}
	if err.ErrCode != domain.ErrBadRequest {
		t.Errorf("ErrCode: got %d, want ErrBadRequest", err.ErrCode)
	}
}

func TestGetOneProductStock_NotFound(t *testing.T) {
	repo := &mocks.MockProductStockRepository{
		GetOneByIDFn: func(id string) (*domain.ProductStock, *domain.Error) {
			return nil, domain.NewError("not found", domain.ErrNotFound)
		},
	}
	uc := usecases.NewGetOneProductStockUseCase(repo)

	_, err := uc.Execute(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected not found error")
	}
	if err.ErrCode != domain.ErrNotFound {
		t.Errorf("ErrCode: got %d, want ErrNotFound", err.ErrCode)
	}
}
