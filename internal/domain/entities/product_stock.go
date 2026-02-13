package entities

type ProductStock struct {
	ID                *string          `json:"id"`
	Name              string           `json:"name"`
	Category          ProductCategory  `json:"category"`
	CurrentStock      int              `json:"currentStock"`
	MinimumStock      int              `json:"minimumStock"`
	AverageDailySales int              `json:"averageDailySales"`
	LeadTimeDays      int              `json:"leadTimeDays"`
	UnitCost          float64          `json:"unitCost"`
	CriticalityLevel  CriticalityLevel `json:"criticalityLevel"`
}
