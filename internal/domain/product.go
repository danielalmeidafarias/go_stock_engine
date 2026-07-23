package domain

type ProductStock struct {
	ID                *string
	Name              string
	Category          string
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
	category string,
	currentStock, minimumStock, averageDailySales, leadTimeDays int,
	unitCost float64,
	criticalityLevel CriticalityLevel,
) (*ProductStock, *Error) {

	errValidation := func() string {
		if name == "" {
			return "name is required"
		}

		if minimumStock < 0 {
			return "minimum stock must be non-negative"
		}

		if averageDailySales < 0 {
			return "average daily sales must be non-negative"
		}

		if leadTimeDays < 0 {
			return "lead time must be non-negative"
		}

		if unitCost <= 0 {
			return "unit cost must be greater than zero"
		}

		if category == "" {
			return "category is required"
		}

		if !IsValidCriticalityLevel(criticalityLevel) {
			return "criticality level must be between 1 and 5"
		}

		return ""
	}()

	if errValidation != "" {
		return nil, NewError(errValidation, ErrBadRequest)
	}

	p := &ProductStock{
		Name:              name,
		Category:          category,
		CurrentStock:      currentStock,
		MinimumStock:      minimumStock,
		AverageDailySales: averageDailySales,
		LeadTimeDays:      leadTimeDays,
		UnitCost:          unitCost,
		CriticalityLevel:  criticalityLevel,
	}

	if id != nil {
		p.ID = id
	}

	return p, nil
}
