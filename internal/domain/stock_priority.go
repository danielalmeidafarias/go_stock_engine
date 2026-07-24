package domain

import "strings"

type CriticalityLevel int

const (
	Low CriticalityLevel = iota + 1
	Moderate
	High
	VeryHigh
	Critical
)

func IsValidCriticalityLevel(c CriticalityLevel) bool {
	switch c {
	case Low, Moderate, High, VeryHigh, Critical:
		return true
	default:
		return false
	}
}

type ProductStockPriority struct {
	ExpectedConsumption int
	ProjectedStock      int
	UrgencyScore        float64
	RestockNeeded       bool
	ProductStock        *ProductStock
}

type PriorityPolicy struct {
	UsePolicy           bool
	NegativeStockFactor float64
	LeadTimeFactor      float64
	ZeroSalesFactor     float64
}

// CalculateStockPriority computes the stock replenishment priority for the product
// based on projected consumption during lead time and business priority policies.
//
// The urgency score is primarily determined by the stock deficit, calculated as the
// difference between the minimum required stock and the projected stock after expected
// consumption. This deficit is weighted by the product's criticality level.
//
// Optional policy adjustments, when enabled:
//   - NegativeStockFactor increases urgency when the current stock is below zero.
//   - LeadTimeFactor increases urgency proportionally to the lead time.
//   - ZeroSalesFactor reduces urgency when the product has no average daily sales.
//
// The method returns a ProductStockPriority
func (p ProductStock) CalculateStockPriority(policy PriorityPolicy) ProductStockPriority {
	avgSales := max(p.AverageDailySales, 0)
	leadTime := p.LeadTimeDays
	minStock := max(p.MinimumStock, 0)

	criticality := 1
	if IsValidCriticalityLevel(p.CriticalityLevel) {
		criticality = int(p.CriticalityLevel)
	}

	expectedConsumption := avgSales * leadTime
	projectedStock := p.CurrentStock - expectedConsumption
	restockNeeded := projectedStock < minStock

	urgencyScore := float64((minStock - projectedStock) * criticality)
	urgencyScore = p.ApplyPolicies(urgencyScore, policy)

	return ProductStockPriority{
		ExpectedConsumption: expectedConsumption,
		ProjectedStock:      projectedStock,
		UrgencyScore:        urgencyScore,
		RestockNeeded:       restockNeeded,
		ProductStock:        &p,
	}
}

// ApplyPolicies adjusts the urgency score when priority policies are enabled.
func (p ProductStock) ApplyPolicies(urgencyScore float64, policy PriorityPolicy) float64 {
	if !policy.UsePolicy {
		return urgencyScore
	}

	avgSales := max(p.AverageDailySales, 0)
	leadTime := p.LeadTimeDays

	if p.CurrentStock < 0 && policy.NegativeStockFactor > 1 {
		urgencyScore += (policy.NegativeStockFactor * (float64(p.CurrentStock) * -1))
	}

	if policy.LeadTimeFactor > 1 {
		urgencyScore += (policy.LeadTimeFactor * float64(leadTime))
	}

	if avgSales == 0 && (policy.ZeroSalesFactor < 1 && policy.ZeroSalesFactor > 0) {
		if urgencyScore > 0 {
			urgencyScore *= policy.ZeroSalesFactor
		} else if urgencyScore < 0 {
			urgencyScore /= policy.ZeroSalesFactor
		}
	}

	return urgencyScore
}

// HasHigherPriorityThan determines whether the current ProductStockPriority
// has higher priority than other
//
// The comparison follows this order:
//  1. Higher UrgencyScore takes precedence.
//  2. If tied, higher CriticalityLevel takes precedence.
//  3. If still tied, higher AverageDailySales takes precedence.
//  4. As a final deterministic tie-breaker, products are ordered
//     alphabetically by name (case-insensitive).
//
// This method defines the rule for priority ordering
func (p ProductStockPriority) HasHigherPriorityThan(other ProductStockPriority) bool {
	if p.UrgencyScore != other.UrgencyScore {
		return p.UrgencyScore > other.UrgencyScore
	}

	if p.ProductStock.CriticalityLevel != other.ProductStock.CriticalityLevel {
		return p.ProductStock.CriticalityLevel > other.ProductStock.CriticalityLevel
	}

	if p.ProductStock.AverageDailySales != other.ProductStock.AverageDailySales {
		return p.ProductStock.AverageDailySales > other.ProductStock.AverageDailySales
	}

	return strings.ToLower(p.ProductStock.Name) < strings.ToLower(other.ProductStock.Name)
}
