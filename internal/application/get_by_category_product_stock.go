package usecases

import (
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/entities"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/repository"
)

type GetByCategoryProductStockUseCase struct {
	repo             repository.IProductStockRepository
	paginationConfig domain.PaginationConfig
}

func NewGetByCategoryProductStockUseCase(repo repository.IProductStockRepository, config domain.PaginationConfig) *GetByCategoryProductStockUseCase {
	return &GetByCategoryProductStockUseCase{
		repo:             repo,
		paginationConfig: config,
	}
}

type GetByCategoryDTO struct {
	Category   string
	Pagination domain.Pagination
}

func (uc *GetByCategoryProductStockUseCase) Execute(dto GetByCategoryDTO) ([]*entities.ProductStock, *domain.Error) {
	category := entities.ProductCategory(dto.Category)

	if !entities.IsValidProductCategory(category) {
		return nil, domain.NewError("invalid product category", domain.ErrBadRequest)
	}

	domain.ApplyPaginationRules(&dto.Pagination, uc.paginationConfig)

	products, err := uc.repo.GetByCategory(category, &dto.Pagination)
	if err != nil {
		return nil, err
	}

	return products, nil
}
