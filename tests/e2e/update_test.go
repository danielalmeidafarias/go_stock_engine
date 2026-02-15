package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"time"
)

func runUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	updateCount := 200
	if updateCount > len(seededIDs) {
		updateCount = len(seededIDs)
	}

	start := time.Now()
	var wg sync.WaitGroup
	sem := make(chan struct{}, 20)

	for i := 0; i < updateCount; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int) {
			defer wg.Done()
			defer func() { <-sem }()

			productID := seededIDs[rand.Intn(len(seededIDs))]

			currentStock := rand.Intn(500)
			criticalityLevel := rand.Intn(5) + 1
			req := updateRequest{
				CurrentStock:     &currentStock,
				CriticalityLevel: &criticalityLevel,
			}

			data, _ := json.Marshal(req)
			httpReq, _ := http.NewRequest(http.MethodPut, baseURL+"/stock/"+productID, bytes.NewReader(data))
			httpReq.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(httpReq)
			if err != nil {
				t.Errorf("PUT %s: %v", baseURL+"/stock/"+productID, err)
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			if resp.StatusCode != http.StatusNoContent {
				t.Errorf("PUT /stock/%s: expected 204, got %d: %s", productID, resp.StatusCode, string(body))
			}
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(start)
	t.Logf("UPDATE: %d requests in %v (%.1f req/sec)", updateCount, elapsed, float64(updateCount)/elapsed.Seconds())
}
