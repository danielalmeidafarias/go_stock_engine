package db

import (
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
)

type ProductStockModel struct {
	ID                string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name              string  `gorm:"type:varchar(255);not null"`
	Category          string  `gorm:"type:varchar(100);not null"`
	CurrentStock      int     `gorm:"not null"`
	MinimumStock      int     `gorm:"not null"`
	AverageDailySales int     `gorm:"not null"`
	LeadTimeDays      int     `gorm:"not null"`
	UnitCost          float64 `gorm:"type:numeric(10,2);not null"`
	CriticalityLevel  int     `gorm:"not null"`
}

func (m *ProductStockModel) ToDomain() *domain.ProductStock {
	id := m.ID
	return &domain.ProductStock{
		ID:                &id,
		Name:              m.Name,
		Category:          domain.ProductCategory(m.Category),
		CurrentStock:      m.CurrentStock,
		MinimumStock:      m.MinimumStock,
		AverageDailySales: m.AverageDailySales,
		LeadTimeDays:      m.LeadTimeDays,
		UnitCost:          m.UnitCost,
		CriticalityLevel:  domain.CriticalityLevel(m.CriticalityLevel),
	}
}

func MapProductStockToModel(e *domain.ProductStock) *ProductStockModel {
	model := &ProductStockModel{
		Name:              e.Name,
		Category:          string(e.Category),
		CurrentStock:      e.CurrentStock,
		MinimumStock:      e.MinimumStock,
		AverageDailySales: e.AverageDailySales,
		LeadTimeDays:      e.LeadTimeDays,
		UnitCost:          e.UnitCost,
		CriticalityLevel:  int(e.CriticalityLevel),
	}

	if e.ID != nil {
		model.ID = *e.ID
	}

	return model
}
