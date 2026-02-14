package usecases

import (
	"sort"
	"strings"
	"sync"

	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/entities"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/repository"
)

type GetProductPriorityUseCase struct {
	repo             repository.IProductStockRepository
	paginationConfig domain.PaginationConfig
}

func NewGetProductPriorityUseCase(repo repository.IProductStockRepository, paginationConfig domain.PaginationConfig) *GetProductPriorityUseCase {
	return &GetProductPriorityUseCase{
		repo:             repo,
		paginationConfig: paginationConfig,
	}
}

type ProductStockPriority struct {
	ExpectedConsumption int
	ProjectedStock      int
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
			expectedConsumption := p.AverageDailySales * p.LeadTimeDays
			projectedStock := p.CurrentStock - expectedConsumption
			isRepositionNeeded := projectedStock < p.MinimumStock

			if isRepositionNeeded {
				mu.Lock()
				priorityList = append(priorityList, ProductStockPriority{
					ProductStock:        p,
					ExpectedConsumption: expectedConsumption,
					ProjectedStock:      projectedStock,
					UrgencyScore:        (p.MinimumStock - projectedStock) * int(p.CriticalityLevel),
				})
				mu.Unlock()
			}
		})
	}
	wg.Wait()

	sort.Slice(priorityList, func(i, j int) bool {
		x := priorityList[i]
		y := priorityList[j]

		if x.UrgencyScore != y.UrgencyScore {
			return x.UrgencyScore > y.UrgencyScore
		}

		if x.ProductStock.CriticalityLevel != y.ProductStock.CriticalityLevel {
			return x.ProductStock.CriticalityLevel > y.ProductStock.CriticalityLevel
		}

		if x.ProductStock.AverageDailySales != y.ProductStock.AverageDailySales {
			return x.ProductStock.AverageDailySales > y.ProductStock.AverageDailySales
		}

		return strings.ToLower(x.ProductStock.Name) < strings.ToLower(y.ProductStock.Name)
	})

	domain.ApplyPaginationRules(&pagination, uc.paginationConfig)

	return domain.PaginatedSlice(priorityList, &pagination), nil
}
