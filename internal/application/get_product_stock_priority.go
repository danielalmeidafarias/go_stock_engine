package usecases

import (
	"sort"
	"sync"

	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/entities"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/repository"
)

type GetProductPriorityUseCase struct {
	repo   repository.IProductStockRepository
	config domain.AppConfig
}

func NewGetProductPriorityUseCase(repo repository.IProductStockRepository, config domain.AppConfig) *GetProductPriorityUseCase {
	return &GetProductPriorityUseCase{
		repo:   repo,
		config: config,
	}
}

type ProductStockPriority struct {
	ExpectedConsumption int
	ProjectedStock      int
	IsRepositionNeeded  bool
	UrgencyScore        int
	ProductStock        *entities.ProductStock
}

func (uc *GetProductPriorityUseCase) Execute(pagination domain.Pagination) ([]ProductStockPriority, *domain.Error) {
	products, err := uc.repo.GetAll(nil)
	if err != nil {
		return nil, err
	}

	var priorityList []ProductStockPriority
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, p := range products {
		wg.Go(func() {
			stockPriority := getStockPriority(p)

			if stockPriority.IsRepositionNeeded {
				mu.Lock()
				priorityList = append(priorityList, *stockPriority)
				mu.Unlock()
			}
		})
	}
	wg.Wait()

	sort.Slice(priorityList, func(i, j int) bool {
		return priorityList[i].UrgencyScore > priorityList[j].UrgencyScore
	})

	uc.config.ApplyPaginationConfig(&pagination)

	offset := (pagination.Page - 1) * pagination.Limit
	if offset >= len(priorityList) {
		return []ProductStockPriority{}, nil
	}

	end := min(offset+pagination.Limit, len(priorityList))

	return priorityList[offset:end], nil
}

func getStockPriority(p *entities.ProductStock) *ProductStockPriority {
	expectedConsumption := p.AverageDailySales * p.LeadTimeDays
	projectedStock := p.CurrentStock - expectedConsumption

	return &ProductStockPriority{
		ProductStock:        p,
		ExpectedConsumption: expectedConsumption,
		ProjectedStock:      projectedStock,
		IsRepositionNeeded:  projectedStock < p.MinimumStock,
		UrgencyScore:        (p.MinimumStock - projectedStock) * int(p.CriticalityLevel),
	}
}
