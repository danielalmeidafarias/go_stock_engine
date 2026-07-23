package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func runGetAllProductStock(t *testing.T) {
	pages := []struct{ page, limit int }{
		{1, 20},
		{1, 100},
		{5, 100},
		{50, 100},
	}

	for _, p := range pages {
		start := time.Now()
		url := fmt.Sprintf("%s/stock?page=%d&limit=%d", baseURL, p.page, p.limit)
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
		t.Logf("GET all page=%d limit=%d: %v", p.page, p.limit, elapsed)
	}
}
