package http

import (
	"testing"

	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
)

func TestToRestockPriorityResponseDTOFloorsUrgencyScore(t *testing.T) {
	id := "uuid"
	response := toRestockPriorityResponseDTO(domain.ProductStockPriority{
		UrgencyScore: 45.9,
		ProductStock: &domain.ProductStock{ID: &id},
	})

	if response.UrgencyScore != 45 {
		t.Fatalf("UrgencyScore = %d, want 45", response.UrgencyScore)
	}
}
