package e2e

import (
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
		t.Logf("GET category=%s: %v", cat, elapsed)
	}
}
