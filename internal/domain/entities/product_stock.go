package entities

import (
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
)

type ProductStock struct {
	ID                *string
	Name              string
	Category          ProductCategory
	CurrentStock      int
	MinimumStock      int
	AverageDailySales int
	LeadTimeDays      int
	UnitCost          float64
	CriticalityLevel  CriticalityLevel
}

func NewProductStock(
	id *string,
	name string,
	category ProductCategory,
	currentStock, minimumStock, averageDailySales, leadTimeDays int,
	unitCost float64,
	criticalityLevel CriticalityLevel,
) (*ProductStock, *domain.Error) {

	errValidation := func() string {
		if name == "" {
			return "name is required"
		}

		if currentStock < 0 || minimumStock < 0 || averageDailySales < 0 || leadTimeDays < 0 {
			return "numeric fields must be non-negative"
		}

		if unitCost <= 0 {
			return "unit cost must be greater than zero"
		}

		if !IsValidProductCategory(category) {
			return "invalid product category"
		}

		if !IsValidCriticalityLevel(criticalityLevel) {
			return "criticality level must be between 1 and 5"
		}

		return ""
	}()

	if errValidation != "" {
		return nil, domain.NewError(errValidation, domain.ErrBadRequest)
	}

	return &ProductStock{
		Name:              name,
		Category:          category,
		CurrentStock:      currentStock,
		MinimumStock:      minimumStock,
		AverageDailySales: averageDailySales,
		LeadTimeDays:      leadTimeDays,
		UnitCost:          unitCost,
		CriticalityLevel:  criticalityLevel,
	}, nil
}
