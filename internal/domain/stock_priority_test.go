package domain

import (
	"testing"
)

func TestCalculateStockPriority_NormalScenario(t *testing.T) {
	p := ProductStock{
		Name:              "Oil Filter",
		Category:          "oil",
		CurrentStock:      100,
		MinimumStock:      50,
		AverageDailySales: 10,
		LeadTimeDays:      5,
		UnitCost:          25.0,
		CriticalityLevel:  High,
	}
	policy := PriorityPolicy{
		UsePolicy:           true,
		NegativeStockFactor: 1.5,
		LeadTimeFactor:      1.2,
		ZeroSalesFactor:     0.5,
	}

	result := p.CalculateStockPriority(policy)

	// expectedConsumption = 10 * 5 = 50
	// projectedStock = 100 - 50 = 50
	// restockNeeded = 50 < 50 => false
	// urgencyScore = ((50 - 50) * 3) + (1.2 * 5) = 6
	if result.ExpectedConsumption != 50 {
		t.Errorf("ExpectedConsumption: got %d, want 50", result.ExpectedConsumption)
	}
	if result.ProjectedStock != 50 {
		t.Errorf("ProjectedStock: got %d, want 50", result.ProjectedStock)
	}
	if result.RestockNeeded != false {
		t.Errorf("RestockNeeded: got %v, want false", result.RestockNeeded)
	}
	if result.UrgencyScore != 6 {
		t.Errorf("UrgencyScore: got %v, want 6", result.UrgencyScore)
	}
}

func TestCalculateStockPriority_PolicyDisabled(t *testing.T) {
	p := ProductStock{
		CurrentStock:      100,
		MinimumStock:      50,
		AverageDailySales: 10,
		LeadTimeDays:      5,
		CriticalityLevel:  High,
	}

	result := p.CalculateStockPriority(PriorityPolicy{
		NegativeStockFactor: 2.0,
		LeadTimeFactor:      2.0,
		ZeroSalesFactor:     0.5,
	})

	if result.UrgencyScore != 0 {
		t.Errorf("UrgencyScore: got %v, want 0", result.UrgencyScore)
	}
}

func TestCalculateStockPriority_PolicyDisabledIgnoresNegativeStockAndLeadTime(t *testing.T) {
	p := ProductStock{
		CurrentStock:      -10,
		MinimumStock:      50,
		AverageDailySales: 5,
		LeadTimeDays:      3,
		CriticalityLevel:  High,
	}

	result := p.CalculateStockPriority(PriorityPolicy{
		NegativeStockFactor: 2.0,
		LeadTimeFactor:      1.2,
	})

	if result.UrgencyScore != 225 {
		t.Errorf("UrgencyScore: got %v, want 225", result.UrgencyScore)
	}
}

func TestCalculateStockPriority_PolicyDisabledIgnoresZeroSalesFactor(t *testing.T) {
	p := ProductStock{
		CurrentStock:     10,
		MinimumStock:     50,
		LeadTimeDays:     5,
		CriticalityLevel: Moderate,
	}

	result := p.CalculateStockPriority(PriorityPolicy{
		LeadTimeFactor:  1.2,
		ZeroSalesFactor: 0.3,
	})

	if result.UrgencyScore != 80 {
		t.Errorf("UrgencyScore: got %v, want 80", result.UrgencyScore)
	}
}

func TestCalculateStockPriority_RestockNeeded(t *testing.T) {
	p := ProductStock{
		Name:              "Engine Part A",
		Category:          "engine",
		CurrentStock:      30,
		MinimumStock:      50,
		AverageDailySales: 10,
		LeadTimeDays:      5,
		UnitCost:          100.0,
		CriticalityLevel:  Critical,
	}
	policy := PriorityPolicy{
		UsePolicy:           true,
		NegativeStockFactor: 1.5,
		LeadTimeFactor:      1.2,
		ZeroSalesFactor:     0.5,
	}

	result := p.CalculateStockPriority(policy)

	// expectedConsumption = 10 * 5 = 50
	// projectedStock = 30 - 50 = -20
	// restockNeeded = -20 < 50 => true
	// urgencyScore = (50 - (-20)) * 5 = 350 + (1.2 * 5) = 356
	expectedScore := 350.0 + (1.2 * 5)
	if result.ExpectedConsumption != 50 {
		t.Errorf("ExpectedConsumption: got %d, want 50", result.ExpectedConsumption)
	}
	if result.ProjectedStock != -20 {
		t.Errorf("ProjectedStock: got %d, want -20", result.ProjectedStock)
	}
	if result.RestockNeeded != true {
		t.Errorf("RestockNeeded: got %v, want true", result.RestockNeeded)
	}
	if result.UrgencyScore != expectedScore {
		t.Errorf("UrgencyScore: got %f, want %f", result.UrgencyScore, expectedScore)
	}
}

func TestCalculateStockPriority_NegativeStock(t *testing.T) {
	p := ProductStock{
		Name:              "Engine Part B",
		Category:          "engine",
		CurrentStock:      -10,
		MinimumStock:      50,
		AverageDailySales: 5,
		LeadTimeDays:      3,
		UnitCost:          50.0,
		CriticalityLevel:  High,
	}
	policy := PriorityPolicy{
		UsePolicy:           true,
		NegativeStockFactor: 2.0,
		LeadTimeFactor:      1.2,
		ZeroSalesFactor:     0.5,
	}

	result := p.CalculateStockPriority(policy)

	// expectedConsumption = 5 * 3 = 15
	// projectedStock = -10 - 15 = -25
	// urgencyScore = (50 - (-25)) * 3 = 225
	// NegativeStockFactor: 2.0 > 1 -> urgencyScore += 2.0 * 10 = 20 -> 245
	// LeadTimeFactor: 1.2 > 1 -> urgencyScore += 1.2 * 3 = 3.6 -> 248.6
	expectedScore := 225.0 + 20.0 + 3.6
	if result.UrgencyScore != expectedScore {
		t.Errorf("UrgencyScore: got %f, want %f", result.UrgencyScore, expectedScore)
	}
	if result.RestockNeeded != true {
		t.Errorf("RestockNeeded: got %v, want true", result.RestockNeeded)
	}
}

func TestCalculateStockPriority_NegativeStockFactorNotApplied(t *testing.T) {
	p := ProductStock{
		Name:              "Part X",
		Category:          "engine",
		CurrentStock:      -5,
		MinimumStock:      20,
		AverageDailySales: 2,
		LeadTimeDays:      2,
		UnitCost:          10.0,
		CriticalityLevel:  Low,
	}
	policy := PriorityPolicy{
		UsePolicy:           true,
		NegativeStockFactor: 0.8,
		LeadTimeFactor:      0.5,
		ZeroSalesFactor:     0.5,
	}

	result := p.CalculateStockPriority(policy)

	// expectedConsumption = 2 * 2 = 4
	// projectedStock = -5 - 4 = -9
	// urgencyScore = (20 - (-9)) * 1 = 29
	if result.UrgencyScore != 29.0 {
		t.Errorf("UrgencyScore: got %f, want 29.0", result.UrgencyScore)
	}
}

func TestCalculateStockPriority_ZeroSalesPositiveUrgency(t *testing.T) {
	p := ProductStock{
		Name:              "Slow Mover",
		Category:          "oil",
		CurrentStock:      10,
		MinimumStock:      50,
		AverageDailySales: 0,
		LeadTimeDays:      5,
		UnitCost:          15.0,
		CriticalityLevel:  Moderate,
	}
	policy := PriorityPolicy{
		UsePolicy:           true,
		NegativeStockFactor: 1.5,
		LeadTimeFactor:      1.2,
		ZeroSalesFactor:     0.3,
	}

	result := p.CalculateStockPriority(policy)

	// urgencyScore = (50 - 10) * 2 = 80
	// LeadTimeFactor: 1.2 > 1 -> urgencyScore += 1.2 * 5 = 6 -> 86
	// ZeroSalesFactor: urgencyScore > 0 -> urgencyScore *= 0.3 -> 25.8
	expectedScore := (80.0 + 6.0) * 0.3
	if result.UrgencyScore != expectedScore {
		t.Errorf("UrgencyScore: got %f, want %f", result.UrgencyScore, expectedScore)
	}
	if result.RestockNeeded != true {
		t.Errorf("RestockNeeded: got %v, want true", result.RestockNeeded)
	}
}

func TestCalculateStockPriority_ZeroSalesNegativeUrgency(t *testing.T) {
	p := ProductStock{
		Name:              "Overstocked No Sales",
		Category:          "oil",
		CurrentStock:      200,
		MinimumStock:      10,
		AverageDailySales: 0,
		LeadTimeDays:      1,
		UnitCost:          5.0,
		CriticalityLevel:  Low,
	}
	policy := PriorityPolicy{
		UsePolicy:           true,
		NegativeStockFactor: 1.5,
		LeadTimeFactor:      0.5,
		ZeroSalesFactor:     0.5,
	}

	result := p.CalculateStockPriority(policy)

	// urgencyScore = (10 - 200) * 1 = -190
	// ZeroSalesFactor: urgencyScore < 0 -> urgencyScore /= 0.5 -> -380
	expectedScore := -190.0 / 0.5
	if result.UrgencyScore != expectedScore {
		t.Errorf("UrgencyScore: got %f, want %f", result.UrgencyScore, expectedScore)
	}
	if result.RestockNeeded != false {
		t.Errorf("RestockNeeded: got %v, want false", result.RestockNeeded)
	}
}

func TestCalculateStockPriority_ZeroSalesFactorNotApplied(t *testing.T) {
	p := ProductStock{
		Name:              "No Sales Factor High",
		Category:          "oil",
		CurrentStock:      10,
		MinimumStock:      50,
		AverageDailySales: 0,
		LeadTimeDays:      2,
		UnitCost:          10.0,
		CriticalityLevel:  Low,
	}
	policy := PriorityPolicy{
		UsePolicy:           true,
		NegativeStockFactor: 1.5,
		LeadTimeFactor:      0.5,
		ZeroSalesFactor:     1.5,
	}

	result := p.CalculateStockPriority(policy)

	// urgencyScore = (50 - 10) * 1 = 40
	if result.UrgencyScore != 40.0 {
		t.Errorf("UrgencyScore: got %f, want 40.0", result.UrgencyScore)
	}
}

func TestCalculateStockPriority_HighLeadTime(t *testing.T) {
	p := ProductStock{
		Name:              "Long Lead Part",
		Category:          "engine",
		CurrentStock:      100,
		MinimumStock:      80,
		AverageDailySales: 5,
		LeadTimeDays:      30,
		UnitCost:          200.0,
		CriticalityLevel:  VeryHigh,
	}
	policy := PriorityPolicy{
		UsePolicy:           true,
		NegativeStockFactor: 1.5,
		LeadTimeFactor:      1.5,
		ZeroSalesFactor:     0.5,
	}

	result := p.CalculateStockPriority(policy)

	// expectedConsumption = 5 * 30 = 150
	// projectedStock = 100 - 150 = -50
	// urgencyScore = (80 - (-50)) * 4 = 520
	// LeadTimeFactor: 1.5 > 1 -> urgencyScore += 1.5 * 30 = 45 -> 565
	expectedScore := 520.0 + 45.0
	if result.UrgencyScore != expectedScore {
		t.Errorf("UrgencyScore: got %f, want %f", result.UrgencyScore, expectedScore)
	}
	if result.RestockNeeded != true {
		t.Errorf("RestockNeeded: got %v, want true", result.RestockNeeded)
	}
}

func TestCalculateStockPriority_InvalidCriticalityDefaultsToOne(t *testing.T) {
	p := ProductStock{
		Name:              "Bad Criticality",
		Category:          "oil",
		CurrentStock:      10,
		MinimumStock:      50,
		AverageDailySales: 5,
		LeadTimeDays:      2,
		UnitCost:          10.0,
		CriticalityLevel:  CriticalityLevel(99),
	}
	policy := PriorityPolicy{
		NegativeStockFactor: 1.5,
		LeadTimeFactor:      0.5,
		ZeroSalesFactor:     0.5,
	}

	result := p.CalculateStockPriority(policy)

	// urgencyScore = (50 - 0) * 1 = 50 (criticality defaults to 1)
	if result.UrgencyScore != 50.0 {
		t.Errorf("UrgencyScore: got %f, want 50.0", result.UrgencyScore)
	}
}
func TestHasHigherPriorityThan_ByUrgencyScore(t *testing.T) {
	a := ProductStockPriority{
		UrgencyScore: 100,
		ProductStock: &ProductStock{Name: "A", CriticalityLevel: Low, AverageDailySales: 1},
	}
	b := ProductStockPriority{
		UrgencyScore: 50,
		ProductStock: &ProductStock{Name: "B", CriticalityLevel: Low, AverageDailySales: 1},
	}

	if !a.HasHigherPriorityThan(b) {
		t.Error("A should have higher priority than B")
	}
}

func TestHasHigherPriorityThan_ByCriticality(t *testing.T) {
	a := ProductStockPriority{
		UrgencyScore: 100,
		ProductStock: &ProductStock{Name: "A", CriticalityLevel: Critical, AverageDailySales: 1},
	}
	b := ProductStockPriority{
		UrgencyScore: 100,
		ProductStock: &ProductStock{Name: "B", CriticalityLevel: Low, AverageDailySales: 1},
	}

	if !a.HasHigherPriorityThan(b) {
		t.Error("A should have higher priority")
	}
}

func TestHasHigherPriorityThan_ByAverageDailySales(t *testing.T) {
	a := ProductStockPriority{
		UrgencyScore: 100,
		ProductStock: &ProductStock{Name: "A", CriticalityLevel: High, AverageDailySales: 20},
	}
	b := ProductStockPriority{
		UrgencyScore: 100,
		ProductStock: &ProductStock{Name: "B", CriticalityLevel: High, AverageDailySales: 5},
	}

	if !a.HasHigherPriorityThan(b) {
		t.Error("A should have higher priority")
	}
}

func TestHasHigherPriorityThan_ByNameAlphabetical(t *testing.T) {
	a := ProductStockPriority{
		UrgencyScore: 100,
		ProductStock: &ProductStock{Name: "Alpha", CriticalityLevel: High, AverageDailySales: 10},
	}
	b := ProductStockPriority{
		UrgencyScore: 100,
		ProductStock: &ProductStock{Name: "Zeta", CriticalityLevel: High, AverageDailySales: 10},
	}

	if !a.HasHigherPriorityThan(b) {
		t.Error("Alpha should have higher priority than Zeta")
	}
}

func TestHasHigherPriorityThan_CaseInsensitiveName(t *testing.T) {
	a := ProductStockPriority{
		UrgencyScore: 100,
		ProductStock: &ProductStock{Name: "alpha", CriticalityLevel: High, AverageDailySales: 10},
	}
	b := ProductStockPriority{
		UrgencyScore: 100,
		ProductStock: &ProductStock{Name: "ALPHA", CriticalityLevel: High, AverageDailySales: 10},
	}

	r1 := a.HasHigherPriorityThan(b)
	r2 := b.HasHigherPriorityThan(a)
	if r1 != r2 {
		t.Error("Same name (case insensitive) should produce symmetric result")
	}
}
