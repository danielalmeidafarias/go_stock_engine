package tests

import (
	"testing"

	usecases "github.com/danielalmeidafarias/go_stock_engine/internal/application"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/entities"
)

// --- Mock ---
type mockProductStockRepository struct {
	CreateFn             func(in *entities.ProductStock) (string, *domain.Error)
	UpdateFn             func(in *entities.ProductStock) *domain.Error
	GetAllFn             func(pagination *domain.Pagination) ([]*entities.ProductStock, *domain.Error)
	GetOneByIDFn         func(id string) (*entities.ProductStock, *domain.Error)
	GetByCategoryFn      func(category entities.ProductCategory, pagination *domain.Pagination) ([]*entities.ProductStock, *domain.Error)
	DeleteProductStockFn func(id string) *domain.Error
}

func (m *mockProductStockRepository) Create(in *entities.ProductStock) (string, *domain.Error) {
	return m.CreateFn(in)
}
func (m *mockProductStockRepository) Update(in *entities.ProductStock) *domain.Error {
	return m.UpdateFn(in)
}
func (m *mockProductStockRepository) GetAll(p *domain.Pagination) ([]*entities.ProductStock, *domain.Error) {
	return m.GetAllFn(p)
}
func (m *mockProductStockRepository) GetOneByID(id string) (*entities.ProductStock, *domain.Error) {
	return m.GetOneByIDFn(id)
}
func (m *mockProductStockRepository) GetByCategory(c entities.ProductCategory, p *domain.Pagination) ([]*entities.ProductStock, *domain.Error) {
	return m.GetByCategoryFn(c, p)
}
func (m *mockProductStockRepository) DeleteProductStock(id string) *domain.Error {
	return m.DeleteProductStockFn(id)
}

// --- Helpers ---

func cfg() domain.AppConfig {
	return domain.AppConfig{PaginationConfig: domain.PaginationConfig{DefaultLimit: 10, MaxLimit: 100}}
}

func product(name string, stock, min, dailySales, lead int, crit entities.CriticalityLevel) *entities.ProductStock {
	return &entities.ProductStock{
		Name: name, Category: entities.Engine,
		CurrentStock: stock, MinimumStock: min,
		AverageDailySales: dailySales, LeadTimeDays: lead,
		UnitCost: 10.0, CriticalityLevel: crit,
	}
}

func repoWith(products []*entities.ProductStock) *mockProductStockRepository {
	return &mockProductStockRepository{
		GetAllFn: func(_ *domain.Pagination) ([]*entities.ProductStock, *domain.Error) {
			return products, nil
		},
	}
}

func execute(t *testing.T, products []*entities.ProductStock, pagination domain.Pagination) []usecases.ProductStockPriority {
	t.Helper()
	uc := usecases.NewGetProductPriorityUseCase(repoWith(products), cfg())
	result, err := uc.Execute(pagination)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Message)
	}
	return result
}

func page(p, l int) domain.Pagination { return domain.Pagination{Page: p, Limit: l} }

// --- Testes: Cálculo de prioridade (regras de negócio do README) ---

func TestPriority_Calculations(t *testing.T) {
	// README:
	//   expectedConsumption = averageDailySales * leadTimeDays
	//   projectedStock      = currentStock - expectedConsumption
	//   isRepositionNeeded  = projectedStock < minimumStock
	//   urgencyScore        = (minimumStock - projectedStock) * criticalityLevel
	tests := []struct {
		name           string
		p              *entities.ProductStock
		wantConsump    int
		wantProj       int
		wantReposition bool
		wantUrgency    int
	}{
		{"estoque suficiente", product("A", 100, 20, 5, 7, entities.High), 35, 65, false, -135},
		{"reposição necessária", product("B", 10, 20, 5, 7, entities.Critical), 35, -25, true, 225},
		{"projeção exata no mínimo", product("C", 55, 20, 5, 7, entities.Moderate), 35, 20, false, 0},
		{"logo abaixo do mínimo", product("D", 54, 20, 5, 7, entities.VeryHigh), 35, 19, true, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := execute(t, []*entities.ProductStock{tt.p}, page(1, 10))
			if len(r) != 1 {
				t.Fatalf("expected 1 result, got %d", len(r))
			}
			got := r[0]
			if got.ExpectedConsumption != tt.wantConsump {
				t.Errorf("ExpectedConsumption: got %d, want %d", got.ExpectedConsumption, tt.wantConsump)
			}
			if got.ProjectedStock != tt.wantProj {
				t.Errorf("ProjectedStock: got %d, want %d", got.ProjectedStock, tt.wantProj)
			}
			if got.IsRepositionNeeded != tt.wantReposition {
				t.Errorf("IsRepositionNeeded: got %v, want %v", got.IsRepositionNeeded, tt.wantReposition)
			}
			if got.UrgencyScore != tt.wantUrgency {
				t.Errorf("UrgencyScore: got %d, want %d", got.UrgencyScore, tt.wantUrgency)
			}
		})
	}
}

// --- Testes: Cenários extremos (README: estoque negativo, venda zero, lead time alto) ---

func TestPriority_ExtremeScenarios(t *testing.T) {
	tests := []struct {
		name        string
		p           *entities.ProductStock
		wantProj    int
		wantUrgency int
	}{
		// Estoque negativo (currentStock=0, consumption=100 → proj=-100)
		{"estoque negativo", product("X", 0, 100, 10, 10, entities.Critical), -100, 1000},
		// Venda zero (consumption=0 → proj=currentStock)
		{"venda zero", product("Y", 50, 20, 0, 7, entities.Low), 50, -30},
		// Lead time alto (consumption=5*100=500 → proj=100-500=-400)
		{"lead time alto", product("Z", 100, 50, 5, 100, entities.VeryHigh), -400, 1800},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := execute(t, []*entities.ProductStock{tt.p}, page(1, 10))
			got := r[0]
			if got.ProjectedStock != tt.wantProj {
				t.Errorf("ProjectedStock: got %d, want %d", got.ProjectedStock, tt.wantProj)
			}
			if got.UrgencyScore != tt.wantUrgency {
				t.Errorf("UrgencyScore: got %d, want %d", got.UrgencyScore, tt.wantUrgency)
			}
		})
	}
}

// --- Teste: Ordenação por urgência decrescente ---

func TestPriority_SortedByUrgencyDescending(t *testing.T) {
	products := []*entities.ProductStock{
		product("Low", 100, 20, 5, 7, entities.High),       // urgency = -135
		product("High", 30, 50, 8, 5, entities.Critical),   // urgency = 300
		product("Medium", 40, 30, 4, 5, entities.Moderate), // urgency = 20
	}

	r := execute(t, products, page(1, 10))

	expected := []string{"High", "Medium", "Low"}
	for i, name := range expected {
		if r[i].ProductStock.Name != name {
			t.Errorf("position %d: got '%s', want '%s'", i, r[i].ProductStock.Name, name)
		}
	}

	for i := 0; i < len(r)-1; i++ {
		if r[i].UrgencyScore < r[i+1].UrgencyScore {
			t.Errorf("not sorted desc at %d: %d < %d", i, r[i].UrgencyScore, r[i+1].UrgencyScore)
		}
	}
}

// --- Testes: Paginação ---

func TestPriority_Pagination(t *testing.T) {
	products := []*entities.ProductStock{
		product("P1", 10, 50, 5, 7, entities.Critical),
		product("P2", 20, 40, 3, 5, entities.High),
		product("P3", 100, 20, 2, 3, entities.Low),
	}

	tests := []struct {
		name    string
		page    domain.Pagination
		wantLen int
	}{
		{"primeira página", page(1, 2), 2},
		{"segunda página", page(2, 2), 1},
		{"página além do total", page(5, 10), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := execute(t, products, tt.page)
			if len(r) != tt.wantLen {
				t.Errorf("got %d results, want %d", len(r), tt.wantLen)
			}
		})
	}
}

// --- Teste: Erro do repositório ---

func TestPriority_RepoError(t *testing.T) {
	repo := &mockProductStockRepository{
		GetAllFn: func(_ *domain.Pagination) ([]*entities.ProductStock, *domain.Error) {
			return nil, domain.NewError("db error", domain.ErrInternal)
		},
	}

	uc := usecases.NewGetProductPriorityUseCase(repo, cfg())
	result, err := uc.Execute(page(1, 10))

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.ErrCode != domain.ErrInternal {
		t.Errorf("expected ErrInternal, got %d", err.ErrCode)
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

// --- Teste: Lista vazia ---

func TestPriority_EmptyList(t *testing.T) {
	r := execute(t, []*entities.ProductStock{}, page(1, 10))
	if len(r) != 0 {
		t.Errorf("expected 0 results, got %d", len(r))
	}
}
