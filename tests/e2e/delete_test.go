package e2e

import (
	"net/http"
	"sync"
	"testing"
	"time"
)

func runDelete(t *testing.T) {
	if len(seededIDs) == 0 {
		t.Skip("no products to delete")
	}

	deleteCount := len(seededIDs)
	var wg sync.WaitGroup
	var errors int64
	var mu sync.Mutex
	concurrency := 50
	sem := make(chan struct{}, concurrency)

	t.Logf("Deleting %d products with concurrency %d...", deleteCount, concurrency)
	start := time.Now()

	for _, id := range seededIDs {
		wg.Add(1)
		sem <- struct{}{}
		go func(productID string) {
			defer wg.Done()
			defer func() { <-sem }()

			req, _ := http.NewRequest(http.MethodDelete, baseURL+"/stock/"+productID, nil)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("DELETE /stock/%s: %v", productID, err)
				mu.Lock()
				errors++
				mu.Unlock()
				return
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusNoContent {
				t.Errorf("DELETE /stock/%s: expected 204, got %d", productID, resp.StatusCode)
				mu.Lock()
				errors++
				mu.Unlock()
			}
		}(id)
	}

	wg.Wait()
	elapsed := time.Since(start)
	t.Logf("DELETE: %d requests in %v (%.0f req/sec, %d errors)",
		deleteCount, elapsed, float64(deleteCount)/elapsed.Seconds(), errors)
}
