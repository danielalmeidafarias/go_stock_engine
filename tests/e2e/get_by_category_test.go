package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func runGetByCategory(t *testing.T) {
	for _, cat := range categories {
		start := time.Now()
		url := fmt.Sprintf("%s/stock/category/%s?page=1&limit=100", baseURL, cat)
		resp, body := doGet(t, url)
		elapsed := time.Since(start)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("GET %s: expected 200, got %d: %s", url, resp.StatusCode, string(body))
			continue
		}

		var products []productStockResponse
		if err := json.Unmarshal(body, &products); err != nil {
			t.Errorf("GET %s: decode response: %v", url, err)
			continue
		}
		for _, product := range products {
			if product.Category != cat {
				t.Errorf("GET %s: category %q, want %q", url, product.Category, cat)
			}
		}
		t.Logf("GET category=%s: %v", cat, elapsed)
	}
}
