package usecases_test

import (
	"testing"

	usecases "github.com/danielalmeidafarias/go_stock_engine/internal/application"
	"github.com/danielalmeidafarias/go_stock_engine/internal/application/tests/mocks"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
)

func TestDeleteProductStock_Success(t *testing.T) {
	id := "uuid-1"
	repo := &mocks.MockProductStockRepository{
		GetOneByIDFn: func(reqId string) (*domain.ProductStock, *domain.Error) {
			return &domain.ProductStock{ID: &id, Name: "Test"}, nil
		},
		DeleteProductStockFn: func(reqId string) *domain.Error {
			return nil
		},
	}
	uc := usecases.NewDeleteProductStockUseCase(repo)

	err := uc.Execute("uuid-1")
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Message)
	}
}

func TestDeleteProductStock_EmptyID(t *testing.T) {
	repo := &mocks.MockProductStockRepository{}
	uc := usecases.NewDeleteProductStockUseCase(repo)

	err := uc.Execute("")
	if err == nil {
		t.Fatal("expected error for empty id")
	}
	if err.ErrCode != domain.ErrBadRequest {
		t.Errorf("ErrCode: got %d, want ErrBadRequest", err.ErrCode)
	}
}

func TestDeleteProductStock_NotFound(t *testing.T) {
	repo := &mocks.MockProductStockRepository{
		GetOneByIDFn: func(id string) (*domain.ProductStock, *domain.Error) {
			return nil, domain.NewError("not found", domain.ErrNotFound)
		},
	}
	uc := usecases.NewDeleteProductStockUseCase(repo)

	err := uc.Execute("nonexistent")
	if err == nil {
		t.Fatal("expected not found error")
	}
	if err.ErrCode != domain.ErrNotFound {
		t.Errorf("ErrCode: got %d, want ErrNotFound", err.ErrCode)
	}
}
