package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"
)

func runRestockPriorities(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	pages := []struct{ page, limit int }{
		{1, 20},
		{1, 100},
		{1, 100},
		{5, 100},
	}

	// Warmup
	doGet(t, baseURL+"/restock/priorities?page=1&limit=1")

	for _, p := range pages {
		start := time.Now()
		url := fmt.Sprintf("%s/restock/priorities?page=%d&limit=%d", baseURL, p.page, p.limit)
		resp, body := doGet(t, url)
		elapsed := time.Since(start)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("GET %s: expected 200, got %d: %s", url, resp.StatusCode, string(body))
			continue
		}

		var results []json.RawMessage
		json.Unmarshal(body, &results)
		t.Logf("PRIORITIES page=%d limit=%d: %v (returned %d items)", p.page, p.limit, elapsed, len(results))
	}

	// Measure multiple sequential requests for consistency
	t.Log("--- Priority endpoint latency (10 sequential requests) ---")
	var totalDuration time.Duration
	for i := 0; i < 10; i++ {
		start := time.Now()
		resp, _ := doGet(t, baseURL+"/restock/priorities?page=1&limit=100")
		d := time.Since(start)
		totalDuration += d
		if resp.StatusCode != http.StatusOK {
			t.Errorf("priority request %d: status %d", i, resp.StatusCode)
		}
	}
	avg := totalDuration / 10
	t.Logf("PRIORITIES avg latency over 10 requests: %v (total: %v)", avg, totalDuration)

	// Concurrent priority requests
	t.Log("--- Priority endpoint under load (50 concurrent requests) ---")
	{
		concurrency := 50
		start := time.Now()
		var wg sync.WaitGroup
		var errors int64
		var mu sync.Mutex

		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				resp, _ := doGet(t, baseURL+"/restock/priorities?page=1&limit=100")
				if resp.StatusCode != http.StatusOK {
					mu.Lock()
					errors++
					mu.Unlock()
				}
			}()
		}
		wg.Wait()
		elapsed := time.Since(start)
		t.Logf("PRIORITIES concurrent: %d requests in %v (%.1f req/sec, %d errors)",
			concurrency, elapsed, float64(concurrency)/elapsed.Seconds(), errors)
	}
}
