package http

import (
	"math"

	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
)

type errorResponseDTO struct {
	Error string `json:"error" example:"error message"`
}

type createResponseDTO struct {
	ID string `json:"id" example:"uuid"`
}

type createProductStockRequestDTO struct {
	Name              string  `json:"name"`
	Category          string  `json:"category"`
	CurrentStock      int     `json:"current_stock"`
	MinimumStock      int     `json:"minimum_stock"`
	AverageDailySales int     `json:"average_daily_sales"`
	LeadTimeDays      int     `json:"lead_time_days"`
	UnitCost          float64 `json:"unit_cost"`
	CriticalityLevel  int     `json:"criticality_level"`
}

type updateProductStockRequestDTO struct {
	CurrentStock      *int     `json:"current_stock"`
	MinimumStock      *int     `json:"minimum_stock"`
	AverageDailySales *int     `json:"average_daily_sales"`
	LeadTimeDays      *int     `json:"lead_time_days"`
	UnitCost          *float64 `json:"unit_cost"`
	CriticalityLevel  *int     `json:"criticality_level"`
}

// productStockResponseDTO represents a product stock item.
type productStockResponseDTO struct {
	ID                string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name              string  `json:"name" example:"Engine Oil Filter"`
	Category          string  `json:"category" example:"engine"`
	CurrentStock      int     `json:"currentStock" example:"150"`
	MinimumStock      int     `json:"minimumStock" example:"50"`
	AverageDailySales int     `json:"averageDailySales" example:"10"`
	LeadTimeDays      int     `json:"leadTimeDays" example:"7"`
	UnitCost          float64 `json:"unitCost" example:"25.50"`
	CriticalityLevel  int     `json:"criticalityLevel" example:"3"`
}

type restockPrioritiesResponseDTO struct {
	Priorities []restockPriorityResponseDTO `json:"priorities"`
}

// restockPriorityResponseDTO represents a product restock priority.
type restockPriorityResponseDTO struct {
	PartID         string `json:"partId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name           string `json:"name" example:"Engine Oil Filter"`
	CurrentStock   int    `json:"currentStock" example:"15"`
	ProjectedStock int    `json:"projectedStock" example:"-20"`
	MinimumStock   int    `json:"minimumStock" example:"20"`
	UrgencyScore   int    `json:"urgencyScore" example:"45"`
}

func toProductStockResponseDTO(product *domain.ProductStock) productStockResponseDTO {
	return productStockResponseDTO{
		ID:                productID(product),
		Name:              product.Name,
		Category:          product.Category,
		CurrentStock:      product.CurrentStock,
		MinimumStock:      product.MinimumStock,
		AverageDailySales: product.AverageDailySales,
		LeadTimeDays:      product.LeadTimeDays,
		UnitCost:          product.UnitCost,
		CriticalityLevel:  int(product.CriticalityLevel),
	}
}

func toProductStockResponseDTOs(products []*domain.ProductStock) []productStockResponseDTO {
	responses := make([]productStockResponseDTO, len(products))
	for i, product := range products {
		responses[i] = toProductStockResponseDTO(product)
	}

	return responses
}

func toRestockPrioritiesResponseDTO(priorities []domain.ProductStockPriority) restockPrioritiesResponseDTO {
	responses := make([]restockPriorityResponseDTO, len(priorities))
	for i, priority := range priorities {
		responses[i] = toRestockPriorityResponseDTO(priority)
	}

	return restockPrioritiesResponseDTO{Priorities: responses}
}

func toRestockPriorityResponseDTO(priority domain.ProductStockPriority) restockPriorityResponseDTO {
	return restockPriorityResponseDTO{
		PartID:         productID(priority.ProductStock),
		Name:           priority.ProductStock.Name,
		CurrentStock:   priority.ProductStock.CurrentStock,
		ProjectedStock: priority.ProjectedStock,
		MinimumStock:   priority.ProductStock.MinimumStock,
		UrgencyScore:   int(math.Floor(priority.UrgencyScore)),
	}
}

func productID(product *domain.ProductStock) string {
	if product.ID == nil {
		return ""
	}

	return *product.ID
}
