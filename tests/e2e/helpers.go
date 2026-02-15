package e2e

import (
	"io"
	"net/http"
	"testing"
)

const baseURL = "http://localhost:8080"

type createRequest struct {
	Name              string  `json:"name"`
	Category          string  `json:"category"`
	CurrentStock      int     `json:"current_stock"`
	MinimumStock      int     `json:"minimum_stock"`
	AverageDailySales int     `json:"average_daily_sales"`
	LeadTimeDays      int     `json:"lead_time_days"`
	UnitCost          float64 `json:"unit_cost"`
	CriticalityLevel  int     `json:"criticality_level"`
}

type createResponse struct {
	ID string `json:"id"`
}

type updateRequest struct {
	CurrentStock      *int     `json:"current_stock,omitempty"`
	MinimumStock      *int     `json:"minimum_stock,omitempty"`
	AverageDailySales *int     `json:"average_daily_sales,omitempty"`
	LeadTimeDays      *int     `json:"lead_time_days,omitempty"`
	UnitCost          *float64 `json:"unit_cost,omitempty"`
	CriticalityLevel  *int     `json:"criticality_level,omitempty"`
}

var categories = []string{"engine", "oil"}

func doGet(t *testing.T, url string) (*http.Response, []byte) {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp, b
}
