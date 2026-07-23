package usecases

import (
	"sort"

	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/repository"
)

type GetProductPriorityUseCase struct {
	repo             repository.IProductStockRepository
	paginationConfig domain.PaginationConfig
	policy           domain.PriorityPolicy
}

func NewGetProductPriorityUseCase(
	repo repository.IProductStockRepository,
	paginationConfig domain.PaginationConfig,
	priorityPolicy domain.PriorityPolicy,
) *GetProductPriorityUseCase {
	return &GetProductPriorityUseCase{
		repo:             repo,
		paginationConfig: paginationConfig,
		policy:           priorityPolicy,
	}
}

func (uc *GetProductPriorityUseCase) Execute(pagination domain.Pagination) ([]domain.ProductStockPriority, *domain.Error) {
	domain.ApplyPaginationRules(&pagination, uc.paginationConfig)

	products, err := uc.repo.GetAll(nil)
	if err != nil {
		return nil, err
	}

	var priorityList []domain.ProductStockPriority

	for _, p := range products {
		stockPriority := p.CalculateStockPriority(uc.policy)
		if stockPriority.RestockNeeded {
			priorityList = append(priorityList, stockPriority)
		}
	}

	sort.Slice(priorityList, func(i, j int) bool {
		return priorityList[i].HasHigherPriorityThan(priorityList[j])
	})

	return domain.PaginatedSlice(priorityList, &pagination), nil
}
