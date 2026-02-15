package e2e

import (
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"time"
)

func runGetByID(t *testing.T) {
	sampleSize := 1000
	if sampleSize > len(seededIDs) {
		sampleSize = len(seededIDs)
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, CONCURRENCY)
	var errors int64
	var mu sync.Mutex

	start := time.Now()

	for i := 0; i < sampleSize; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			idx := rand.Intn(len(seededIDs))
			id := seededIDs[idx]
			resp, body := doGet(t, baseURL+"/stock/"+id)
			if resp.StatusCode != http.StatusOK {
				mu.Lock()
				errors++
				mu.Unlock()
				t.Errorf("GET /stock/%s: expected 200, got %d: %s", id, resp.StatusCode, string(body))
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)
	t.Logf("GET by ID: %d requests in %v (%.1f req/sec, avg %.2fms, %d errors)",
		sampleSize, elapsed,
		float64(sampleSize)/elapsed.Seconds(),
		float64(elapsed.Milliseconds())/float64(sampleSize),
		errors)
}
